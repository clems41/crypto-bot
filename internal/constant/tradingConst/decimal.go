package tradingConst

var (
	MaxDecimalByPair = map[Pair]int{
		BitcoinEuro:     1,
		DashEuro:        1,
		EthereumEuro:    1,
		EthereumBitcoin: 1,
		CardanoEuro:     5,
	}
	MaxDecimalByCurrency = map[Currency]int{
		Bitcoin:  6,
		Dash:     6,
		Euro:     2,
		Ethereum: 6,
		Cardano:  6,
	}
)
