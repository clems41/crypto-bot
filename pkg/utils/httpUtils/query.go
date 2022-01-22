package httpUtils

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"reflect"
	"strconv"
	"strings"
)

// UnmarshallQueryParameters will loop overs fields from out interface one by one and fill them with value from query parameters based on json tag.
// Out interface should be pointer on struct.
// Out interface fields should be of following types : int (classic, 8, 16, 32 or 64), string, float (32 or 64), bool or array of previous types (ex: []string).
// Error will be returned if not.
// You can specify default values with map[string]string (can be nil), meaning that if there is no query parameter for specific field, defaultValue will be used if found in map.
// In this case, jsonTag name should be used as key of this map.
func UnmarshallQueryParameters(req *restful.Request, out interface{}, defaultValues map[string]string) (err error) {
	// Check that out interface is type of pointer on struct
	outValue := reflect.ValueOf(out)
	outKind := outValue.Kind()
	if outKind != reflect.Ptr {
		return fmt.Errorf("out interface should be ptr but it is %s", outKind.String())
	}
	// Get value of pointer
	outValue = outValue.Elem()
	outKind = outValue.Kind()
	if outKind != reflect.Struct {
		return fmt.Errorf("out interface should be ptr on struct but it is ptr on %s", outKind.String())
	}

	if defaultValues == nil {
		defaultValues = make(map[string]string)
	}

	// Loop over all struct fields with recursive diving until we got a usable field
	err = loopOverAllFieldsRecursively(req, &outValue, defaultValues)
	return
}

func loopOverAllFieldsRecursively(req *restful.Request, outValue *reflect.Value, defaultValues map[string]string) (err error) {
	if err != nil {
		return
	}
	for outValueFieldIdx := 0; outValueFieldIdx < outValue.NumField(); outValueFieldIdx++ {
		outValueFieldValue := outValue.Field(outValueFieldIdx)
		outValueFieldKind := outValueFieldValue.Kind()
		outValueFieldTag := outValue.Type().Field(outValueFieldIdx).Tag.Get(fieldTag)
		switch outValueFieldKind {
		case reflect.Struct:
			err = loopOverAllFieldsRecursively(req, &outValueFieldValue, defaultValues)
			if err != nil {
				return
			}
		case reflect.Slice:
			valuesStr := req.QueryParameters(outValueFieldTag)
			defaultValue := defaultValues[outValueFieldTag]
			if len(valuesStr) == 0 && defaultValue != "" {
				valuesStr = strings.Split(defaultValue, ",")
			}
			if len(valuesStr) != 0 {
				err = setFieldValues(valuesStr, &outValueFieldValue)
				if err != nil {
					return
				}
			}
		case reflect.String:
			valueStr := req.QueryParameter(outValueFieldTag)
			defaultValue := defaultValues[outValueFieldTag]
			if valueStr == "" {
				valueStr = defaultValue
			}
			outValueFieldValue.Set(reflect.ValueOf(valueStr))
		case reflect.Int:
			valueStr := req.QueryParameter(outValueFieldTag)
			defaultValue := defaultValues[outValueFieldTag]
			if valueStr == "" {
				valueStr = defaultValue
			}
			if valueStr != "" {
				var value int
				value, err = strconv.Atoi(valueStr)
				if err != nil {
					return
				}
				outValueFieldValue.Set(reflect.ValueOf(value))
			}
		case reflect.Int8:
			valueStr := req.QueryParameter(outValueFieldTag)
			defaultValue := defaultValues[outValueFieldTag]
			if valueStr == "" {
				valueStr = defaultValue
			}
			if valueStr != "" {
				var value int
				value, err = strconv.Atoi(valueStr)
				if err != nil {
					return
				}
				outValueFieldValue.Set(reflect.ValueOf(int8(value)))
			}
		case reflect.Int16:
			valueStr := req.QueryParameter(outValueFieldTag)
			defaultValue := defaultValues[outValueFieldTag]
			if valueStr == "" {
				valueStr = defaultValue
			}
			if valueStr != "" {
				var value int
				value, err = strconv.Atoi(valueStr)
				if err != nil {
					return
				}
				outValueFieldValue.Set(reflect.ValueOf(int16(value)))
			}
		case reflect.Int32:
			valueStr := req.QueryParameter(outValueFieldTag)
			defaultValue := defaultValues[outValueFieldTag]
			if valueStr == "" {
				valueStr = defaultValue
			}
			if valueStr != "" {
				var value int
				value, err = strconv.Atoi(valueStr)
				if err != nil {
					return
				}
				outValueFieldValue.Set(reflect.ValueOf(int32(value)))
			}
		case reflect.Int64:
			valueStr := req.QueryParameter(outValueFieldTag)
			defaultValue := defaultValues[outValueFieldTag]
			if valueStr == "" {
				valueStr = defaultValue
			}
			if valueStr != "" {
				var value int
				value, err = strconv.Atoi(valueStr)
				if err != nil {
					return
				}
				outValueFieldValue.Set(reflect.ValueOf(int64(value)))
			}
		case reflect.Bool:
			valueStr := req.QueryParameter(outValueFieldTag)
			defaultValue := defaultValues[outValueFieldTag]
			if valueStr == "" {
				valueStr = defaultValue
			}
			if valueStr != "" {
				var value bool
				value, err = strconv.ParseBool(valueStr)
				if err != nil {
					return
				}
				outValueFieldValue.Set(reflect.ValueOf(value))
			}
		case reflect.Float32:
			valueStr := req.QueryParameter(outValueFieldTag)
			defaultValue := defaultValues[outValueFieldTag]
			if valueStr == "" {
				valueStr = defaultValue
			}
			if valueStr != "" {
				var value float64
				value, err = strconv.ParseFloat(valueStr, 64)
				if err != nil {
					return
				}
				outValueFieldValue.Set(reflect.ValueOf(float32(value)))
			}
		case reflect.Float64:
			valueStr := req.QueryParameter(outValueFieldTag)
			defaultValue := defaultValues[outValueFieldTag]
			if valueStr == "" {
				valueStr = defaultValue
			}
			if valueStr != "" {
				var value float64
				value, err = strconv.ParseFloat(valueStr, 64)
				if err != nil {
					return
				}
				outValueFieldValue.Set(reflect.ValueOf(value))
			}
		}
	}
	return
}

