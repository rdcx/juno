package util

import "errors"

func WrapErr(err error, msg string) error {
	return errors.New(err.Error() + ": " + msg)
}
