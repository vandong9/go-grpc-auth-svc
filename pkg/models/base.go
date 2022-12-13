package models

// Error code
const (
	Login_Error = "Login_Err.001"
)

type BaseMessage struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type BaseResponse struct {
	Code    string      `json:"code"`
	Message BaseMessage `json:"message"`
	Data    interface{} `json:"data"`
}
