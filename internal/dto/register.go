package dto

type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
