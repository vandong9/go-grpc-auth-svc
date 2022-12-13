package models

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	Status int64  `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
	Token  string `json:"token,omitempty"`
}
