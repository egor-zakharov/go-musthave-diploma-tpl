package dto

type LoginUserRequest struct {
	Login    string `json:"login" `
	Password string `json:"password"`
}
