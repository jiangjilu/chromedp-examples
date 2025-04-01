package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// API 响应结构
type APIResponse struct {
	URL        string `json:"url,omitempty"`
	Title      string `json:"title,omitempty"`
	Screenshot string `json:"screenshot,omitempty"`
	Message    string `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
}

// 测试 HTTP 请求
func testAPI(t *testing.T, endpoint string, expectedStatus int, responseStruct *APIResponse) {
	req, err := http.NewRequest("GET", endpoint, nil)
	assert.NoError(t, err)

	// 发送请求
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	assert.Equal(t, expectedStatus, resp.StatusCode)

	// 解析 JSON 响应
	err = json.NewDecoder(resp.Body).Decode(responseStruct)
	assert.NoError(t, err)
}

// 测试 /navigate API
func TestNavigateAPI(t *testing.T) {
	var res APIResponse
	testAPI(t, "http://localhost:8080/navigate?url=https://www.example.com", http.StatusOK, &res)

	assert.Equal(t, "https://www.example.com", res.URL)
	assert.Equal(t, "Page loaded successfully", res.Message)
}

// 测试 /title API
func TestTitleAPI(t *testing.T) {
	var res APIResponse
	testAPI(t, "http://localhost:8080/title?url=https://www.example.com", http.StatusOK, &res)

	assert.Equal(t, "https://www.example.com", res.URL)
	assert.NotEmpty(t, res.Title)
}

// 测试 /screenshot API
func TestScreenshotAPI(t *testing.T) {
	var res APIResponse
	testAPI(t, "http://localhost:8080/screenshot?url=https://www.example.com", http.StatusOK, &res)

	assert.Equal(t, "https://www.example.com", res.URL)
	assert.NotEmpty(t, res.Screenshot)
}
