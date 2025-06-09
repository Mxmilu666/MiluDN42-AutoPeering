package handles

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source"
	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/api"
	"github.com/gin-gonic/gin"
)

// GetPeerInfo 获取指定 ASN 的 Peer 信息
func GetPeerInfo(c *gin.Context) {
	asn := c.Query("asn")
	node := c.Query("node") // 支持通过 node 参数指定节点
	if asn == "" {
		SendResponse(c, http.StatusBadRequest, "asn is required", nil)
		return
	}
	if node == "" {
		SendResponse(c, http.StatusBadRequest, "node is required", nil)
		return
	}
	info, err := api.GetPeerInfoByASN(node, asn)
	if err != nil {
		SendResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	SendResponse(c, http.StatusOK, "success", info)
}

// GetAllPeersInfo 获取所有节点上指定 ASN 的 Peer 信息
func GetAllPeersInfo(c *gin.Context) {
	asn := c.Query("asn")
	if asn == "" {
		SendResponse(c, http.StatusBadRequest, "asn is required", nil)
		return
	}

	nodes := source.AppConfig.Nodes
	results := make([]map[string]interface{}, len(nodes))

	type peerResult struct {
		idx  int
		item map[string]interface{}
	}
	ch := make(chan peerResult, len(nodes))
	var wg sync.WaitGroup

	timeout := 10 * time.Second // 每个节点请求超时时间

	for i, n := range nodes {
		wg.Add(1)
		go func(idx int, node source.NodeConfig) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			infoChan := make(chan interface{})
			errChan := make(chan error)

			go func() {
				info, err := api.GetPeerInfoByASN(node.Name, asn)
				if err != nil {
					errChan <- err
					return
				}
				infoChan <- info
			}()

			item := map[string]interface{}{
				"name": node.Name,
			}

			select {
			case <-ctx.Done():
				item["error"] = "Node unreachable or timeout"
			case <-errChan:
				item["error"] = "Node internal error"
			case info := <-infoChan:
				item["peerinfo"] = info
			}

			ch <- peerResult{idx: idx, item: item}
		}(i, n)
	}

	wg.Wait()
	close(ch)

	for res := range ch {
		results[res.idx] = res.item
	}

	SendResponse(c, http.StatusOK, "success", results)
}
