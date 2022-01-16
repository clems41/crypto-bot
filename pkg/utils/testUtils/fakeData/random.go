package fakeData

import "reflect"

func RandomElement(slice interface{}) interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil
	}
	idx := FakeIntBetween(0, int64(s.Len() - 1))
	return s.Index(idx).Interface()
}
