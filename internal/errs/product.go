package errs

import "errors"

var (
	ErrNoProductFound        = errors.New("produk tidak ditemukan")
	ErrNoProductVariantFound = errors.New("varian produk tidak ditemukan")
	ErrNoProductImagesFound  = errors.New("images produk tidak ditemukan")
)
