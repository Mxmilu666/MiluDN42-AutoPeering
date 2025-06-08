package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source"
)

// RequestNodeByName 支持泛型，自动解析json，支持GET/POST
func RequestNodeByName[T any](nodeName, apiPath, method string, bodyData interface{}) (*T, error) {
	var nodeAddr, token string
	for _, node := range source.AppConfig.Nodes {
		if node.Name == nodeName {
			nodeAddr = node.Address
			token = node.Token
			break
		}
	}
	if nodeAddr == "" || token == "" {
		return nil, fmt.Errorf("node not found: %s", nodeName)
	}
	url := nodeAddr + apiPath

	var reqBody io.Reader
	if method == "POST" && bodyData != nil {
		jsonBytes, err := json.Marshal(bodyData)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		// 兼容统一API格式，优先返回data字段（如果是字符串且非空），否则返回msg
		var raw map[string]interface{}
		_ = json.Unmarshal(body, &raw)
		if data, ok := raw["data"].(string); ok && data != "" {
			return nil, fmt.Errorf("%s", data)
		}
		if msg, ok := raw["msg"].(string); ok {
			return nil, fmt.Errorf("%s", msg)
		}
		return nil, fmt.Errorf("request failed: %s", string(body))
	}
	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
