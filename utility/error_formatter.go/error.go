package error_formatter

import "errors"

func NewError(message string) error{
	return errors.New("[ERROR] "+message)
}