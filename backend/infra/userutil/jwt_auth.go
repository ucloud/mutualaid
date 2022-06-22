package userutil

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/ucloud/mutualaid/backend/api/mutualaid/uerrors"
	"google.golang.org/grpc/metadata"
)

var (
	ErrAuthFail     = errors.New("jwt token failed")
	ErrUnknowClaims = errors.New("jwt token unknow claims")
	ErrExpired      = errors.New("jwt is expired")
)

type uidctxkey struct{}
type openidctxkey struct{}

const (
	Authorization = "Authorization"
	Expired       = 60 * 24 * 30 // 分钟
)

type JWT struct {
	key string
}

func NewJWT() *JWT {
	return &JWT{
		key: "5c510195-5545-43c8-93da-bd6db992ddff",
	}
}

type MyClaims struct {
	jwtv4.StandardClaims

	UID    uint64 `json:"uid"`
	OpenID string `json:"openid"`
}

func (j JWT) Auth(uid uint64, openid string) (string, error) {
	myClaims := MyClaims{
		UID:    uid,
		OpenID: openid,
		StandardClaims: jwtv4.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * Expired).Unix(), //设置JWT过期时间
			Issuer:    "mutualaid-csr.ucloud.cn",                    //设置签发人
		},
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims)
	//加盐
	return claims.SignedString([]byte(j.key))
}

// CheckJWT 解析
func (j JWT) CheckJWT(jwtToken string) (map[string]interface{}, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			return nil, ErrExpired
		}

		result := make(map[string]interface{}, 2)
		result["uid"] = claims["uid"]
		result["openid"] = claims["openid"]
		return result, nil
	} else {
		return nil, ErrUnknowClaims
	}
}

//==============

// FromAuthHeader is a "TokenExtractor" that takes a give context and extracts
// the JWT token from the Authorization header.
func FromAuthHeader(tr transport.Transporter) (string, error) {
	authHeader := tr.RequestHeader().Get(Authorization)
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	// TODO: Make this a bit more robust, parsing-wise
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

// NewJWTServerMidware jwt Server中间件
func NewJWTServerMidware(j *JWT) func(handler middleware.Handler) middleware.Handler {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var jwtToken string
			var token map[string]interface{}
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				jwtToken = md.Get("x-md-global-jwt")[0]
			} else if tr, ok := transport.FromServerContext(ctx); ok {
				jwtToken, err = FromAuthHeader(tr)
				if err != nil {
					return nil, err
				}

				// goto PASS 本地测试开关，绕过鉴权

				if jwtToken == "" {
					if tr.Operation() == "/api.mutualaid.MutualAidQuery/Discovery" {
						goto PASS
					}
					return nil, uerrors.ErrorUnloginError("未登录")
				}
			} else {
				// 缺少可认证的token，返回错误
				return nil, ErrAuthFail
			}
			token, err = j.CheckJWT(jwtToken)
			if err != nil {
				// 缺少合法的token，返回错误
				return nil, uerrors.ErrorUnloginError("登录过期")
			}
			ctx = context.WithValue(ctx, uidctxkey{}, token["uid"])
			ctx = context.WithValue(ctx, openidctxkey{}, token["openid"])
		PASS:
			reply, err = handler(ctx, req)
			return
		}
	}
}

// NewJWTClientMidware jwt Client中间件
func NewJWTClientMidware(j *JWT) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				// 存在jwt_token就透传
				jwtToken, err := FromAuthHeader(tr)
				if err != nil {
					return nil, err
				}

				ctx = metadata.AppendToOutgoingContext(ctx, "x-md-global-jwt", jwtToken)
			}
			return handler(ctx, req)
		}
	}
}

func ExtractUID(ctx context.Context) uint64 {
	v := ctx.Value(uidctxkey{})
	if v == nil {
		return 0
	}
	if uid, ok := v.(float64); ok {
		return uint64(uid)
	}

	return 0
}

func ExtractOpenID(ctx context.Context) string {
	v := ctx.Value(openidctxkey{})
	if v == nil {
		return ""
	}
	if openid, ok := v.(string); ok {
		return openid
	}

	return ""
}
