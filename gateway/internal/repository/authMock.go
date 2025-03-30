package repository

type AuthMock struct {
}

func NewAuthMock() *AuthMock {
	return &AuthMock{}
}

func (a *AuthMock) Auth(token string) (map[string]string, error) {
	return map[string]string{
		"mock": "mock",
	}, nil
}
