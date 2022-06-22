package wechat

import (
	"context"
	"errors"
	"fmt"
	"time"

	reqv3 "github.com/imroc/req/v3"
)

const (
	redlockKey     = "redlock:%s"
	redlockExp     = time.Second * 5
	accessTokenKey = "wx:accesstoken:%s"

	retryTimes = 2
)

var ErrRedLockExpired = errors.New("redlock expired")

type RequestWXToken func() (string, int, error)

type _TokenResponse struct {
	wxCommonResp
	AccessToken string `json:"access_token"` // 获取到的凭证
	Expires     int    `json:"expires_in"`   // 凭证有效时间，单位：秒。目前是7200秒之内的值。
}

func (*_TokenResponse) Api() string {
	return "/cgi-bin/token"
}

type _APITicketResponse struct {
	wxCommonResp
	Ticket  string `json:"ticket"`     // 获取到的凭证
	Expires int    `json:"expires_in"` // 凭证有效时间，单位：秒。目前是7200秒之内的值。
}

func (*_APITicketResponse) Api() string {
	return "/cgi-bin/ticket/getticket"
}

//--------------------

func GetAccessToken() string {
	ctx := context.Background()

	for i := 0; i < retryTimes; i++ {
		ckey := fmt.Sprintf(accessTokenKey, wxKey.AppID)
		result := cache.WithContext(ctx).Get(ctx, ckey)
		if result.Err() != nil {
			err := globalCache(ckey, freshWXAccessToken(wxKey.AppID, wxKey.Secret))
			if err != nil {
				if err == ErrRedLockExpired {
					time.Sleep(time.Microsecond * time.Duration(i*50))
					continue
				}

				break
			}
		}

		return result.Val()
	}

	logger.Errorf("无法取得微信access token")
	return ""
}

func GetMPAccessToken() string {
	ctx := context.Background()

	for i := 0; i < retryTimes; i++ {
		ckey := fmt.Sprintf(accessTokenKey, wxKey.WPAppID)
		result := cache.WithContext(ctx).Get(ctx, ckey)
		if result.Err() != nil {
			err := globalCache(ckey, freshWXAccessToken(wxKey.WPAppID, wxKey.WPSecret))
			if err != nil {
				if err == ErrRedLockExpired {
					time.Sleep(time.Microsecond * time.Duration(i*50))
					continue
				}

				break
			}
		}

		return result.Val()
	}

	logger.Errorf("无法取得微信公众号access token")
	return ""
}

func GetJSTicket() string {
	freshJSTicket := func() (string, int, error) {
		var ticket _APITicketResponse
		client := reqv3.C(). // Use C() to create a client
					EnableDumpAll().
					SetTimeout(5 * time.Second)
		mpAccessToken := GetMPAccessToken()
		if mpAccessToken == "" {
			return "", 0, errors.New("no access token")
		}

		fnToken := func() (wxcode, error) {

			_, err := client.R().
				SetResult(&ticket).
				SetQueryParams(map[string]string{
					"type":         "jsapi",
					"access_token": mpAccessToken}).
				Get(wxKey.BaseURL + ticket.Api())
			return &ticket, err
		}

		if err := callWX(fnToken); err != nil {
			return "", 0, err
		}

		return ticket.Ticket, ticket.Expires, nil
	}

	ctx := context.Background()
	for i := 0; i < retryTimes; i++ {
		ckey := fmt.Sprintf("wx:js:ticket:%s", wxKey.WPAppID)
		result := cache.WithContext(ctx).Get(ctx, ckey)
		if result.Err() != nil {
			err := globalCache(ckey, freshJSTicket)
			if err != nil {
				if err == ErrRedLockExpired {
					time.Sleep(time.Microsecond * time.Duration(i*50))
					continue
				}

				break
			}
		}

		return result.Val()
	}

	logger.Errorf("无法取得微信jsapi_ticket")
	return ""
}

//-------------

func freshWXAccessToken(appid, secret string) RequestWXToken {
	return func() (string, int, error) {
		var access _TokenResponse
		client := reqv3.C(). // Use C() to create a client
					EnableDumpAll().
					SetTimeout(5 * time.Second)
		fnToken := func() (wxcode, error) {
			_, err := client.R().
				SetResult(&access).
				SetQueryParams(map[string]string{
					"grant_type": "client_credential",
					"appid":      appid,
					"secret":     secret}).
				Get(wxKey.BaseURL + access.Api())
			return &access, err
		}

		if err := callWX(fnToken); err != nil {
			return "", 0, err
		}

		logger.Infof("got weixin access_token for appid: %s: exp: %d", appid, access.Expires)
		return access.AccessToken, access.Expires, nil
	}
}

func globalCache(ckey string, fnRequest RequestWXToken) error {
	ctx := context.Background()

	begin := time.Now()
	result := cache.WithContext(ctx).SetNX(ctx, fmt.Sprintf(redlockKey, ckey), time.Now().Unix(), redlockExp)
	if result.Err() != nil {
		logger.Warnf("get redlock for wx access token failed")
		return result.Err()
	}
	defer cache.WithContext(ctx).Del(ctx, fmt.Sprintf(redlockKey, ckey))

	// 取token
	data, exp, err := fnRequest()
	if err != nil {
		return err
	}

	if time.Since(begin) > redlockExp {
		// redlock失效了
		logger.Warnf("redlock for wx access token expired")
		return ErrRedLockExpired
	}

	logger.Info("fresh wx access token sucessed")
	cache.WithContext(ctx).SetEX(ctx,
		ckey,
		data,
		time.Second*time.Duration(float64(exp)*3/4),
	)

	return nil
}
