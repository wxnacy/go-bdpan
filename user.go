package bdpan

import (
	"context"
)

func GetPanInfo(accessToken string) (*GetPanInfoRes, error) {
	_, r, err := GetClient().
		UserinfoApi.
		Apiquota(context.Background()).AccessToken(
		accessToken).Checkexpire(1).Checkfree(1).Execute()
	if err != nil {
		return nil, err
	}
	return ToInterface[GetPanInfoRes](r)
}

func GetUserInfo(accessToken string) (*GetUserInfoRes, error) {
	_, r, err := GetClient().
		UserinfoApi.
		Xpannasuinfo(context.Background()).
		AccessToken(accessToken).
		Execute()
	if err != nil {
		return nil, err
	}

	return ToInterface[GetUserInfoRes](r)
}
