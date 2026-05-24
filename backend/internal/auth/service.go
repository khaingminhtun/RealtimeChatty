package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/khaingminhtun/realtimechatty/internal/db"
	"github.com/khaingminhtun/realtimechatty/internal/mail"
	"github.com/khaingminhtun/realtimechatty/internal/pkg/dbutils"
	"github.com/khaingminhtun/realtimechatty/internal/user"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(
		ctx context.Context,
		req *RegisterRequest,
	) error

	VerifyOTP(
		ctx context.Context,
		email string,
		otp string,
	) error
	ResendOTP(ctx context.Context, email string) error
	Login(
		ctx context.Context,
		req *LoginRequest,
		metadata ClientMetadata,
	) (LoginResponse, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
}

type authService struct {
	userRepo     user.UserRepository
	authRepo     AuthRepository
	tokenManager TokenManager // Interface/struct that handles access/refresh generation
	mailer       mail.Mailer  // Interface/struct that dispatches background emails
	rdb          *redis.Client
}

func NewAuthService(
	userRepo user.UserRepository,
	authRepo AuthRepository,
	tokenManager TokenManager,
	mailer mail.Mailer,
	rdb *redis.Client,
) AuthService {
	return &authService{
		userRepo:     userRepo,
		authRepo:     authRepo,
		tokenManager: tokenManager,
		mailer:       mailer,
		rdb:          rdb,
	}
}

func (s *authService) Register(
	ctx context.Context,
	req *RegisterRequest,
) error {

	// 1. Validation Pre-flight Check
	_, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return err
	}

	finalDisplayName := req.DisplayName
	if finalDisplayName == "" {
		finalDisplayName = req.Username
	}

	var rawOTP string

	// 2. Wrap profile generation inside a database transaction
	txErr := s.authRepo.WithTransaction(ctx, func(txRepo AuthRepository) error {

		// 3. Register user profile core records (unverified by default)
		newUser, err := s.userRepo.CreateUser(ctx, db.CreateUserParams{
			Username:    req.Username,
			Email:       req.Email,
			DisplayName: dbutils.NewText(finalDisplayName),
		})
		if err != nil {
			return err
		}

		// 4. Save authentication hashed secrets
		userAuth, err := txRepo.CreateUserAuth(ctx, db.CreateUserAuthParams{
			UserID:       newUser.ID, // Assuming user ID is returned here or retrieved
			PasswordHash: hashedPassword,
		})
		fmt.Printf("userauth %v", userAuth)
		if err != nil {
			return err
		}

		// 5. Generate validation parameters (OTP)
		generatedCode, err := mail.GenerateNumericOTP()
		if err != nil {
			return err
		}
		rawOTP = generatedCode

		// 6. Save OTP to Redis with a 15-Minute expiration
		redisKey := fmt.Sprintf("otp:%s", req.Email)
		err = s.rdb.Set(ctx, redisKey, rawOTP, 15*time.Minute).Err()
		if err != nil {
			return fmt.Errorf("failed to save verification code to cache: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	// 7. Dispatch validation email asynchronously
	go func(email, otp, name string) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		_ = s.mailer.SendVerificationEmail(bgCtx, email, otp, name)
	}(req.Email, rawOTP, finalDisplayName)

	// Return initial state without tokens yet
	return nil
}

// VerifyOTP validates the code, marks user as verified, creates a session, and issues tokens
// VerifyOTP validates the code and marks the user as verified. It does NOT issue tokens.
func (s *authService) VerifyOTP(ctx context.Context, email string, otp string) error {
	redisKey := fmt.Sprintf("otp:%s", email)

	// 1. Fetch the value stored under the email key from Redis
	storedOTP, err := s.rdb.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return errors.New("verification code has expired or does not exist")
	} else if err != nil {
		return fmt.Errorf("failed to access security cache memory: %w", err)
	}

	// 2. Validate user input against cache state
	if storedOTP != otp {
		return errors.New("invalid verification code provided")
	}

	// 3. Mark user verified in the database
	err = s.authRepo.MarkUserAsVerified(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to update user verification profile status: %w", err)
	}

	// 4. Clean Up: Delete the OTP key immediately out of Redis cache
	s.rdb.Del(ctx, redisKey)

	return nil
}

func (s *authService) ResendOTP(ctx context.Context, email string) error {
	// 1. Ensure the user profile exists before doing anything
	dbUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("no registration record found for this email address")
	}

	if dbUser.IsVerified.Valid && dbUser.IsVerified.Bool {
		return errors.New("this account has already been verified; please login instead")
	}

	// 3. Generate a brand new validation parameters block (6-digit OTP)
	rawOTP, err := mail.GenerateNumericOTP()
	if err != nil {
		return fmt.Errorf("failed to generate secure security code: %w", err)
	}

	// 4. Overwrite/Save OTP into Redis cache under the same key with a clean 15-minute TTL
	redisKey := fmt.Sprintf("otp:%s", email)
	err = s.rdb.Set(ctx, redisKey, rawOTP, 15*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to refresh security cache memory: %w", err)
	}

	// 5. Dispatch validation email asynchronously
	displayName := dbUser.DisplayName.String
	if displayName == "" {
		displayName = dbUser.Username
	}

	go func(targetEmail, code, name string) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		_ = s.mailer.SendVerificationEmail(bgCtx, targetEmail, code, name)
	}(email, rawOTP, displayName)

	return nil
}

