package wechat

import (
	"time"

	reqv3 "github.com/imroc/req/v3"
)

type _SnsUserInfoResp struct {
	wxCommonResp
	OpenID     string   `json:"openid"`     // 用户的唯一标识
	Nickname   string   `json:"nickname"`   // 用户昵称
	Sex        int      `json:"sex"`        // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Province   string   `json:"province"`   // 用户个人资料填写的省份
	City       string   `json:"city"`       // 普通用户个人资料填写的城市
	Country    string   `json:"country"`    // 国家，如中国为CN
	HeadImgUrl string   `json:"headimgurl"` // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	Privilege  []string `json:"privilege"`  // 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	UnionID    string   `json:"unionid"`    // 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
}

func (*_SnsUserInfoResp) Api() string {
	return "/sns/userinfo"
}

//----------------

func SnsUserInfo(openID string, oauthAccssToken string) (unionID, nickname, headImgUrl string, err error) {
	var wxUserInfo _SnsUserInfoResp
	client := reqv3.C(). // Use C() to create a client
				EnableDumpAll().
				SetTimeout(5 * time.Second)
	fnUserInfo := func() (wxcode, error) {
		resp, err := client.R(). // Use R() to create a request
						SetQueryParam("access_token", oauthAccssToken).
						SetQueryParam("openid", openID).
						SetQueryParam("lang", "zh_CN").
						Get(wxKey.BaseURL + wxUserInfo.Api())
		if err != nil {
			return nil, err
		}
		err = resp.UnmarshalJson(&wxUserInfo)
		if err != nil {
			return nil, err
		}
		return &wxUserInfo, nil
	}
	if err = callWX(fnUserInfo); err != nil {
		return
	}

	unionID, nickname, headImgUrl = wxUserInfo.UnionID, wxUserInfo.Nickname, wxUserInfo.HeadImgUrl
	return
}
