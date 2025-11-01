package bdpan

import (
	"context"
)

// 通过 code 获取用户身份 token
func GetDeviceToken(appKey, secretKey, code string) (*GetDeviceTokenRes, error) {
	_, r, err := GetClient().
		AuthApi.
		OauthTokenDeviceToken(context.Background()).
		Code(code).
		ClientId(appKey).
		ClientSecret(secretKey).
		Execute()
	return ToInterface[GetDeviceTokenRes](r, err)
}

// 获取登录使用 code 二维码
// scope: 一般使用 basic,netdisk
func GetDeviceCode(appKey, scope string) (*GetDeviceCodeRes, error) {
	_, r, err := GetClient().
		AuthApi.
		OauthTokenDeviceCode(context.Background()).
		ClientId(appKey).
		Scope(scope).
		Execute()
	return ToInterface[GetDeviceCodeRes](r, err)
}
