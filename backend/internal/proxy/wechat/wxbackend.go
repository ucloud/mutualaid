package wechat

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/ucloud/mutualaid/backend/internal/conf"
	"time"
)

type wxcode interface {
	Code() int
	Message() string
	Api() string
}

type wxCommonResp struct {
	ErrCode int    `json:"errcode"` // 错误码
	ErrMsg  string `json:"errmsg"`  // 错误信息
}

func (w *wxCommonResp) Code() int {
	return w.ErrCode
}
func (w *wxCommonResp) Message() string {
	return w.ErrMsg
}

var (
	wxKey  *conf.WXKey
	logger *log.Helper
	cache  *redis.ClusterClient

	wxRetMessage = map[int]string{
		-1:    "系统繁忙，此时请开发者稍候再试",
		0:     "请求成功",
		40001: "AppSecret 错误或者 AppSecret 不属于这个小程序，请开发者确认 AppSecret 的正确性",
		40002: "请确保 grant_type 字段值为 client_credential",
		40013: "不合法的 AppID，请开发者检查 AppID 的正确性，避免异常字符，注意大小写",
		40029: "不合法的code（code不存在、已过期或者使用过）",
		45011: "频率限制，每个用户每分钟100次",
		40226: "高风险等级用户，小程序登录拦截 。风险等级详见用户安全解方案",
		40129: "场景值错误（目前支持场景 1 资料；2 评论；3 论坛；4 社交日志）",
		43104: "appid与openid不匹配",
		43302: "方法调用错误，请用post方法调用",
		44002: "传递的参数为空",
		47001: "传递的参数格式不对",
		61010: "用户访问记录超时（用户未在近两小时访问小程序）",
		40164: "调用接口的IP地址不在白名单中，请在接口IP白名单中进行设置。",
		89503: "此IP调用需要管理员确认,请联系管理员",
		89501: "此IP正在等待管理员确认,请联系管理员",
		89506: "24小时内该IP被管理员拒绝调用两次，24小时内不可再使用该IP调用",
		89507: "1小时内该IP被管理员拒绝调用一次，1小时内不可再使用该IP调用",
		40003: "touser字段openid为空或者不正确",
		40037: "订阅模板id为空不正确",
		43101: "用户拒绝接受消息，如果用户之前曾经订阅过，则表示用户取消了订阅关系",
		47003: "模板参数不准确，可能为空或者不满足规则，errmsg会提示具体是哪个字段出错",
		41030: "page路径不正确，需要保证在现网版本小程序中存在，与app.json保持一致",
	}
)

func InitWXBackend(conf *conf.Proxy, c *conf.Data, l log.Logger) (func(), error) {
	wxKey = conf.GetWxkey()
	logger = log.NewHelper(l)

	const (
		defaultPoolSize    = 60
		defaultMinIdleConn = 20
		defaultMaxRetries  = 3
		defaultTimeOut     = 1 * time.Second
		defaultIdleTimeout = 5 * time.Minute
	)
	poolSize := c.Redis.GetPoolSize()
	if poolSize == 0 {
		poolSize = defaultPoolSize
	}
	minIdleConn := c.Redis.MinIdleConn
	if minIdleConn <= 0 {
		minIdleConn = defaultMinIdleConn
	}

	cache = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        c.Redis.GetAddrs(),
		Username:     c.Redis.GetUserName(),
		Password:     c.Redis.GetPassword(),
		MaxRetries:   defaultMaxRetries,
		ReadTimeout:  defaultTimeOut,
		WriteTimeout: defaultTimeOut,
		PoolSize:     int(poolSize),
		MinIdleConns: int(minIdleConn),
		IdleTimeout:  defaultIdleTimeout,
	})
	status := cache.Ping(context.Background())
	if status.Err() != nil {
		return nil, fmt.Errorf("redis open error: %v", status.Err())
	}

	cleanup := func() {
		cache.Close()
		log.NewHelper(l).Info("redis closed")
	}

	return cleanup, nil
}

func callWX(fn func() (wxcode, error)) error {
	var result wxcode
	var err error
	for i := 0; i < 3; i++ {
		result, err = fn()
		if err != nil {
			logger.Errorf("request wx backend failed, retry now", err)
			continue
		}
		if result.Code() != 0 {
			if result.Code() == -1 {
				// 立即重试
				time.Sleep(time.Millisecond * 50 * time.Duration(i+1))
				continue
			}
			err = errors.Newf(result.Code(), "call weixin api failed", "api: %s; code: %d; %s; %s", result.Api(), result.Code(), result.Message(), wxRetMessage[result.Code()])
			logger.Error(err)
			return err
		}

		break
	}

	return err
}
