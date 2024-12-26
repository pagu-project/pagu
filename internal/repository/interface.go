package repository

type IDatabase interface {
	IUser
	IVoucher
	IFaucet
	IZealy
	IValidator
	INotification
}
