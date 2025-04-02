package dto

type CreateUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	UUID string `json:"uuid"`
}

type User struct {
	UUID  string `json:"uuid"`
	Login string `json:"login"`
}
