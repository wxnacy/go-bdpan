package bdpan

import (
	"net/http"
	"sync"
	"time"

	sdk "github.com/wxnacy/go-bdpan/openapi"
)

var (
	client   *sdk.APIClient
	once     sync.Once
	clientMu sync.RWMutex
)

// GetClient 返回 sdk.APIClient 的单例实例
func GetClient() *sdk.APIClient {
	clientMu.RLock()
	current := client
	clientMu.RUnlock()
	if current != nil {
		return current
	}

	once.Do(func() {
		clientMu.Lock()
		defer clientMu.Unlock()
		if client == nil {
			client = initClient()
		}
	})

	clientMu.RLock()
	defer clientMu.RUnlock()
	return client
}

// SetHTTPClient 替换全局 HTTP client，用于配置超时、代理、连接池等。
//
// 应在程序启动阶段调用一次，不要在并发请求过程中反复替换；传 nil 恢复默认 client。
func SetHTTPClient(httpClient *http.Client) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	clientMu.Lock()
	defer clientMu.Unlock()
	if client == nil {
		client = initClient()
	}
	client.GetConfig().HTTPClient = httpClient
}

// SetHTTPTimeout 是 SetHTTPClient 的便捷封装，只设置请求总超时。
func SetHTTPTimeout(timeout time.Duration) {
	SetHTTPClient(&http.Client{Timeout: timeout})
}

func initClient() *sdk.APIClient {
	configuration := sdk.NewConfiguration()
	return sdk.NewAPIClient(configuration)
}
