package wechat

import (
	"crypto/sha1"
	"fmt"
	"net/url"
	"time"

	uuid "github.com/satori/go.uuid"
)

// 按照微信文档说明，为前端请求进行签名：
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#62
func JSSign(api string) (nonce string, ts int64, sign string) {
	nonce = uuid.NewV4().String()
	ts = time.Now().Unix()

	shouldSign := url.Values{}
	shouldSign.Set("noncestr", nonce)
	shouldSign.Set("jsapi_ticket", GetJSTicket())
	shouldSign.Set("timestamp", fmt.Sprint(ts))
	//shouldSign.Set("url", api)

	rawstr := shouldSign.Encode()
	rawstr = rawstr + fmt.Sprintf("&url=%s", api)
	logger.Debugf("js api sign: %s", rawstr)

	hashed := sha1.Sum([]byte(rawstr))
	sign = fmt.Sprintf("%x", hashed[:])

	return
}
