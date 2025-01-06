package utils

import "errors"

func ValidateDepartment(name string) error {
	if len(name) < 4 || len(name) > 33 {
		return errors.New("name must be between 4 and 33 characters")
	}
	return nil
}
