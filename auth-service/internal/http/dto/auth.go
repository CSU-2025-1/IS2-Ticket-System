package dto

type AuthenticateRequest struct {
	Challenge string `form:"challenge"`
	Login     string `form:"login"`
	Password  string `form:"password"`
}

type ConsentRequest struct {
	Challenge string `form:"challenge"`
}
