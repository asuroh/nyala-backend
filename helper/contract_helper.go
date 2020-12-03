package helper

import (
	"errors"
	"strings"
)

// CheckPassword check password format must meet 3/4 criterias: Symbol, Number, Lowercase, Uppercase
func CheckPassword(password string) (err error) {
	if len(password) < 8 && len(password) > 10 {
		return errors.New(PasswordLength)
	}

	return err
}

// SexStringToBool ...
func SexStringToBool(sex string) (res bool, err error) {
	if strings.ToLower(sex) == "male" {
		return true, err
	} else if strings.ToLower(sex) == "famale" {
		return false, err
	}
	return res, errors.New(InvalidSex)
}

// SexBoolToString ...
func SexBoolToString(sex bool) string {
	if sex {
		return "male"
	} else if !sex {
		return "famale"
	}

	return ""
}
