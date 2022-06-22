package server

import (
	"context"
	stdhttp "net/http"

	"github.com/ucloud/mutualaid/backend/infra/userutil"

	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	pb "github.com/ucloud/mutualaid/backend/api/mutualaid"
	"github.com/ucloud/mutualaid/backend/internal/conf"
	"github.com/ucloud/mutualaid/backend/internal/service"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(
	c *conf.Server,
	userService *service.UserService,
	mutualAidQueryService *service.MutualAidQueryService,
	ExamineAidService *service.ExamineAidService,
	aidQueryService *service.UserAidQueryService,
	aidManagerService *service.UserAidManagerService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.ResponseEncoder(responseEncoder),
		http.Middleware(
			recovery.Recovery(),
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
				metrics.WithRequests(prom.NewCounter(_metricRequests)),
			),
			logging.Server(logger),
			validate.Validator(),
			selector.Server( //jwt中间件
				userutil.NewJWTServerMidware(userutil.NewJWT()),
			).Match(NewWhiteListMatcher()).Build(), ratelimit.Server(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	pb.RegisterUserServiceHTTPServer(srv, userService)
	pb.RegisterExamineAidHTTPServer(srv, ExamineAidService)
	pb.RegisterUserAidQueryHTTPServer(srv, aidQueryService)
	pb.RegisterUserAidManagerHTTPServer(srv, aidManagerService)
	pb.RegisterMutualAidQueryHTTPServer(srv, mutualAidQueryService)
	h := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", h)

	return srv
}

// jwt白名单，在名单中的路由不用校验jwt
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/api.mutualaid.UserService/SendAuthCode"] = struct{}{}
	whiteList["/api.mutualaid.UserService/Auth"] = struct{}{}
	whiteList["/api.mutualaid.UserService/ActiveUser"] = struct{}{}
	whiteList["/api.mutualaid.UserService/JSAPISign"] = struct{}{}
	whiteList["/api.mutualaid.UserService/WxLogin"] = struct{}{}
	whiteList["/api.mutualaid.UserService/WxOAuth2"] = struct{}{}

	whiteList["/api.mutualaid.ExamineAid/ExamineLogin"] = struct{}{}
	//whiteList["/api.mutualaid.MutualAidQuery/Discovery"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

func responseEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, v interface{}) error {
	codec, _ := http.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	if r.Header.Get("JWT-Token") != "" {
		w.Header().Set("JWT-Token", r.Header.Get("JWT-Token"))
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	_, err = w.Write([]byte(data))
	if err != nil {
		return err
	}
	return nil
}
