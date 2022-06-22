package wechat

import (
	"fmt"
	"time"

	reqv3 "github.com/imroc/req/v3"
)

type _WXPhoneResp struct {
	wxCommonResp
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`     // 完整电话号码，包括区号
		PurePhoneNumber string `json:"purePhoneNumber"` // 没有区号的手机号
		CountryCode     string `json:"countryCode"`     // 区号
	} `json:"phone_info"`
}

func (*_WXPhoneResp) Api() string {
	return "/wxa/business/getuserphonenumber"
}

type _WXLoginResp struct {
	wxCommonResp
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符，若当前小程序已绑定到微信开放平台帐号下会返回，详见 UnionID 机制说明。
}

func (*_WXLoginResp) Api() string {
	return "/sns/jscode2session"
}

//----------------

func GetPhone(code string) (phoneNumber, purePhoneNumber string, err error) {
	if code == "" {
		logger.Warn("no phone code, keep empty")
		return
	}

	var wxPhone _WXPhoneResp
	client := reqv3.C(). // Use C() to create a client
				EnableDebugLog().
				SetTimeout(5 * time.Second)
	fnPhonenumber := func() (wxcode, error) {
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"code": "%s"}`, code)).
			Post(wxKey.BaseURL + fmt.Sprintf("%s?access_token=%s", wxPhone.Api(), GetAccessToken()))
		if err != nil {
			return nil, err
		}
		err = resp.UnmarshalJson(&wxPhone)
		if err != nil {
			return nil, err
		}
		return &wxPhone, err
	}
	if err = callWX(fnPhonenumber); err != nil {
		return
	}

	phoneNumber, purePhoneNumber = wxPhone.PhoneInfo.PhoneNumber, wxPhone.PhoneInfo.PurePhoneNumber
	return
}

func Login(code string) (openID, unionID, sessionKey string, err error) {
	var wxLogin _WXLoginResp
	client := reqv3.C(). // Use C() to create a client
				SetTimeout(5 * time.Second)
	fnLogin := func() (wxcode, error) {
		resp, err := client.R(). // Use R() to create a request
						SetQueryParam("appid", wxKey.AppID).
						SetQueryParam("secret", wxKey.Secret).
						SetQueryParam("js_code", code).
						SetQueryParam("grant_type", "authorization_code").
						Get(wxKey.BaseURL + wxLogin.Api())
		if err != nil {
			return nil, err
		}
		err = resp.UnmarshalJson(&wxLogin)
		if err != nil {
			return nil, err
		}
		return &wxLogin, nil
	}
	if err = callWX(fnLogin); err != nil {
		return
	}

	openID, unionID, sessionKey = wxLogin.OpenID, wxLogin.UnionID, wxLogin.SessionKey
	return
}
