package errs

import "errors"

var (
	ErrNoRatingFound = errors.New("rating tidak ditemukan")
	ErrAlreadyRated  = errors.New("produk sudah dinilai")
)