func setFieldValues(valuesStr []string, fieldValue *reflect.Value) (err error) {
	// Check if field is slice
	fieldKind := fieldValue.Kind()
	if fieldKind != reflect.Slice {
		return fmt.Errorf("fieldValue should be type of slice but it is %s", fieldKind.String())
	}
	// Check if field can be set
	if !fieldValue.CanSet() {
		return fmt.Errorf("out interface field %v cannot be set", fieldValue)
	}
	// Get slice element type
	sliceElementKind := reflect.TypeOf(fieldValue.Interface()).Elem().Kind()
	if sliceElementKind == reflect.String {
		var valueSlice []string
		for _, valueStr := range valuesStr {
			valueSlice = append(valueSlice, valueStr)
		}
		fieldValue.Set(reflect.ValueOf(valueSlice))
	} else if sliceElementKind == reflect.Int {
		var valueSlice []int
		for _, valueStr := range valuesStr {
			var value int
			value, err = strconv.Atoi(valueStr)
			if err != nil {
				return
			}
			valueSlice = append(valueSlice, value)
		}
		// Set field only if kind are equals
		if fieldKind == reflect.TypeOf(valueSlice).Kind() {
			fieldValue.Set(reflect.ValueOf(valueSlice))
		}
	} else if sliceElementKind == reflect.Int8 {
		var valueSlice []int8
		for _, valueStr := range valuesStr {
			var value int
			value, err = strconv.Atoi(valueStr)
			if err != nil {
				return
			}
			valueSlice = append(valueSlice, int8(value))
		}
		fieldValue.Set(reflect.ValueOf(valueSlice))
	} else if sliceElementKind == reflect.Int16 {
		var valueSlice []int16
		for _, valueStr := range valuesStr {
			var value int
			value, err = strconv.Atoi(valueStr)
			if err != nil {
				return
			}
			valueSlice = append(valueSlice, int16(value))
		}
		fieldValue.Set(reflect.ValueOf(valueSlice))
	} else if sliceElementKind == reflect.Int32 {
		var valueSlice []int32
		for _, valueStr := range valuesStr {
			var value int
			value, err = strconv.Atoi(valueStr)
			if err != nil {
				return
			}
			valueSlice = append(valueSlice, int32(value))
		}
		fieldValue.Set(reflect.ValueOf(valueSlice))
	} else if sliceElementKind == reflect.Int64 {
		var valueSlice []int64
		for _, valueStr := range valuesStr {
			var value int
			value, err = strconv.Atoi(valueStr)
			if err != nil {
				return
			}
			valueSlice = append(valueSlice, int64(value))
		}
		// Set field only if kind are equals
		if fieldKind == reflect.TypeOf(valueSlice).Kind() {
			fieldValue.Set(reflect.ValueOf(valueSlice))
		}
	} else if sliceElementKind == reflect.Bool {
		var valueSlice []bool
		for _, valueStr := range valuesStr {
			var value bool
			value, err = strconv.ParseBool(valueStr)
			if err != nil {
				return
			}
			valueSlice = append(valueSlice, value)
		}
		// Set field only if kind are equals
		if fieldKind == reflect.TypeOf(valueSlice).Kind() {
			fieldValue.Set(reflect.ValueOf(valueSlice))
		}
	} else if sliceElementKind == reflect.Float32 {
		var valueSlice []float32
		for _, valueStr := range valuesStr {
			var value float64
			value, err = strconv.ParseFloat(valueStr, 32)
			if err != nil {
				return
			}
			valueSlice = append(valueSlice, float32(value))
		}
		fieldValue.Set(reflect.ValueOf(valueSlice))
	} else if sliceElementKind == reflect.Float64 {
		var valueSlice []float64
		for _, valueStr := range valuesStr {
			var value float64
			value, err = strconv.ParseFloat(valueStr, 32)
			if err != nil {
				return
			}
			valueSlice = append(valueSlice, value)
		}
		// Set field only if kind are equals
		if fieldKind == reflect.TypeOf(valueSlice).Kind() {
			fieldValue.Set(reflect.ValueOf(valueSlice))
		}
	}
	return
}
