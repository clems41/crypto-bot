package httpUtils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Url               string
	Method            string
	ContentType       string
	QueryParameters   map[string]string
	HeadersParameters map[string]string
	Body              interface{}
}

// SendRequest will send http request based on fields from Request.
// If body response is expected, you can specify your custom struct with bodyResponse parameter to unmarshall http response into this struct.
// You also can specify custom http.Client in order specify some parameters (timeout for example). If nil, default http.Client will be used.
//  - request : contains all information needed to send request through HTTP protocol
//  - bodyResponse (optional) : can be nil, if specified, response will unmarshall into this structure
//  - customClient (optional) : can be nil, if specified, this client will be used to send request
func SendRequest(request Request, bodyResponse interface{}, customClient *http.Client) (httpStatusCode int, err error) {
	// Create body reader if needed
	var bodyBytes []byte
	if request.Body != nil {
		bodyBytes, err = json.Marshal(request.Body)
		if err != nil {
			return
		}
	}
	bodyReader := bytes.NewBuffer(bodyBytes)

	// Prepare request
	var httpClient http.Client
	if customClient != nil {
		httpClient = *customClient
	}
	httpRequest, err := http.NewRequest(request.Method, request.Url, bodyReader)
	if err != nil {
		return
	}

	// Set content type if specified
	if request.ContentType != "" {
		httpRequest.Header.Set(contentTypeHeaderName, request.ContentType)
	}

	// Adding header parameters
	for headerName, headerValue := range request.HeadersParameters {
		httpRequest.Header.Set(headerName, headerValue)
	}

	// Adding query parameters
	queryParameters := httpRequest.URL.Query()
	for queryName, queryValue := range request.QueryParameters {
		queryParameters.Set(queryName, queryValue)
	}
	httpRequest.URL.RawQuery = queryParameters.Encode()

	// Send request and get response
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return
	}

	// Decode response
	httpStatusCode = httpResponse.StatusCode
	defer httpResponse.Body.Close()
	if bodyResponse != nil && httpStatusCode < http.StatusBadRequest {
		var bodyResponseBytes []byte
		bodyResponseBytes, err = ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			return
		}
		if len(bodyResponseBytes) > 0 && string(bodyResponseBytes) != "" {
			err = json.Unmarshal(bodyResponseBytes, bodyResponse)
		}
	}
	return
}
