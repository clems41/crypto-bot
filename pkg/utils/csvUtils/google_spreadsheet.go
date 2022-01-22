package csvUtils

import (
	"crypto-bot/pkg/utils/httpUtils"
	"fmt"
	"net/http"
)

const (
	updateApiBase             = "https://sheets.googleapis.com/v4/spreadsheets"
	valueInputOptionQueryName = "valueInputOption"
	valueInputOptionRaw       = "RAW"
)

func UpdateSpreadsheet(spreadsheetID string, rangeValue string, values [][]string) (err error) {
	url := fmt.Sprintf("%s/%s/values/%s", updateApiBase, spreadsheetID, rangeValue)
	queryParameters := map[string]string{
		valueInputOptionQueryName: valueInputOptionRaw,
	}
	type bodyRequest struct {
		Values [][]string `json:"values"`
	}
	body := bodyRequest{Values: values}
	request := httpUtils.Request{
		Url:               url,
		Method:            http.MethodPut,
		ContentType:       "application/json",
		QueryParameters:   queryParameters,
		HeadersParameters: nil,
		Body:              body,
	}

	statusCode, err := httpUtils.SendRequest(request, nil, nil)
	if err != nil {
		return
	}
	if statusCode > http.StatusCreated {
		return fmt.Errorf("google API return bad code %d", statusCode)
	}
	return
}
