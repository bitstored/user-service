package validator

import (
	errs "github.com/bitstored/auth-service/errors"
	"regexp"
)

var (
	upperCaseRegex       = regexp.MustCompile(".*[A-Z].*")
	lowerCaseRegex       = regexp.MustCompile(".*[a-z].*")
	numbersRegex         = regexp.MustCompile(".*[0-9].*")
	symbolRegex          = regexp.MustCompile(".*[!@#$&*()_+=?,].*")
	emailPattern         = regexp.MustCompile("^[-a-zA-Z0-9_.+_]{5,}@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]{2,4}$")
	usernamePattern      = regexp.MustCompile("^[-_A-Za-z0-9.$")
	ErrMsgEmailInvalid   = "email is not a valid email address format"
	ErrMsgPasswordShort  = "password is too short: minimal length is 8 symbols"
	ErrMsgPasswordLong   = "password is too long: maximal length is 100 symbols"
	ErrMsgPasswordSymbol = "password should contain at least one symbol"
	ErrMsgPasswordDigit  = "password should contain at least one digit"
	ErrMsgPasswordUpper  = "password should contain at least one uppercase letter"
	ErrMsgPasswordLower  = "password should contain at least one lowercase letter"
)

func Password(passw string) (bool, *errs.Err) {
	pw := []byte(passw)
	if len(passw) > 100 {
		return false, errs.NewError(errs.ErrKindPasswordDoesntComply, ErrMsgPasswordLong)
	}
	if len(passw) < 8 {
		return false, errs.NewError(errs.ErrKindPasswordDoesntComply, ErrMsgPasswordShort)
	}
	complies := true
	var err *errs.Err = nil
	match := upperCaseRegex.Match(pw)
	if !match {
		complies = false
		if err == nil {
			err = errs.NewError(errs.ErrKindPasswordDoesntComply, ErrMsgPasswordUpper)
		} else {
			err.AddMessage(ErrMsgPasswordUpper)
		}
	}
	match = lowerCaseRegex.Match(pw)
	if !match {
		complies = false
		if err == nil {
			err = errs.NewError(errs.ErrKindPasswordDoesntComply, ErrMsgPasswordLower)
		} else {
			err.AddMessage(ErrMsgPasswordLower)
		}
	}
	match = numbersRegex.Match(pw)
	if !match {
		complies = false
		if err == nil {
			err = errs.NewError(errs.ErrKindPasswordDoesntComply, ErrMsgPasswordDigit)
		} else {
			err.AddMessage(ErrMsgPasswordDigit)
		}
	}
	match = symbolRegex.Match(pw)
	if !match {
		complies = false
		if err == nil {
			err = errs.NewError(errs.ErrKindPasswordDoesntComply, ErrMsgPasswordSymbol)
		} else {
			err.AddMessage(ErrMsgPasswordSymbol)
		}
	}

	return complies, err
}

func Email(email string) bool {
	match := emailPattern.Match([]byte(email))

	return match
}

func Username(username string) bool {
	match := usernamePattern.Match([]byte(username))

	return match
}
