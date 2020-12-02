package helper

import (
	"errors"
)

// CheckPassword check password format must meet 3/4 criterias: Symbol, Number, Lowercase, Uppercase
func CheckPassword(password string) (err error) {
	if len(password) < 8 && len(password) > 10 {
		return errors.New(PasswordLength)
	}

	return err
}
