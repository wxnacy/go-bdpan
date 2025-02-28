package bdpan

import (
	"sync"

	sdk "github.com/wxnacy/go-bdpan/openapi"
)

var (
	client *sdk.APIClient
	once   sync.Once
)

// GetClient 返回 sdk.APIClient 的单例实例
func GetClient() *sdk.APIClient {
	if client == nil {
		once.Do(func() {
			client = initClient()
		})
	}
	return client
}

func initClient() *sdk.APIClient {
	configuration := sdk.NewConfiguration()
	return sdk.NewAPIClient(configuration)
}
