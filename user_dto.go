package bdpan

import (
	"github.com/wxnacy/go-bdpan/openapi"
	"github.com/wxnacy/go-tools"
)

type GetPanInfoRes struct {
	*openapi.Quotaresponse
}

func (p GetPanInfoRes) GetUsedStr() string {
	return tools.FormatSize(p.GetUsed())
}

func (p GetPanInfoRes) GetTotalStr() string {
	return tools.FormatSize(p.GetTotal())
}

type GetUserInfoRes struct {
	*openapi.Uinforesponse
}

func (u GetUserInfoRes) GetVipName() string {
	switch u.GetVipType() {
	case 0:
		return "普通用户"
	case 1:
		return "普通会员"
	case 2:
		return "超级会员"
	}
	return "未知身份"
}
