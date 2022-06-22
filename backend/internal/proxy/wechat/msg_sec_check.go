package wechat

import (
	"fmt"
	"time"

	reqv3 "github.com/imroc/req/v3"
)

var label = map[int]string{
	100:   "正常",
	10001: "广告",
	20001: "时政",
	20002: "色情",
	20003: "辱骂",
	20006: "违法犯罪",
	20008: "欺诈",
	20012: "低俗",
	20013: "版权",
	21000: "其他",
}

type _MsgCheckResp struct {
	wxCommonResp
	TraceID string `json:"trace_id"` // 唯一请求标识，标记单次请求
	Result  struct {
		Suggest string `json:"suggest"` // 建议，有risky、pass、review三种值
		Label   int    `json:"label"`   // 命中标签枚举值
	} `json:"result"` // 综合结果
	Detail []struct {
		Strategy string `json:"strategy"` // 策略类型
		Errcode  int    `json:"errcode"`  // 错误码，仅当该值为0时，该项结果有效
		Suggest  string `json:"suggest"`  // 建议，有risky、pass、review三种值
		Label    int    `json:"label"`    // 命中标签枚举值
		Prob     int    `json:"prob"`     // 0-100，代表置信度，越高代表越有可能属于当前返回的标签（label）
		Keyword  string `json:"keyword"`  // 命中的自定义关键词
	} `json:"detail"` // 详细检测结果
}

func (*_MsgCheckResp) Api() string {
	return "/wxa/msg_sec_check"
}

func MessageSecCheck(openid, content string) (err error) {
	if openid == "" {
		// 没有小程序openid, 应该是公众号H5用户
		return nil
	}

	var msgCheckResp _MsgCheckResp
	var msgseccheckreq = struct {
		Version   int    `json:"version"`             // 是	接口版本号，2.0版本为固定值2
		Openid    string `json:"openid"`              // 是	用户的openid（用户需在近两小时访问过小程序）
		Scene     int    `json:"scene"`               // 是	场景枚举值（1 资料；2 评论；3 论坛；4 社交日志）
		Content   string `json:"content"`             // 是
		Nickname  string `json:"nickname",omitempty`  // 否
		Title     string `json:"title",omitempty`     // 否
		Signature string `json:"signature",omitempty` //	否	个性签名，该参数仅在资料类场景有效(scene=1)，需使用UTF-8编码
	}{
		Version: 2,
		Openid:  openid,
		Scene:   2,
		Content: content,
	}

	client := reqv3.C(). // Use C() to create a client
				EnableDumpAll().
				SetTimeout(5 * time.Second)
	fnMsgCheck := func() (wxcode, error) {
		resp, err := client.R(). // Use R() to create a request
						SetQueryParam("access_token", GetAccessToken()).
						SetBodyJsonMarshal(msgseccheckreq).
						Post(wxKey.BaseURL + msgCheckResp.Api())
		if err != nil {
			return nil, err
		}
		err = resp.UnmarshalJson(&msgCheckResp)
		if err != nil {
			return nil, err
		}
		return &msgCheckResp, nil
	}
	if err = callWX(fnMsgCheck); err != nil {
		return
	}

	if msgCheckResp.Result.Suggest == "risky" {
		logger.Errorf("用户发布的内容涉及%s的不当言论; %s", label[msgCheckResp.Result.Label], content)
		for _, d := range msgCheckResp.Detail {
			if d.Prob > 90 {
				return fmt.Errorf("您发表的内容有涉及%s的不当用词, %s", label[d.Label], d.Keyword)
			}
		}
	}
	if msgCheckResp.Result.Suggest == "review" {
		logger.Warnf("用户发布的内容涉及：%s，需要人工复核; %s", label[msgCheckResp.Result.Label], content)
	}

	return
}
