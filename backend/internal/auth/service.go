package auth

import (
	"context"
	"errors"

	"github.com/khaingminhtun/realtimechatty/internal/db"
	"github.com/khaingminhtun/realtimechatty/internal/pkg/dbutils"
	"github.com/khaingminhtun/realtimechatty/internal/user"
)

type AuthService interface {
	Register(
		ctx context.Context,
		req *RegisterRequest,
	) (RegisterResponse, error)
}

type authService struct {
	userRepo user.UserRepository
	authRepo AuthRepository
}

func NewAuthService(userRepo user.UserRepository, authRepo AuthRepository) AuthService {
	return &authService{userRepo: userRepo, authRepo: authRepo}
}

func (s *authService) Register(
	ctx context.Context,
	req *RegisterRequest,
) (RegisterResponse, error) {

	_, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return RegisterResponse{}, errors.New("email already exists")
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return RegisterResponse{}, err
	}

	finalDisplayName := req.DisplayName
	if finalDisplayName == "" {
		finalDisplayName = req.Username
	}

	user, err := s.userRepo.CreateUser(ctx, db.CreateUserParams{
		Username:    req.Username,
		Email:       req.Email,
		DisplayName: dbutils.NewText(finalDisplayName),
	})
	if err != nil {
		return RegisterResponse{}, err
	}

	// 4. Create auth-specific security credentials
	_, err = s.authRepo.CreateUserAuth(ctx, db.CreateUserAuthParams{
		UserID:       user.ID,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return RegisterResponse{}, err
	}

	// FIX: Returns a flat struct value instead of an address pointer
	return RegisterResponse{
		UserID:      user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName.String,
	}, nil
}
