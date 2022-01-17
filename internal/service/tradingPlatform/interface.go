package tradingPlatform

type Api interface {
	GetWalletBalance() (view WalletView, err error)
	OpenPosition(form OpenPositionForm) (view OpenPositionView, err error)
	ClosePosition(form ClosePositionForm) (view ClosePositionView, err error)
	GetPrice(form GetPriceForm) (view GetPriceView, err error)
}
