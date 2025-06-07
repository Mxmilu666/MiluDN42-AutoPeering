package handles

import (
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/helper"
	"github.com/gin-gonic/gin"
)

// SendASNVerifyCode 处理 ASN 邮箱获取
func SendASNVerifyCode(c *gin.Context) {
	asn := c.Query("asn")
	if asn == "" {
		SendResponse(c, http.StatusBadRequest, "asn is required", nil)
		return
	}

	emails, err := helper.GetASNMaintainerEmails(asn)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	SendResponse(c, http.StatusOK, "success", emails)
}
