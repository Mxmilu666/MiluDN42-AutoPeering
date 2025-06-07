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

// GetAllNodesInfo 获取所有 node 的详细信息
func GetAllNodesInfo(c *gin.Context) {
	nodes := source.AppConfig.Nodes
	results := make([]map[string]interface{}, len(nodes))

	type nodeResult struct {
		idx  int
		item map[string]interface{}
	}
	ch := make(chan nodeResult, len(nodes))
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
				info, err := api.GetNodeInfoByName(node.Name)
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
				item["info"] = info
			}

			ch <- nodeResult{idx: idx, item: item}
		}(i, n)
	}

	wg.Wait()
	close(ch)

	for res := range ch {
		results[res.idx] = res.item
	}

	SendResponse(c, http.StatusOK, "success", results)
}

// GetNodeInfo 获取指定 node 的详细信息
func GetNodeInfo(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		SendResponse(c, http.StatusBadRequest, "node name is required", nil)
		return
	}

	info, err := api.GetNodeInfoByName(name)
	if err != nil {
		SendResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	SendResponse(c, http.StatusOK, "success", info)
}
