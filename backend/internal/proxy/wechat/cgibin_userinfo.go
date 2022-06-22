package wechat

import (
	"time"

	reqv3 "github.com/imroc/req/v3"
)

type _UserInfoResp struct {
	wxCommonResp
	Subscribe      int    `json:"subscribe"`       // 用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。
	OpenID         string `json:"openid"`          // 用户的标识，对当前公众号唯一
	Language       string `json:"language"`        // 用户的语言，简体中文为zh_CN
	SubscribeTime  int    `json:"subscribe_time"`  // 用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
	UnionID        string `json:"unionid"`         // 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	Remark         string `json:"remark"`          // 公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	GroupID        int    `json:"groupid"`         // 用户所在的分组ID（兼容旧的用户分组接口）
	TagidList      []int  `json:"tagid_list"`      // 用户被打上的标签ID列表
	SubscribeScene string `json:"subscribe_scene"` // 返回用户关注的渠道来源，ADD_SCENE_SEARCH 公众号搜索，ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，ADD_SCENE_PROFILE_CARD 名片分享，ADD_SCENE_QR_CODE 扫描二维码，ADD_SCENE_PROFILE_LINK 图文页内名称点击，ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，ADD_SCENE_PAID 支付后关注，ADD_SCENE_WECHAT_ADVERTISEMENT 微信广告，ADD_SCENE_REPRINT 他人转载 ,ADD_SCENE_LIVESTREAM 视频号直播，ADD_SCENE_CHANNELS 视频号 , ADD_SCENE_OTHERS 其他
	QrScene        int    `json:"qr_scene"`        // 二维码扫码场景（开发者自定义）
	QrSceneStr     string `json:"qr_scene_str"`    // 二维码扫码场景描述（开发者自定义）
}

func (*_UserInfoResp) Api() string {
	return "/cgi-bin/user/info"
}

//----------------

func UserInfo(openID string) (unionID string, err error) {
	var wxUserInfo _UserInfoResp
	client := reqv3.C(). // Use C() to create a client
				EnableDumpAll().
				SetTimeout(5 * time.Second)
	fnUserInfo := func() (wxcode, error) {
		resp, err := client.R(). // Use R() to create a request
						SetQueryParam("access_token", GetMPAccessToken()).
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

	unionID = wxUserInfo.UnionID
	return
}
