package errs

import "errors"

var (
	ErrCheckoutRequestZero      = errors.New("request produk kosong")
	ErrNoOrderFound             = errors.New("produk order tidak ditemukan")
	ErrNoCheckoutFound          = errors.New("produk checkout tidak ditemukan")
	ErrNoPaymentFound           = errors.New("jenis pembayaran tidak ditemukan")
	ErrNoOrderSubscriptionFound = errors.New("Order langganan tidak ditemukan")
	ErrNoSellerFound            = errors.New("Seller tidak ditemukan")
)
