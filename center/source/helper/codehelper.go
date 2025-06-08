package helper

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source"
	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/untils"

	"gopkg.in/gomail.v2"
)

var ErrVerificationCodeExists = errors.New("verification code already exists and has not expired")

// 验证码结构体
type VerificationCode struct {
	Code      string
	ExpiresAt time.Time
	Email     string
	ASN       string
}

// 内存缓存，用于存储验证码
var codeCache = struct {
	sync.RWMutex
	m map[string]VerificationCode
}{m: make(map[string]VerificationCode)}

type VerificationRequest struct {
	Email string `json:"email"`
}

func init() {
	// 启动定期清理过期验证码
	go cleanupExpiredCodes()
}

// cleanupExpiredCodes 定期清理过期的验证码和临时码
func cleanupExpiredCodes() {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟清理一次
	defer ticker.Stop()

	for range ticker.C {
		// 清理过期的验证码
		removeExpiredVerificationCodes()
	}
}

// removeExpiredVerificationCodes 清理过期的验证码
func removeExpiredVerificationCodes() {
	now := time.Now()

	codeCache.Lock()
	defer codeCache.Unlock()

	// 遍历所有验证码，删除已过期的
	for email, code := range codeCache.m {
		if now.After(code.ExpiresAt) {
			delete(codeCache.m, email)
		}
	}
}

func SendEmail(to, subject, body string) error {
	// 从配置中读取 SMTP 相关信息
	username := source.AppConfig.SMTP.Username
	from := source.AppConfig.SMTP.From
	password := source.AppConfig.SMTP.Password
	smtpHost := source.AppConfig.SMTP.Host
	smtpPort := source.AppConfig.SMTP.Port

	// 创建一封邮件
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, username, password)
	// 发送邮件
	return d.DialAndSend(m)
}

// SendVerificationCodeByEmail 发送验证码到用户的电子邮件
func SendVerificationCodeByEmail(to, asn string) error {
	// 检查是否存在未过期的验证码
	codeCache.RLock()
	existingCode, exists := codeCache.m[to]
	codeCache.RUnlock()

	// 如果已存在未过期的验证码，则不发送验证码防止被刷爆接口
	if exists && time.Now().Before(existingCode.ExpiresAt) {
		return fmt.Errorf("%w: please check your email or wait for expiration", ErrVerificationCodeExists)
	}

	// 生成验证码
	code, err := untils.GenerateRandomCode(6, true)
	if err != nil {
		return fmt.Errorf("failed to generate verification code: %w", err)
	}

	// 验证码有效期为10分钟
	expirationMinutes := 10

	// 将验证码存储在缓存中
	expiration := time.Now().Add(time.Duration(expirationMinutes) * time.Minute)
	codeCache.Lock()
	codeCache.m[to] = VerificationCode{
		Code:      code,
		ExpiresAt: expiration,
		Email:     to,
		ASN:       asn,
	}
	codeCache.Unlock()

	subject := "[MiluDN42-AutoPeering] Your ASN verification code is here!"
	body := fmt.Sprintf("Your ASN verification code is: %s. It is valid for %d minutes.",
		code,
		expirationMinutes)

	if err := SendEmail(to, subject, body); err != nil {
		return fmt.Errorf("failed to send verification code: %w", err)
	}

	return nil
}

// VerifyCode 验证用户输入的验证码是否正确
func VerifyCode(email, code string) (bool, string) {
	codeCache.RLock()
	storedCode, exists := codeCache.m[email]
	codeCache.RUnlock()

	if !exists {
		return false, ""
	}

	// 检查验证码是否是对应邮箱
	if storedCode.Email != email {
		return false, ""
	}

	// 检查验证码是否已过期
	if time.Now().After(storedCode.ExpiresAt) {
		return false, ""
	}

	// 验证码正确，删除缓存中的验证码
	if storedCode.Code == code {
		codeCache.Lock()
		delete(codeCache.m, email)
		codeCache.Unlock()
		return true, storedCode.ASN
	}

	return false, ""
}
