package errs

import "errors"

var (
	ErrEmailFound = errors.New("email sudah terdaftar")
)
