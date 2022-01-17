package tradingPlatform

type Api interface {
	GetWalletBalance() (amount float32, err error)
}
