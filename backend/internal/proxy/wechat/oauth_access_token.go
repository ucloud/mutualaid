package wechat

import (
	"time"

	reqv3 "github.com/imroc/req/v3"
)

type _OAuth2Resp struct {
	wxCommonResp
	AccessToken  string `json:"access_token"`  // 网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	ExpiresIn    int    `json:"expires_in"`    // access_token接口调用凭证超时时间，单位（秒）
	RefreshToken string `json:"refresh_token"` // 用户刷新access_token
	OpenID       string `json:"openid"`        // 用户唯一标识，请注意，在未关注公众号时，用户访问公众号的网页，也会产生一个用户和公众号唯一的OpenID
	Scope        string `json:"scope"`         // 用户授权的作用域，使用逗号（,）分隔
}

func (*_OAuth2Resp) Api() string {
	return "/sns/oauth2/access_token"
}

//----------------

func OAuth2(code string) (openID string, oauthAccssToken string, err error) {
	var wxOAuth2 _OAuth2Resp
	client := reqv3.C(). // Use C() to create a client
				SetTimeout(5 * time.Second)
	fnOAuth := func() (wxcode, error) {
		resp, err := client.R(). // Use R() to create a request
						SetQueryParam("appid", wxKey.WPAppID).
						SetQueryParam("secret", wxKey.WPSecret).
						SetQueryParam("code", code).
						SetQueryParam("grant_type", "authorization_code").
						Get(wxKey.BaseURL + wxOAuth2.Api())
		if err != nil {
			return nil, err
		}
		err = resp.UnmarshalJson(&wxOAuth2)
		if err != nil {
			return nil, err
		}
		return &wxOAuth2, nil
	}
	if err = callWX(fnOAuth); err != nil {
		return
	}

	openID, oauthAccssToken = wxOAuth2.OpenID, wxOAuth2.AccessToken
	return
}
