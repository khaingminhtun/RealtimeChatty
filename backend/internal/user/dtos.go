package user

import "time"

// UpdateRequest handles optional patches
type UpdateRequest struct {
	DisplayName *string `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
	Timezone    *string `json:"timezone"`
	PushToken   *string `json:"push_token"`
}

type UpdateProfileDTO struct {
	ID          int64
	DisplayName *string
	AvatarURL   *string
	Timezone    *string
	PushToken   *string
}

type UpdateResponse struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	AvatarURL   string    `json:"avatar_url"`
	Bio         string    `json:"bio"`
	Timezone    string    `json:"timezone"`
	IsVerified  bool      `json:"is_verified"`
	UpdatedAt   time.Time `json:"updated_at"`
}
