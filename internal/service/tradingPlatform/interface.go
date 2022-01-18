package tradingPlatform

type Api interface {
	GetOpenedPositions() (view GetOpenedPositionsView, err error)
	GetAllPositions() (view GetAllPositionsView, err error)
	GetWalletBalance() (view WalletView, err error)
	OpenPosition(form OpenPositionForm) (view OpenPositionView, err error)
	ClosePosition(form ClosePositionForm) (view ClosePositionView, err error)
	GetPrice(form GetPriceForm) (view GetPriceView, err error)
}
