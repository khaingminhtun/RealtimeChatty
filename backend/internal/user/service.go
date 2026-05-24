package user

import (
	"context"
	"fmt"

	"github.com/khaingminhtun/realtimechatty/internal/db"
	"github.com/khaingminhtun/realtimechatty/internal/pkg/dbutils"
)

type UserService interface {
	GetProfile(ctx context.Context, userID int64) (UpdateResponse, error)
	UpdateProfile(ctx context.Context, dto UpdateProfileDTO) (UpdateResponse, error)
	DeleteAccount(ctx context.Context, userID int64) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetProfile(ctx context.Context, userID int64) (UpdateResponse, error) {
	user, err := s.repo.GetProfile(ctx, userID)
	if err != nil {
		return UpdateResponse{}, fmt.Errorf("can't get user profile: %w", err)
	}
	return mapToResponseDTO(user), nil
}

func (s *userService) UpdateProfile(ctx context.Context, dto UpdateProfileDTO) (UpdateResponse, error) {
	// 1. Build params safely by handling optional nil pointers
	params := db.UpdateUserProfileParams{
		ID: dto.ID,
	}

	if dto.DisplayName != nil {
		params.DisplayName = dbutils.NewText(*dto.DisplayName)
	}
	if dto.AvatarURL != nil {
		params.AvatarUrl = dbutils.NewText(*dto.AvatarURL)
	}
	if dto.Timezone != nil {
		params.Timezone = dbutils.NewText(*dto.Timezone)
	}
	if dto.PushToken != nil {
		params.PushToken = dbutils.NewText(*dto.PushToken)
	}

	// 2. Execute database transaction/query via repository
	user, err := s.repo.UpdateProfile(ctx, params)
	if err != nil {
		return UpdateResponse{}, fmt.Errorf("can't update user profile: %w", err)
	}

	// 3. Return mapped clean DTO
	return mapToResponseDTO(user), nil
}

func (s *userService) DeleteAccount(ctx context.Context, userID int64) error {
	error := s.repo.Delete(ctx, userID)
	if error != nil {
		return fmt.Errorf("can't delete user: %w", error)
	}
	return nil
}

// Mapper helper keeps your core service logic incredibly clean
func mapToResponseDTO(user db.User) UpdateResponse {
	return UpdateResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName.String,
		AvatarURL:   user.AvatarUrl.String,
		Bio:         user.Bio.String,
		Timezone:    user.Timezone.String,
		IsVerified:  user.IsVerified.Bool,
		UpdatedAt:   user.UpdatedAt.Time,
	}
}
