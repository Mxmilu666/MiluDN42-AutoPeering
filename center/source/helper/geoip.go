package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// IsIPCNOnline 通过在线API判断IP是否属于中国
func IsIPCNOnline(ip string) (bool, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=countryCode", ip)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var result struct {
		CountryCode string `json:"countryCode"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	return result.CountryCode == "CN", nil
}

// IsIPCNByIPInfo 通过 ipinfo.io 在线API判断IP是否属于中国
func IsIPCNByIPInfo(ip string) (bool, error) {
	url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var result struct {
		Country string `json:"country"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	return result.Country == "CN", nil
}

// IsIPCNAny 并发查询两个API，只要有一个判断为中国IP就返回true，否则false
func IsIPCNAny(ip string) (bool, error) {
	var (
		wg       sync.WaitGroup
		resultCh = make(chan bool, 2)
		errCh    = make(chan error, 2)
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		ok, err := IsIPCNOnline(ip)
		if err != nil {
			errCh <- err
			return
		}
		if ok {
			resultCh <- true
		}
	}()
	go func() {
		defer wg.Done()
		ok, err := IsIPCNByIPInfo(ip)
		if err != nil {
			errCh <- err
			return
		}
		if ok {
			resultCh <- true
		}
	}()

	// 等待结果
	var errs []error
	for i := 0; i < 2; i++ {
		select {
		case <-resultCh:
			return true, nil
		case err := <-errCh:
			errs = append(errs, err)
		}
	}
	if len(errs) == 2 {
		return false, fmt.Errorf("ip-api error: %v; ipinfo error: %v", errs[0], errs[1])
	}
	return false, nil
}
