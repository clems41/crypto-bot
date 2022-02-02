package tradingConst

type Pair string

/* Pairs using EURO */
const (
	BitcoinEuro  Pair = "BTC_EUR"
	DashEuro     Pair = "DASH_EUR"
	EthereumEuro Pair = "ETH_EUR"
	CardanoEuro  Pair = "ADA_EUR"
)

/* Pairs using only crypto */

const (
	EthereumBitcoin Pair = "ETH_BTC"
)

/* Pairs using USD */

const (
	BitcoinUSDollar Pair = "BTC_USD"
)
