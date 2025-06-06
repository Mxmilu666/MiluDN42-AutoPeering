package handles

import (
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/node/source/helper"
	"github.com/gin-gonic/gin"
)

// RestartHandler 重启指定 ASN 的 WireGuard 隧道和 BIRD 会话
func RestartHandler(c *gin.Context) {
	asn := c.Query("asn")
	if asn == "" {
		SendResponse(c, http.StatusBadRequest, "error", "asn is required")
		return
	}

	_, _, err := helper.RestartTunnelAndBird(asn)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, "error", err.Error())
		return
	}

	SendResponse(c, http.StatusOK, "success", nil)
}
