package domain

// UserError is a custom error from Go built-in error.
type UserError struct {
	Code int
}

const (
	UserErrorUsernameEmptyCode = iota
	UserErrorInvalidUsernameLengthCode
	UserErrorPasswordEmptyCode
	UserErrorWrongPasswordCode
	UserErrorUsernameExistsCode
	UserErrorPasswordConfirmationNotMatchCode
	UserChangePasswordErrorWrongOldPasswordCode
)

func (e UserError) Error() string {
	switch e.Code {
	case UserErrorUsernameEmptyCode:
		return "Username cannot be empty"
	case UserErrorInvalidUsernameLengthCode:
		return "Username is too short"
	case UserErrorPasswordEmptyCode:
		return "Password cannot be empty"
	case UserErrorWrongPasswordCode:
		return "Wrong password"
	case UserErrorUsernameExistsCode:
		return "Username already exists"
	case UserErrorPasswordConfirmationNotMatchCode:
		return "Password confirmation didn't match"
	case UserChangePasswordErrorWrongOldPasswordCode:
		return "Invalid old password"
	default:
		return "Unrecognized user error code"
	}
}
