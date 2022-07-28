package http

type BaseResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type RegisterRequest struct {
	Account  string `json:"account" validate:"required"`
	NickName string `json:"nickName" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	BaseResponse
}

type AuthRequest struct {
	Account  string `json:"account" validate:"required"`
	DeviceID string `json:"deviceID" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	BaseResponse
	Token string `json:"token"`
}

type SendMessageRequest struct {
	To      int64  `json:"to" validate:"required"`
	Content string `json:"content" validate:"required"`
	Flag    int32  `json:"flag" validate:"required"`
}

type SendMessageResponse struct {
	BaseResponse
	MessageID int64 `json:"messageID"`
	Timestamp int64 `json:"timestamp"`
}
