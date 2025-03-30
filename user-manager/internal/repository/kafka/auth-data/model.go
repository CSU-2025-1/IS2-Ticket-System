package auth_data

type User struct {
	UserUUID string `json:"user_uuid"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
