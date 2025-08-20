package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/uozi-tech/cosy/logger"
)

const (
	endPoint = ""
)


func Post[T any](url string, headers http.Header, queryParams url.Values, data interface{}) (*T, error) {
	if len(queryParams) > 0 {
		url += "?" + queryParams.Encode()
	}
	var jsonData []byte
	if data != nil {
		var err error
		jsonData, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest("POST", endPoint+url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	headers.Add("Content-Type", "application/json")
	req.Header = headers
	return doRequest[T](req)
}
func Get[T any](url string, headers http.Header, queryParams url.Values, data interface{}) (*T, error) {
	if len(queryParams) > 0 {
		url += "?" + queryParams.Encode()
	}
	var jsonData []byte
	if data != nil {
		var err error
		jsonData, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest("GET", endPoint+url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header = headers

	return doRequest[T](req)
}

func doRequest[T any](req *http.Request) (*T, error) {
	logger.Debug(req.Method, req.URL)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		logger.Error("doRequest: ", resp.StatusCode, string(body))
		var errRes map[string]interface{}
		err := json.Unmarshal(body, &errRes)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		if errRes["msg"] != nil {
			return nil, errors.New(errRes["msg"].(string))
		}
		if errRes["message"] != nil {
			return nil, errors.New(errRes["message"].(string))
		}
		if errRes["code"] != nil {
			return nil, errors.New(errRes["code"].(string))
		}
		return nil, fmt.Errorf("get request failed with status code: %d, error: %v", resp.StatusCode, errRes)
	}
	if string(body) == "" {
		return nil, nil
	}
	logger.Debug(resp.StatusCode)
	var res T
	if err := json.Unmarshal(body, &res); err != nil {
		logger.Error(err)
		return nil, err
	}
	return &res, nil
}
