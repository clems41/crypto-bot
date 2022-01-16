package fakeData

func FakeRange(min, max int64) []int {
	rangeValue := FakeIntBetween(min, max)
	return make([]int, rangeValue)
}
