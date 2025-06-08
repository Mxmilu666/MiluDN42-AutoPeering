package helper

import (
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source"
)

// whoisQuery 查询 whois
func whoisQuery(server string, query string) (string, error) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(query + "\r\n"))
	if err != nil {
		return "", err
	}

	result, err := io.ReadAll(conn)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// extractField  从 whois 查询结果中提取指定字段
func extractField(data, field string) []string {
	regex := regexp.MustCompile(field + `:\s+(.+)`)
	matches := regex.FindAllStringSubmatch(data, -1)

	var results []string
	for _, match := range matches {
		results = append(results, strings.TrimSpace(match[1]))
	}
	return results
}

// GetASNMaintainerEmails 获取 ASN 的维护者邮箱
func GetASNMaintainerEmails(asn string) ([]string, error) {
	if source.AppConfig == nil {
		return nil, fmt.Errorf("AppConfig is not initialized")
	}
	whoisHost := source.AppConfig.Whois.Host + ":" + fmt.Sprint(source.AppConfig.Whois.Port)
	if whoisHost == "" {
		return nil, fmt.Errorf("whois server is not configured")
	}

	asnResult, err := whoisQuery(whoisHost, asn)
	if err != nil {
		return nil, fmt.Errorf("ASN query failed: %v", err)
	}

	// 优先查找 admin-c 和 tech-c
	contacts := extractField(asnResult, "admin-c")
	if len(contacts) == 0 {
		contacts = extractField(asnResult, "tech-c")
	}
	if len(contacts) == 0 {
		return nil, fmt.Errorf("admin-c or tech-c field not found")
	}

	var emails []string
	for _, contact := range contacts {
		contactResult, err := whoisQuery(whoisHost, contact)
		if err != nil {
			return nil, fmt.Errorf("contact query failed: %v", err)
		}
		emailList := extractField(contactResult, "e-mail")
		emails = append(emails, emailList...)
	}

	if len(emails) == 0 {
		return nil, fmt.Errorf("email not found")
	}

	return emails, nil
}

// GetASNByIP 根据 IP 查询 ASN
func GetASNByIP(ip string) (string, error) {
	if source.AppConfig == nil {
		return "", fmt.Errorf("AppConfig is not initialized")
	}
	whoisHost := source.AppConfig.Whois.Host + ":" + fmt.Sprint(source.AppConfig.Whois.Port)
	if whoisHost == ":" {
		return "", fmt.Errorf("whois server is not configured")
	}

	ipResult, err := whoisQuery(whoisHost, ip)
	if err != nil {
		return "", fmt.Errorf("IP whois query failed: %v", err)
	}

	// 优先查找 origin 字段
	asns := extractField(ipResult, "origin")
	if len(asns) == 0 {
		// 某些 whois 服务器可能用 originas 字段
		asns = extractField(ipResult, "originas")
	}
	if len(asns) == 0 {
		return "", fmt.Errorf("ASN (origin) field not found in whois result")
	}

	return asns[0], nil
}
