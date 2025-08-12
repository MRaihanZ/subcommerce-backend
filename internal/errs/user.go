package errs

import "errors"

var (
	ErrNoDirectoryFound      = errors.New("direktori tidak ditemukan")
	ErrFailedToDeleteProfile = errors.New("gagal menghapus gambar profil")
	ErrUserNotFound          = errors.New("user tidak ditemukan")
)
