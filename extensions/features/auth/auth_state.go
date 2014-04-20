package auth

type AuthState struct {
	UserName              string
	Mechanism             string
	GetPasswordByUserName func(string) string
}
