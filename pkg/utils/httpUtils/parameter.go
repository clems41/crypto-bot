package httpUtils

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"strconv"
	"strings"
)

// Parameter can be used to specify path and query parameter in swagger but can also be used to get value or default value
//  - Name : indicate name of query/path parameter that will be received in request
//  - Description : use in swagger to explain the purpose of the parameter
//  - DefaultValueStr : value that will be return if not received from request (should be string as long as all query parameters are string types).
//  If more than one query parameters are expected, DefaultValueStr can contain multiple default values separated by comma.
//  - DataType : use in swagger to specify data type
//  - Required : if true, will throw an error if not define when getting query/path parameter
type Parameter struct {
	Name            string
	Description     string
	DefaultValueStr string
	DataType        string
	Required        bool
}

func (param Parameter) PathParameter() *restful.Parameter {
	return restful.
		PathParameter(param.Name, param.Description).
		DataType(param.DataType).
		Required(param.Required).
		DefaultValue(param.DefaultValueStr)
}

func (param Parameter) QueryParameter() *restful.Parameter {
	return restful.
		QueryParameter(param.Name, param.Description).
		DataType(param.DataType).
		Required(param.Required).
		DefaultValue(param.DefaultValueStr)
}

func (param Parameter) Joker() string {
	return fmt.Sprintf("{%s}", param.Name)
}

func (param Parameter) GetValueFromQueryParameter(req *restful.Request) (value string, err error) {
	value = req.QueryParameter(param.Name)
	if value == "" {
		if param.Required {
			return value, fmt.Errorf("query parameter %s is not defined but it is required", param.Name)
		} else {
			value = param.DefaultValueStr
		}
	}
	return
}

func (param Parameter) GetValueFromQueryParameters(req *restful.Request) (values []string, err error) {
	values = req.QueryParameters(param.Name)
	if len(values) == 0 {
		if param.Required {
			return values, fmt.Errorf("query parameter %s is not defined but it is required", param.Name)
		} else {
			values = strings.Split(param.DefaultValueStr, ",")
		}
	}
	return
}

func (param Parameter) GetValueFromPathParameter(req *restful.Request) (value string, err error) {
	value = req.PathParameter(param.Name)
	if value == "" {
		if param.Required {
			return value, fmt.Errorf("path parameter %s is not defined but it is required", param.Name)
		} else {
			value = param.DefaultValueStr
		}
	}
	return
}

func (param Parameter) GetValueFromPathParameterAsUint(req *restful.Request) (value uint, err error) {
	valueStr := req.PathParameter(param.Name)
	if valueStr == "" {
		if param.Required {
			return value, fmt.Errorf("path parameter %s is not defined but it is required", param.Name)
		} else {
			valueStr = param.DefaultValueStr
		}
	}
	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return
	}
	value = uint(valueInt)
	return
}
