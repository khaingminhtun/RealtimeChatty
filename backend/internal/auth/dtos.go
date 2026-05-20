package auth

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=32,alphanum"`
	Email       string `json:"email" binding:"required,email,max=255"`
	Password    string `json:"password" binding:"required,min=8,max=72"`
	DisplayName string `json:"displayName" binding:"omitempty,max=80"`
}

type RegisterResponse struct {
	UserID      int64  `json:"userId"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName,omitempty"`
}