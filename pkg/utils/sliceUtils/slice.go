package sliceUtils

// SliceContainsStr will check if string is contains in slice of string
func SliceContainsStr(slice []string, str string) bool {
	if slice == nil {
		return false
	}
	for _, elem := range slice {
		if elem == str {
			return true
		}
	}
	return false
}

// StrSlicesHaveAtLeastOneMatch will check if 2 slice of string has at least one identical element
func StrSlicesHaveAtLeastOneMatch(slice1, slice2 []string) bool {
	if slice1 == nil || slice2 == nil {
		return false
	}
	for _, elem1 := range slice1 {
		for _, elem2 := range slice2 {
			if elem1 == elem2 {
				return true
			}
		}
	}
	return false
}

// StrSlicesContainsAll will check if all elements from slice1 are included in slice2
func StrSlicesContainsAll(slice1, slice2 []string) bool {
	if slice1 == nil || slice2 == nil {
		return false
	}

	for _, elem1 := range slice1 {
		var found bool
		for _, elem2 := range slice2 {
			if elem1 == elem2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
