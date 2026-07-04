package sdk

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/pkg/logger"
)

// 官方 SDK API 基址
const officialSDKBase = "https://jp-sdk-api.bluearchive.cafe"

// officialHTTPClient 共享的 HTTP 客户端(用于代理请求到官方)
var officialHTTPClient = &http.Client{
	Timeout: 15 * time.Second,
}

// proxyToOfficial 将请求代理到官方 SDK 服务器
// path 是不带域名的路径(如 /common/config)
// 客户端收到的是官方真实响应,避免因字段缺失导致客户端崩溃
func proxyToOfficial(c *gin.Context, path string) {
	officialURL := officialSDKBase + path

	// 读取请求体
	reqBody, _ := io.ReadAll(c.Request.Body)
	logger.Debug("proxyToOfficial %s, 请求体: %s", path, string(reqBody))

	// 创建转发请求
	officialReq, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, officialURL, nil)
	if err != nil {
		logger.Debug("proxyToOfficial 创建请求失败: %s", err)
		c.JSON(500, gin.H{"Code": 500, "Msg": "create request error"})
		return
	}

	// 复制所有请求头(让官方服务器看到和客户端直连时一样的头)
	for k, vs := range c.Request.Header {
		for _, v := range vs {
			officialReq.Header.Add(k, v)
		}
	}
	// 确保有 Content-Type
	if officialReq.Header.Get("Content-Type") == "" {
		officialReq.Header.Set("Content-Type", "application/json")
	}

	// 设置请求体
	if len(reqBody) > 0 {
		officialReq.Body = io.NopCloser(newBytesReader(reqBody))
		officialReq.ContentLength = int64(len(reqBody))
	}

	// 发送请求
	resp, err := officialHTTPClient.Do(officialReq)
	if err != nil {
		logger.Debug("proxyToOfficial 请求官方失败: %s, path: %s", err, path)
		// 降级:返回最小有效响应
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.String(200, `{"Code":200,"Data":{},"Msg":"OK"}`)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Debug("proxyToOfficial 读取响应失败: %s", err)
		c.JSON(500, gin.H{"Code": 500, "Msg": "read response error"})
		return
	}

	logger.Debug("proxyToOfficial %s, 官方响应状态: %d, 长度: %d", path, resp.StatusCode, len(respBody))

	// 复制响应头
	for k, vs := range resp.Header {
		for _, v := range vs {
			c.Header(k, v)
		}
	}

	// 返回官方响应
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

// bytesReader 避免引入 bytes 包(虽然其实没冲突,但保持独立)
type bytesReader struct {
	data []byte
	off  int
}

func newBytesReader(b []byte) *bytesReader {
	return &bytesReader{data: b}
}

func (r *bytesReader) Read(p []byte) (n int, err error) {
	if r.off >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.off:])
	r.off += n
	return n, nil
}
