package authlib

type AuthError struct {
	Message string
	Code    int
	Status  int
	err     error
}

func (ae *AuthError) Error() string {
	return ae.Message
}

func (ae *AuthError) Cause() error {
	return ae.err
}