// Login validates user credentials, checks verification status, and spins up a brand new session
func (s *authService) Login(ctx context.Context, req *LoginRequest, metadata ClientMetadata) (LoginResponse, error) {
	// 1. Fetch user profile from the database by email
	dbUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return LoginResponse{}, errors.New("invalid email address or password provided")
	}

	// 2. Security Guardrail: Enforce email verification status check
	if !dbUser.IsVerified.Bool {
		return LoginResponse{}, errors.New("please verify your email address before logging in")
	}

	// 3. Fetch hashed authentication secrets
	userAuth, err := s.authRepo.GetUserAuthByUserID(ctx, dbUser.ID)
	if err != nil {
		return LoginResponse{}, errors.New("invalid email address or password provided")
	}

	// 4. Validate incoming password against the hashed string in the database
	err = CheckPassword(userAuth.PasswordHash, req.Password)
	if err != nil {
		return LoginResponse{}, errors.New("invalid email address or password provided")
	}

	var res LoginResponse

	// 5. Execute fresh session token generation inside an atomic database transaction
	txErr := s.authRepo.WithTransaction(ctx, func(txRepo AuthRepository) error {

		// Generate Opaque Refresh and Session tokens bundled together via tokenManager
		tokenPair, err := s.tokenManager.GenerateAuthTokens(dbUser.ID)
		if err != nil {
			return fmt.Errorf("failed to generate login token infrastructure: %w", err)
		}

		// Store device tracking details in the sessions table for this login event
		session, err := txRepo.CreateSession(ctx, db.CreateSessionParams{
			UserID:     dbUser.ID,
			TokenHash:  tokenPair.HashedSessionToken,
			DeviceName: dbutils.NewText(metadata.DeviceName),
			DeviceType: dbutils.NewText(metadata.DeviceType),
			IpAddress:  dbutils.ParseIPAddress(metadata.IPAddress),
			UserAgent:  dbutils.NewText(metadata.UserAgent),
			ExpiresAt:  dbutils.NewTimestamp(tokenPair.SessionTokenExpiry),
		})
		if err != nil {
			return fmt.Errorf("failed to map login session context parameters: %w", err)
		}

		// Link the refresh token explicitly back to the new session record ID
		err = txRepo.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
			UserID:    dbUser.ID,
			SessionID: session.ID,
			TokenHash: tokenPair.HashedRefreshToken,
			ExpiresAt: dbutils.NewTimestamp(tokenPair.RefreshTokenExpiry),
		})
		if err != nil {
			return fmt.Errorf("failed to bind login refresh token: %w", err)
		}

		// Prepare the identical operational response payload
		res = LoginResponse{
			UserID:             dbUser.ID,
			Username:           dbUser.Username,
			Email:              dbUser.Email,
			DisplayName:        dbUser.DisplayName.String,
			AccessToken:        tokenPair.AccessToken,
			RefreshToken:       tokenPair.RawRefreshToken,
			SessionToken:       tokenPair.RawSessionToken,
			RefreshTokenExpiry: int(time.Until(tokenPair.RefreshTokenExpiry).Seconds()),
			SessionTokenExpiry: int(time.Until(tokenPair.SessionTokenExpiry).Seconds()),
		}

		return nil
	})

	if txErr != nil {
		return LoginResponse{}, txErr
	}

	return res, nil
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	// 1. Verify user exists
	dbUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// Security tip: Return nil or a generic success message so hackers can't scan for registered emails!
		return nil
	}

	// 2. Generate security reset token
	resetToken, err := mail.GenerateNumericOTP() // Returns a secure 6-digit string like "482910"
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// 3. Save to Redis with a 15-Minute Expiration
	redisKey := fmt.Sprintf("reset:%s", email)
	err = s.rdb.Set(ctx, redisKey, resetToken, 15*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to process password reset cache: %w", err)
	}

	// 4. Send the custom styled recovery email asynchronously
	displayName := dbUser.DisplayName.String
	if displayName == "" {
		displayName = dbUser.Username
	}
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		_ = s.mailer.SendPasswordResetEmail(bgCtx, email, resetToken, displayName)
	}()

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	redisKey := fmt.Sprintf("reset:%s", req.Email)

	// 1. Fetch token from Redis
	storedToken, err := s.rdb.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return errors.New("reset token has expired or does not exist")
	} else if err != nil {
		return fmt.Errorf("failed to query safety storage: %w", err)
	}

	// 2. Verify token equality
	if storedToken != strings.ToUpper(req.Token) {
		return errors.New("invalid or incorrect reset token provided")
	}

	// 3. Hash the brand new password securely
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to process incoming credentials string: %w", err)
	}

	// 4. Update password inside PostgreSQL database using a transaction or repository action
	dbUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("user context correlation lost")
	}

	// Assuming your txRepo or authRepo updates credentials records:
	err = s.authRepo.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		UserID:       dbUser.ID,
		PasswordHash: string(newHashedPassword),
	})
	if err != nil {
		return fmt.Errorf("failed to modify security database entries: %w", err)
	}

	// 5. Success! Evict token immediately from Redis cache so it can't be reused
	s.rdb.Del(ctx, redisKey)

	return nil
}
