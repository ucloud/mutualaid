package wechat

import (
	"strings"
	"time"

	reqv3 "github.com/imroc/req/v3"
)

type TplItem struct {
	Key   string
	Value string
}

type _SubscribeMsgResp struct {
	wxCommonResp
}

func (*_SubscribeMsgResp) Api() string {
	return "/cgi-bin/message/subscribe/send"
}

type _SubscribeMsgRequest struct {
	AccessToken      string                   `json:"access_token"` //	接口调用凭证
	ToUser           string                   `json:"touser"`       //	用户openid，可以是小程序的openid，也可以是mp_template_msg.appid对应的公众号的openid
	TemplateID       string                   `json:"template_id"`
	Page             string                   `json:"page"`
	MiniProgramState string                   `json:"miniprogram_state"` // 跳转小程序类型：developer为开发版；trial为体验版；formal为正式版；默认为正式版
	Lang             string                   `json:"lang"`              // 入小程序查看”的语言类型，支持zh_CN(简体中文)、en_US(英文)、zh_HK(繁体中文)、zh_TW(繁体中文)，默认为zh_CN
	Data             map[string]_MessageValue `json:"data"`              //	公众号模板消息的数据
}

var (
	argLenLimit = map[string]int{
		"thing":            20,
		"number":           32,
		"letter":           32,
		"symbol":           5,
		"character_string": 32,
		"amount":           10,
		"phone_number":     17,
		"car_number":       8,
		"name":             10,
		"phrase":           5,
	}
)

func SubscribeMsgSend(openid, tplName, tplID string, content map[string]TplItem) (err error) {
	if openid == "" {
		// 没有小程序openid, 应该是公众号H5用户
		return nil
	}

	tplArg := make(map[string]_MessageValue)
	for _, v := range content {
		mv := []rune(v.Value)
		for lk, lv := range argLenLimit {
			if strings.HasPrefix(v.Key, lk) {
				if len(v.Value) > lv {
					mv = mv[:lv]
					break
				}
			}
		}

		tplArg[v.Key] = _MessageValue{Value: string(mv)}
	}

	var msgResp _SubscribeMsgResp
	var msgreq = _SubscribeMsgRequest{
		AccessToken: GetAccessToken(),
		TemplateID:  tplID,
		ToUser:      openid,
		Page:        "/pages/detail/detail?id=" + content["求助编号"].Value,
		Data:        tplArg,
	}
	logger.Infof("send subscribe message, templte: %s, args: %+v", tplName, content)

	client := reqv3.C(). // Use C() to create a client
				EnableDumpAll().
				SetTimeout(5 * time.Second)
	fnSubMsgSend := func() (wxcode, error) {
		resp, err := client.R(). // Use R() to create a request
						SetQueryParam("access_token", GetAccessToken()).
						SetBodyJsonMarshal(msgreq).
						Post(wxKey.BaseURL + msgResp.Api())
		if err != nil {
			return nil, err
		}
		err = resp.UnmarshalJson(&msgResp)
		if err != nil {
			return nil, err
		}
		return &msgResp, nil
	}
	if err = callWX(fnSubMsgSend); err != nil {
		return err
	}

	return nil
}
