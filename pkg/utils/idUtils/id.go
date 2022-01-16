package idUtils

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

// ContainsID return true if value is in slice
func ContainsID(slice []uint, value uint) bool {
	for _, elem := range slice {
		if elem == value {
			return true
		}
	}
	return false
}

//StrToID convert string to uint
func StrToID(str string) (id uint, err error) {
	integer, err := strconv.Atoi(str)
	if err != nil {
		return id, errors.WithStack(err)
	}
	id = uint(integer)
	return
}

//IDToStr convert uint to string
func IDToStr(id uint) string {
	return fmt.Sprintf("%d", id)
}
