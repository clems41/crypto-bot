package tradingUtils

func GetProfit(askPrice, bidPrice, amount float64) (profit float64) {
	profit = bidPrice * amount / askPrice
	return
}

func GetResult(askPrice, bidPrice, amount float64) (result float64) {
	profit := GetProfit(askPrice, bidPrice, amount)
	result = profit - amount
	return
}

func GetResultInPercent(amount, result float64) (percent float64) {
	if result != 0 {
		percent = result / amount * 100
	}
	return
}
