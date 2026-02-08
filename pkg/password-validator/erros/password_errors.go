package password_errors

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrPasswordValidator = errors.New("password validation failed")
	ErrTooShort          = errors.New("password too short")
	ErrNoUppercase       = errors.New("missing uppercase letters")
	ErrNoLowercase       = errors.New("missing lowercase letters")
	ErrNoDigits          = errors.New("missing digits")
	ErrNoSpecialChars    = errors.New("missing special characters")
	ErrLowEntropy        = errors.New("password entropy too low")
)

type PasswordError struct {
	BaseErr error
	Reasons []error
}

func (e *PasswordError) Error() string {
	if len(e.Reasons) == 0 {
		return fmt.Sprintf("%v", e.BaseErr)
	}
	reasonStrs := make([]string, 0, len(e.Reasons))
	for _, r := range e.Reasons {
		if r == nil {
			continue
		}
		reasonStrs = append(reasonStrs, r.Error())
	}
	return fmt.Sprintf("%v: %s", e.BaseErr, strings.Join(reasonStrs, ", "))
}

func (e *PasswordError) Unwrap() error {
	return e.BaseErr
}

func (e *PasswordError) ReasonsList() []string {
	reasonStrs := make([]string, 0, len(e.Reasons))
	for _, v := range e.Reasons {
		if v != nil {
			reasonStrs = append(reasonStrs, v.Error())
		}
	}
	return reasonStrs
}
