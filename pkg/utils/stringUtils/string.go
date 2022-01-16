package stringUtils

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// StringHasCorrectLength check that string size is between min and max
func StringHasCorrectLength(str string, min, max int) bool {
	return len(str) >= min && len(str) <= max
}

// ToSnakeCase convert string as snake_case string
// Example : NomUtilisateur --> nom_utilisateur (useful to convert golang fields as SQL columns)
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// RemoveDiacritics returns string without any diacritics
func RemoveDiacritics(str string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, str)
	return result
}

// CutWithPrefixAndSuffix will remove from string all char before prefix (including prefix) and all char after suffix (including suffix)
func CutWithPrefixAndSuffix(str, prefix, suffix string) (result string) {
	result = str
	prefixIndex := strings.Index(result, prefix)
	if prefixIndex > 0 {
		result = result[prefixIndex:]
	}
	result = strings.Replace(result, prefix, "", 1)
	suffixIndex := strings.Index(result, suffix)
	if suffixIndex > 0 {
		result = result[:suffixIndex]
	}
	result = strings.Replace(result, suffix, "", 1)
	return
}