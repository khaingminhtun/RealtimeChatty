package auth

type RegisterRequest struct {
	Username    string         `json:"username" binding:"required,min=3,max=32,alphanum"`
	Email       string         `json:"email" binding:"required,email,max=255"`
	Password    string         `json:"password" binding:"required,min=8,max=72"`
	DisplayName string         `json:"displayName" binding:"omitempty,max=80"`
	Metadata    ClientMetadata `json:"metadata"`
}

type RegisterResponse struct {
	UserID             int64  `json:"userId"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	DisplayName        string `json:"displayName,omitempty"`
	AccessToken        string
	RefreshToken       string
	SessionToken       string
	RefreshTokenExpiry int // cookie Max-Age in seconds
	SessionTokenExpiry int // cookie Max-Age in seconds
}

type ClientMetadata struct {
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"`
	IPAddress  string `json:"ip_address"`
	UserAgent  string `json:"user_agent"`
}

type VerifyOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type VerifyOTPUserResponse struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	AccessToken string `json:"access_token"`
}

type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}
