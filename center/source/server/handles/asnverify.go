package handles

import (
	"net/http"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/helper"
	"github.com/gin-gonic/gin"
)

// SendASNVerifyCode 处理 ASN 邮箱获取并发送验证码
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

	if len(emails) == 0 {
		SendResponse(c, http.StatusNotFound, "no maintainer email found", nil)
		return
	}

	// 默认只给第一个邮箱发送验证码
	email := emails[0]
	err = helper.SendVerificationCodeByEmail(email, asn)
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	SendResponse(c, http.StatusOK, "verification code sent", map[string]interface{}{"email": email})
}

// 验证并下发 JWT
func VerifyASNCodeAndIssueJWT(c *gin.Context) {
	type Req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, http.StatusBadRequest, "invalid request", nil)
		return
	}

	ok, asn := helper.VerifyCode(req.Email, req.Code)
	if !ok {
		SendResponse(c, http.StatusUnauthorized, "invalid or expired code", nil)
		return
	}

	// 生成 JWT，payload 包含 asn
	claims := map[string]interface{}{
		"asn":   asn,
		"email": req.Email,
	}
	token, err := helper.JwtHelper.IssueToken(claims, "asn-verify", 3600) // 1小时有效
	if err != nil {
		SendResponse(c, http.StatusInternalServerError, "failed to issue token", nil)
		return
	}

	SendResponse(c, http.StatusOK, "success", map[string]interface{}{"token": token})
}
