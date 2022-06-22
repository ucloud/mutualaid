package wechat

import (
	"time"

	reqv3 "github.com/imroc/req/v3"
)

type _UniformMsgResp struct {
	wxCommonResp
}

func (*_UniformMsgResp) Api() string {
	return "/cgi-bin/message/wxopen/template/uniform_send"
}

type _UniformMsgRequest struct {
	AccessToken   string             `json:"access_token"`    //	接口调用凭证
	ToUser        string             `json:"touser"`          //	用户openid，可以是小程序的openid，也可以是mp_template_msg.appid对应的公众号的openid
	MPTemplateMsg _MPMessageTemplate `json:"mp_template_msg"` //	公众号模板消息相关的信息，可以参考公众号模板消息接口；有此节点并且没有weapp_template_msg节点时，发送公众号模板消息
}

type _MessageValue struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

type _MiniP struct {
	AppID string `json:"appid"` //	公众号appid，要求与小程序有绑定且同主体
	Page  string `json:"page"`  //	公众号模板消息所要跳转的url
}

type _MPMessageTemplate struct {
	AppID       string                   `json:"appid"`       //	公众号appid，要求与小程序有绑定且同主体
	TemplateID  string                   `json:"template_id"` //	公众号模板id
	Url         string                   `json:"url"`         //	公众号模板消息所要跳转的url
	MiniProgram _MiniP                   `json:"miniprogram"` //	公众号模板消息所要跳转的小程序，小程序的必须与公众号具有绑定关系
	Data        map[string]_MessageValue `json:"data"`        //	公众号模板消息的数据
}

func UniformMessageSend(openid, content string) (err error) {
	if openid == "" {
		// 没有小程序openid, 应该是公众号H5用户
		return nil
	}

	tplArg := make(map[string]_MessageValue)
	tplArg["first"] = _MessageValue{Value: "ucloud"}
	tplArg["keyword1"] = _MessageValue{Value: content}
	tplArg["remark"] = _MessageValue{Value: "优刻得"}

	var msgResp _UniformMsgResp
	var msgreq = _UniformMsgRequest{
		AccessToken: GetAccessToken(),
		ToUser:      openid,
		MPTemplateMsg: _MPMessageTemplate{
			AppID:      wxKey.WPAppID,
			TemplateID: "TODO",
			MiniProgram: _MiniP{
				AppID: wxKey.AppID,
				Page:  "TODO",
			},
			Data: tplArg,
		},
	}

	client := reqv3.C(). // Use C() to create a client
				EnableDumpAll().
				SetTimeout(5 * time.Second)
	fnUniformMsgSend := func() (wxcode, error) {
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
	if err = callWX(fnUniformMsgSend); err != nil {
		return err
	}

	return nil
}
