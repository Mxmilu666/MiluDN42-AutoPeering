package handles

import (
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/backend/source/helper"
	"github.com/gin-gonic/gin"
)

// GetPeerInfoHandler 查询指定 ASN 的 WireGuard 隧道信息
func GetPeerInfoHandler(c *gin.Context) {
	asn := c.Query("asn")
	if asn == "" {
		SendResponse(c, http.StatusBadRequest, "error", "asn is required")
		return
	}
	wginfo, err := helper.RawWGShow("as_" + asn)
	if err != nil {
		SendResponse(c, http.StatusNotFound, "error", "asn not found")
		return
	}

	info := map[string]string{
		"asn":    asn,
		"wginfo": wginfo,
	}
	SendResponse(c, http.StatusNotFound, "success", info)
}
