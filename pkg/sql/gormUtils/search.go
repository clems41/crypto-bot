package gormUtils

import (
	"cine-circle-api/pkg/utils/searchUtils"
	"cine-circle-api/pkg/utils/stringUtils"
	"fmt"
	"strings"
)

type SortField struct {
	Field string
	Asc   bool
}

type SortQuery struct {
	Sort []SortField
}

func (query SortQuery) OrderSQL() string {
	var orders []string
	for _, sortField := range query.Sort {
		var ascString string
		if sortField.Asc {
			ascString = "asc"
		} else {
			ascString = "desc"
		}
		orders = append(orders, fmt.Sprintf("%s %s", stringUtils.ToSnakeCase(sortField.Field), ascString))
	}
	return strings.Join(orders, ",")
}

type PaginationQuery struct {
	Page     int
	PageSize int
}

// Offset will add offset to SQL query depending on page and pageSize
func (query PaginationQuery) Offset() int {

	offset := 0

	if query.Page > 1 {
		offset = (query.Page - 1) * query.PageSize
	}

	return offset
}

func FromSortRequestToSortQuery(sortRequest searchUtils.SortingRequest) (sortQuery SortQuery) {
	for _, request := range sortRequest.Sort {
		split := strings.Split(request, ":")
		if len(split) == 2 {
			sortField := SortField{
				Field: split[0],
			}
			if split[1] == "desc" {
				sortField.Asc = false
			} else {
				sortField.Asc = true
			}
			sortQuery.Sort = append(sortQuery.Sort, sortField)
		} else if len(split) != 0 {
			sortField := SortField{
				Field: split[0],
				Asc:   true,
			}
			sortQuery.Sort = append(sortQuery.Sort, sortField)
		}
	}
	return
}
