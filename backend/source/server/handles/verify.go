package handles

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type VerifyRequest struct {
	IPs []string `json:"ips" binding:"required"`
}

type VerifyInfo struct {
	Token   string
	Dir     string
	IPAllow map[string]struct{}
	Expire  time.Time
}

var (
	verifyStore = make(map[string]*VerifyInfo)
	verifyMu    sync.Mutex
)

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// POST /api/verify/request
func RequestVerify(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IPs) == 0 {
		SendResponse(c, http.StatusBadRequest, "invalid request", nil)
		return
	}

	// 检查是否有IP已存在于其他请求中
	verifyMu.Lock()
	for _, info := range verifyStore {
		for _, ip := range req.IPs {
			if _, exists := info.IPAllow[ip]; exists && time.Now().Before(info.Expire) {
				verifyMu.Unlock()
				SendResponse(c, http.StatusConflict, "ip already in use", nil)
				return
			}
		}
	}
	verifyMu.Unlock()

	dir := randomString(8) // 只生成目录名，不带/verify/
	token := randomString(32)
	ipMap := make(map[string]struct{})
	for _, ip := range req.IPs {
		ipMap[ip] = struct{}{}
	}
	info := &VerifyInfo{
		Token:   token,
		Dir:     dir,
		IPAllow: ipMap,
		Expire:  time.Now().Add(5 * time.Minute), // 有效期5分钟
	}
	verifyMu.Lock()
	verifyStore[dir] = info // 只用dir作为key
	verifyMu.Unlock()
	SendResponse(c, http.StatusOK, "success", gin.H{"dir": dir, "token": token})
}

// GET /verify/:dir
func VerifyHandler(c *gin.Context) {
	dir := c.Param("dir")
	verifyMu.Lock()
	info, ok := verifyStore[dir] // 只用dir查找
	verifyMu.Unlock()
	if !ok || time.Now().After(info.Expire) {
		SendResponse(c, http.StatusNotFound, "not found or expired", nil)
		return
	}
	clientIP := c.ClientIP()
	if _, allowed := info.IPAllow[clientIP]; !allowed {
		SendResponse(c, http.StatusForbidden, "ip not allowed", nil)
		return
	}
	c.String(http.StatusOK, info.Token)
}
