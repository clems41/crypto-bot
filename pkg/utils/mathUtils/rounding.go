package mathUtils

import (
	"fmt"
	"math"
	"strconv"
)

// RoundingAfterNDecimal will return value but without exceeding decimals using classic rounding.
// Ex : 1.656 and nb = 2 will return 1.66
func RoundingAfterNDecimal(value float64, nbDecimal int) (result float64, err error) {
	decimalCoefficient := math.Pow(10, float64(nbDecimal))
	valueStrWithoutDecimal := fmt.Sprint(math.Round(value*decimalCoefficient) / decimalCoefficient)
	result, err = strconv.ParseFloat(valueStrWithoutDecimal, 64)
	return
}

// FloorRoundingAfterNDecimal will return value but without exceeding decimals using floor rounding.
// Ex : 1.656 and nb = 2 will return 1.65
func FloorRoundingAfterNDecimal(value float64, nbDecimal int) (result float64, err error) {
	decimalCoefficient := math.Pow(10, float64(nbDecimal))
	valueStrWithoutDecimal := fmt.Sprint(math.Floor(value*decimalCoefficient) / decimalCoefficient)
	result, err = strconv.ParseFloat(valueStrWithoutDecimal, 64)
	return
}

// CeilRoundingAfterNDecimal will return value but without exceeding decimals using ceil rounding.
// Ex : 1.652 and nb = 2 will return 1.66
func CeilRoundingAfterNDecimal(value float64, nbDecimal int) (result float64, err error) {
	decimalCoefficient := math.Pow(10, float64(nbDecimal))
	valueStrWithoutDecimal := fmt.Sprint(math.Ceil(value*decimalCoefficient) / decimalCoefficient)
	result, err = strconv.ParseFloat(valueStrWithoutDecimal, 64)
	return
}
