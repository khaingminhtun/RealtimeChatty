package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/khaingminhtun/realtimechatty/internal/db"
	"github.com/khaingminhtun/realtimechatty/internal/mail"
	"github.com/khaingminhtun/realtimechatty/internal/pkg/dbutils"
	"github.com/khaingminhtun/realtimechatty/internal/user"
	"github.com/redis/go-redis/v9"
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
		metadata ClientMetadata,
	) (RegisterResponse, error)
	ResendOTP(ctx context.Context, email string) error
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
func (s *authService) VerifyOTP(ctx context.Context, email string, otp string, metadata ClientMetadata) (RegisterResponse, error) {
	redisKey := fmt.Sprintf("otp:%s", email)

	// 1. Fetch the value stored under the email key from Redis
	storedOTP, err := s.rdb.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return RegisterResponse{}, errors.New("verification code has expired or does not exist")
	} else if err != nil {
		return RegisterResponse{}, fmt.Errorf("failed to access security cache memory: %w", err)
	}

	// 2. Validate user input against cache state
	if storedOTP != otp {
		return RegisterResponse{}, errors.New("invalid verification code provided")
	}

	// 3. Fetch user info to perform verification and token generation
	dbUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("user not found: %w", err)
	}

	var res RegisterResponse

	// 4. Run token issuance and verification inside a clean transaction block
	txErr := s.authRepo.WithTransaction(ctx, func(txRepo AuthRepository) error {

		// Mark user verified in the database
		err = s.authRepo.MarkUserAsVerified(ctx, email)
		if err != nil {
			return fmt.Errorf("failed to update user verification profile status: %w", err)
		}

		// Generate Opaque Refresh tokens
		TokenPairResponse, err := s.tokenManager.GenerateAuthTokens(dbUser.ID)
		if err != nil {
			return fmt.Errorf("token errors")
		}

		// Store device tracing details in the sessions table
		session, err := txRepo.CreateSession(ctx, db.CreateSessionParams{
			UserID:     dbUser.ID,
			TokenHash:  TokenPairResponse.HashedSessionToken,
			DeviceName: dbutils.NewText(metadata.DeviceName),
			DeviceType: dbutils.NewText(metadata.DeviceType),
			IpAddress:  dbutils.ParseIPAddress(metadata.IPAddress),
			UserAgent:  dbutils.NewText(metadata.UserAgent),
			ExpiresAt:  dbutils.NewTimestamp(TokenPairResponse.SessionTokenExpiry),
		})
		if err != nil {
			return fmt.Errorf("failed to map session context parameters: %w", err)
		}

		// Link the refresh token
		err = txRepo.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
			UserID:    dbUser.ID,
			SessionID: session.ID,
			TokenHash: TokenPairResponse.HashedRefreshToken,
			ExpiresAt: dbutils.NewTimestamp(TokenPairResponse.RefreshTokenExpiry),
		})
		if err != nil {
			return fmt.Errorf("failed to bind refresh token: %w", err)
		}

		// Prepare the response payload
		res = RegisterResponse{
			UserID:             dbUser.ID,
			Username:           dbUser.Username,
			Email:              dbUser.Email,
			DisplayName:        dbUser.DisplayName.String,
			AccessToken:        TokenPairResponse.AccessToken,
			RefreshToken:       TokenPairResponse.RawRefreshToken,
			SessionToken:       TokenPairResponse.RawSessionToken,
			RefreshTokenExpiry: int(time.Until(TokenPairResponse.RefreshTokenExpiry).Seconds()),
			SessionTokenExpiry: int(time.Until(TokenPairResponse.SessionTokenExpiry).Seconds()),
		}

		return nil
	})

	if txErr != nil {
		return RegisterResponse{}, txErr
	}

	// 5. Clean Up: Delete the OTP key immediately out of Redis cache
	s.rdb.Del(ctx, redisKey)

	return res, nil
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
