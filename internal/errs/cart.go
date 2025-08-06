package errs

import "errors"

var (
	ErrQuantityLessThanZero     = errors.New("kuantitas tidak boleh kurang dari 1")
	ErrProductNotFound          = errors.New("produk tidak ditemukan")
	ErrQuantityLessThanMinOrder = errors.New("kuantitas kurang dari minimal pembelian")
	ErrNotEnoughStock           = errors.New("kuantitas melebihi stok")
	ErrCartsEmpty               = errors.New("tidak ada produk di cart")
	ErrNoCartProduct            = errors.New("produk tidak ditemukan di cart")
)
