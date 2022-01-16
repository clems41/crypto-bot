package fakeData

import (
	"math/rand"
	"time"
)

func FakeTime() time.Time {
	rand.Seed(time.Now().UnixNano())
	min := time.Date(minYearForFakeTime, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(maxYearForFakeTime, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func FakeTimeBetween(min, max time.Time) time.Time {
	rand.Seed(time.Now().UnixNano())
	delta := max.Unix() - min.Unix()

	sec := rand.Int63n(delta) + min.Unix()
	return time.Unix(sec, 0)
}

func FakeTimeBefore(max time.Time) time.Time {
	rand.Seed(time.Now().UnixNano())
	min := time.Date(minYearForFakeTime, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max.Unix() - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func FakeTimeAfter(min time.Time) time.Time {
	rand.Seed(time.Now().UnixNano())
	max := time.Date(maxYearForFakeTime, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min.Unix()

	sec := rand.Int63n(delta) + min.Unix()
	return time.Unix(sec, 0)
}

func FakeIntBetween(min, max int64) int {
	rand.Seed(time.Now().UnixNano())
	return int(rand.Int63n(max - min) + min)
}
