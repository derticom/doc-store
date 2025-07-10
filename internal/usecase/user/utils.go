package user

import (
	"errors"
	"regexp"
	"unicode"
)

const minLength = 8

var loginRegexp = regexp.MustCompile(`^[a-zA-Z0-9]{8,}$`)

var (
	errInvalidLogin         = errors.New("login must be at least 8 characters, only letters and digits")
	errInvalidPassTooShort  = errors.New("password must be at least 8 characters long")
	errInvalidPassCase      = errors.New("password must contain both uppercase and lowercase letters")
	errInvalidPassNoDigit   = errors.New("password must contain at least one digit")
	errInvalidPassNoSpecial = errors.New("password must contain at least one special character")
)

func validateLogin(login string) error {
	if !loginRegexp.MatchString(login) {
		return errInvalidLogin
	}
	return nil
}

func validatePassword(password string) error {
	var (
		hasUpper  = false
		hasLower  = false
		hasDigit  = false
		hasSymbol = false
	)

	if len(password) < minLength {
		return errInvalidPassTooShort
	}

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSymbol = true
		}
	}

	if !hasUpper || !hasLower {
		return errInvalidPassCase
	}
	if !hasDigit {
		return errInvalidPassNoDigit
	}
	if !hasSymbol {
		return errInvalidPassNoSpecial
	}

	return nil
}
