package helper

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source"
)

// RequestNodeByName 通过 node 名称自动获取 token 并请求 node
func RequestNodeByName(nodeName, apiPath string) ([]byte, error) {
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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
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
		return nil, fmt.Errorf("request failed: %s", string(body))
	}
	return body, nil
}
