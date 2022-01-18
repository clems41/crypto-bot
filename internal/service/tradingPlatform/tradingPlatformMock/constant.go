package tradingPlatformMock

const (
	datasetIDToUse       = 1
	datasetRelativePath1 = "internal/service/tradingPlatform/tradingPlatformMock/data/1_btc_14012022_1s.csv"
	datasetRelativePath2 = "internal/service/tradingPlatform/tradingPlatformMock/data/2_btc_17012022_1min.json"
	initWalletBalance    = float64(100.0)
	initDatasetIndex     = 0
	fakeFeesInPercent    = 0.15 // in percent, ex : if equal to 1 and price = 100 --> ask = 100.5 and bid = 99.5
)
