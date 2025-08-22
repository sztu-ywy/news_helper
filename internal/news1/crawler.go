package news1

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/uozi-tech/cosy/logger"
)

func CrawlHtml(method, url string, headers http.Header, queryParams url.Values, data interface{}) ([]byte, error) {
	var reqBody io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(jsonData)
	}
	// 1. 发起 HTTP 请求
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	req.URL.RawQuery = queryParams.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf("请求失败:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("HTTP 错误: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
