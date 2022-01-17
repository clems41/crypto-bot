package tradingPlatform

type Api interface {
	GetWalletBalance() (amount int, err error)
}
