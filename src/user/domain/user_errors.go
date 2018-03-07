package domain

// UserError is a custom error from Go built-in error
type UserError struct {
	Code int
}

const (
	UserErrorUsernameEmptyCode = iota
	UserErrorInvalidUsernameLengthCode
	UserErrorPasswordEmptyCode
	UserErrorWrongPassword
	UserErrorPasswordConfirmationNotMatchCode
	UserChangePasswordErrorWrongOldPassword
)

func (e UserError) Error() string {
	switch e.Code {
	case UserErrorUsernameEmptyCode:
		return "Username cannot be empty"
	case UserErrorInvalidUsernameLengthCode:
		return "Username is too short"
	case UserErrorPasswordEmptyCode:
		return "Password cannot be empty"
	case UserErrorWrongPassword:
		return "Wrong password"
	case UserErrorPasswordConfirmationNotMatchCode:
		return "Password confirmation didn't match"
	case UserChangePasswordErrorWrongOldPassword:
		return "Invalid old password"
	default:
		return "Unrecognized user error code"
	}
}
