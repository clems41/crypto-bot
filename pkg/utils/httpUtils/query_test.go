package httpUtils

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/require"
	"isi.nc/go/starter-module/utils/httpUtils"
	"isi.nc/go/starter-module/utils/searchUtils"
	"net/http"
	"sort"
	"strings"
	"testing"
)

func TestUnmarshallQueryParameters(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	queryParameters := map[string]string{
		"field1":  "test1",
		"field2":  "test2",
		"field33": "test3",
	}

	type testStruct struct {
		Field1 string `json:"field1"`  // check with same field name and json tag, should be OK
		Field2 string `json:"field2"`  // check with same field name and json tag, should be OK
		Field3 string `json:"field33"` // check with different field name and json tag, should be OK
		Field4 string `json:"field4"`  // check with field not in query parameters, should be OK but empty
		Field5 string `json:"-"`       // check without json tag, should be OK but empty
	}

	query := httpRequest.URL.Query()
	for queryKey, queryValue := range queryParameters {
		query.Add(queryKey, queryValue)
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	require.Equal(t, queryParameters["field1"], test.Field1)
	require.Equal(t, queryParameters["field2"], test.Field2)
	require.Equal(t, queryParameters["field33"], test.Field3)
	require.Empty(t, test.Field4)
	require.Empty(t, test.Field5)
}

func TestUnmarshallQueryParameters_WithDefaultValues(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	queryParameters := map[string]string{
		"field1":  "test1",
		"field2":  "test2",
		"field33": "test3",
	}

	defaultValues := map[string]string{
		"field4": "toto",
		"field5": "20",
	}

	type testStruct struct {
		Field1 string `json:"field1"`  // check with same field name and json tag, should be OK
		Field2 string `json:"field2"`  // check with same field name and json tag, should be OK
		Field3 string `json:"field33"` // check with different field name and json tag, should be OK
		Field4 string `json:"field4"`  // check with field not in query parameters, should be OK but empty
		Field5 int    `json:"field5"`  // check without json tag, should be OK but empty
	}

	query := httpRequest.URL.Query()
	for queryKey, queryValue := range queryParameters {
		query.Add(queryKey, queryValue)
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, defaultValues)
	require.NoError(t, err)
	require.Equal(t, queryParameters["field1"], test.Field1)
	require.Equal(t, queryParameters["field2"], test.Field2)
	require.Equal(t, queryParameters["field33"], test.Field3)
	require.Equal(t, defaultValues["field4"], test.Field4)
	require.Equal(t, 20, test.Field5)
}

func TestUnmarshallQueryParameters_Integer(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	expectedInt := 12

	queryParameters := map[string]string{
		"field1": fmt.Sprintf("%d", expectedInt),
	}

	type testStruct struct {
		Field1 int `json:"field1"`
	}

	query := httpRequest.URL.Query()
	for queryKey, queryValue := range queryParameters {
		query.Add(queryKey, queryValue)
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	require.Equal(t, expectedInt, test.Field1)
}

func TestUnmarshallQueryParameters_Integer16(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	expectedInt := int16(12)

	queryParameters := map[string]string{
		"field1": fmt.Sprintf("%d", expectedInt),
	}

	type testStruct struct {
		Field1 int16 `json:"field1"`
	}

	query := httpRequest.URL.Query()
	for queryKey, queryValue := range queryParameters {
		query.Add(queryKey, queryValue)
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	require.Equal(t, expectedInt, test.Field1)
}

func TestUnmarshallQueryParameters_Float32(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	expectedFloat := float32(12.58)

	queryParameters := map[string]string{
		"field1": fmt.Sprintf("%f", expectedFloat),
	}

	type testStruct struct {
		Field1 float32 `json:"field1"`
	}

	query := httpRequest.URL.Query()
	for queryKey, queryValue := range queryParameters {
		query.Add(queryKey, queryValue)
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	require.Equal(t, expectedFloat, test.Field1)
}

func TestUnmarshallQueryParameters_Float64(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	expectedFloat := float64(12.58)

	queryParameters := map[string]string{
		"field1": fmt.Sprintf("%f", expectedFloat),
	}

	type testStruct struct {
		Field1 float64 `json:"field1"`
	}

	query := httpRequest.URL.Query()
	for queryKey, queryValue := range queryParameters {
		query.Add(queryKey, queryValue)
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	require.Equal(t, expectedFloat, test.Field1)
}

func TestUnmarshallQueryParameters_BoolTrue(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	expectedBool := true

	queryParameters := map[string]string{
		"field1": fmt.Sprintf("%v", expectedBool),
	}

	type testStruct struct {
		Field1 bool `json:"field1"`
	}

	query := httpRequest.URL.Query()
	for queryKey, queryValue := range queryParameters {
		query.Add(queryKey, queryValue)
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	require.Equal(t, expectedBool, test.Field1)
}

func TestUnmarshallQueryParameters_BoolFalse(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	expectedBool := false

	queryParameters := map[string]string{
		"field1": fmt.Sprintf("%v", expectedBool),
	}

	type testStruct struct {
		Field1 bool `json:"field1"`
	}

	query := httpRequest.URL.Query()
	for queryKey, queryValue := range queryParameters {
		query.Add(queryKey, queryValue)
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	require.Equal(t, expectedBool, test.Field1)
}

func TestUnmarshallQueryParameters_StringArray(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	expectedStr := "toto"

	multipleQueryParameters := map[string][]string{
		"field1": {expectedStr},
	}

	type testStruct struct {
		Field1 []string `json:"field1"`
	}

	query := httpRequest.URL.Query()
	for queryKey, multipleQueryValue := range multipleQueryParameters {
		for _, queryValue := range multipleQueryValue {
			query.Add(queryKey, queryValue)
		}
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	require.Len(t, test.Field1, 1)
	require.Equal(t, expectedStr, test.Field1[0])
}

func TestUnmarshallQueryParameters_StringArrayWithMultipleValues(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	expectedStringArray := []string{"toto", "tata", "titi"}

	multipleQueryParameters := map[string][]string{
		"field1": expectedStringArray,
	}

	type testStruct struct {
		Field1 []string `json:"field1"`
	}

	query := httpRequest.URL.Query()
	for queryKey, multipleQueryValue := range multipleQueryParameters {
		for _, queryValue := range multipleQueryValue {
			query.Add(queryKey, queryValue)
		}
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	// First order slices to make comparison easier
	sort.Slice(test.Field1, func(i, j int) bool {
		return test.Field1[i] < test.Field1[j]
	})
	sort.Slice(expectedStringArray, func(i, j int) bool {
		return expectedStringArray[i] < expectedStringArray[j]
	})
	require.Equal(t, expectedStringArray, test.Field1)
}

func TestUnmarshallQueryParameters_IntegerArray(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	expectedInt := 12

	multipleQueryParameters := map[string][]string{
		"field1": {fmt.Sprintf("%d", expectedInt)},
	}

	type testStruct struct {
		Field1 []int `json:"field1"`
	}

	query := httpRequest.URL.Query()
	for queryKey, multipleQueryValue := range multipleQueryParameters {
		for _, queryValue := range multipleQueryValue {
			query.Add(queryKey, queryValue)
		}
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	require.Len(t, test.Field1, 1)
	require.Equal(t, expectedInt, test.Field1[0])
}

func TestUnmarshallQueryParameters_IntegerArrayWithMultipleValues(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	expectedArray := []int{2, 3, 8}

	multipleQueryParameters := map[string][]int{
		"field1": expectedArray,
	}

	type testStruct struct {
		Field1 []int `json:"field1"`
	}

	query := httpRequest.URL.Query()
	for queryKey, multipleQueryValue := range multipleQueryParameters {
		for _, queryValue := range multipleQueryValue {
			query.Add(queryKey, fmt.Sprintf("%d", queryValue))
		}
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	sort.Slice(test.Field1, func(i, j int) bool {
		return test.Field1[i] < test.Field1[j]
	})
	sort.Slice(expectedArray, func(i, j int) bool {
		return expectedArray[i] < expectedArray[j]
	})
	require.Equal(t, expectedArray, test.Field1)
}

func TestUnmarshallQueryParameters_MissingQueryParameter(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	queryParameters := map[string]string{
		"field1": "toto",
	}

	type testStruct struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	query := httpRequest.URL.Query()
	for queryKey, queryValue := range queryParameters {
		query.Add(queryKey, queryValue)
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	require.Equal(t, queryParameters["field1"], test.Field1)
	require.Equal(t, 0, test.Field2)
}

func TestUnmarshallQueryParameters_MissingQueryParameterSlice(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	queryParameters := map[string]string{
		"field1": "toto",
	}

	type testStruct struct {
		Field1 string `json:"field1"`
		Field2 []int  `json:"field2"`
	}

	query := httpRequest.URL.Query()
	for queryKey, queryValue := range queryParameters {
		query.Add(queryKey, queryValue)
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	require.Equal(t, queryParameters["field1"], test.Field1)
	require.Len(t, test.Field2, 0)
}

func TestUnmarshallQueryParameters_WithStruct(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	queryParameters := map[string]string{
		"page": "2",
	}

	type testStruct struct {
		searchUtils.PaginationRequest
	}

	defaultValues := map[string]string{
		"pageSize": "20",
	}

	query := httpRequest.URL.Query()
	for queryKey, queryValue := range queryParameters {
		query.Add(queryKey, queryValue)
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, defaultValues)
	require.NoError(t, err)
	require.Equal(t, 2, test.Page)
	require.Equal(t, 20, test.PageSize)
}

func TestUnmarshallQueryParameters_WithStructAndSlice(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	queryParameters := map[string][]string{
		"sort": {"id:asc", "lastName:desc"},
	}

	type testStruct struct {
		searchUtils.SortingRequest
	}

	query := httpRequest.URL.Query()
	for queryKey, queryValues := range queryParameters {
		for _, queryValue := range queryValues {
			query.Add(queryKey, queryValue)
		}
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, nil)
	require.NoError(t, err)
	require.Equal(t, queryParameters["sort"], test.Sort)
}

func TestUnmarshallQueryParameters_WithStructAndDefault(t *testing.T) {
	httpRequest, err := http.NewRequest("method", "url", nil)
	require.NoError(t, err)

	queryParameters := map[string][]string{}

	type testStruct struct {
		searchUtils.SortingRequest
	}

	expected := []string{"id:asc", "lastName:desc"}

	defaultValues := map[string]string{
		"sort": strings.Join(expected, ","),
	}

	query := httpRequest.URL.Query()
	for queryKey, queryValues := range queryParameters {
		for _, queryValue := range queryValues {
			query.Add(queryKey, queryValue)
		}
	}
	httpRequest.URL.RawQuery = query.Encode()

	var test testStruct
	err = httpUtils.UnmarshallQueryParameters(restful.NewRequest(httpRequest), &test, defaultValues)
	require.NoError(t, err)
	require.Equal(t, expected, test.Sort)
}
