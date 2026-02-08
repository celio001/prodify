package validator

import (
	"strings"
	password_errors "github.com/celio001/prodify/pkg/password-validator/erros"
)

func Validate(password string, minEntropy float64) error {
	entropy := getEntropy(password)
	if entropy >= minEntropy {
		return nil
	}

	hasReplace := false
	hasSep := false
	hasOtherSpecial := false
	hasLower := false
	hasUpper := false
	hasDigits := false
	for _, c := range password {
		if strings.ContainsRune(replaceChars, c) {
			hasReplace = true
			continue
		}
		if strings.ContainsRune(sepChars, c) {
			hasSep = true
			continue
		}
		if strings.ContainsRune(otherSpecialChars, c) {
			hasOtherSpecial = true
			continue
		}
		if strings.ContainsRune(lowerChars, c) {
			hasLower = true
			continue
		}
		if strings.ContainsRune(upperChars, c) {
			hasUpper = true
			continue
		}
		if strings.ContainsRune(digitsChars, c) {
			hasDigits = true
			continue
		}
	}

	var errs []error

	if !hasOtherSpecial || !hasSep || !hasReplace {
		errs = append(errs, password_errors.ErrNoSpecialChars)
	}
	if !hasLower {
		errs = append(errs, password_errors.ErrNoLowercase)
	}
	if !hasUpper {
		errs = append(errs, password_errors.ErrNoUppercase)
	}
	if !hasDigits {
		errs = append(errs, password_errors.ErrNoDigits)
	}
	if entropy < minEntropy {
		errs = append(errs, password_errors.ErrLowEntropy)
	}

	if len(errs) == 0 {
		return nil
	}

	return &password_errors.PasswordError{
		BaseErr: password_errors.ErrPasswordValidator,
		Reasons: errs,
	}
}
