package mathUtils

import (
	"fmt"
	"strconv"
)

func RemoveNDecimal(value float64, nbDecimal int) (result float64, err error) {
	fmtStr := fmt.Sprintf("%%0.%df", nbDecimal)
	valueStrWithoutDecimal := fmt.Sprintf(fmtStr, value)
	result, err = strconv.ParseFloat(valueStrWithoutDecimal, 64)
	return
}
