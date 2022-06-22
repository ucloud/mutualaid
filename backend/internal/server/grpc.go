package server

import (
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	pb "github.com/ucloud/mutualaid/backend/api/mutualaid"
	"github.com/ucloud/mutualaid/backend/internal/conf"
	"github.com/ucloud/mutualaid/backend/internal/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Server,
	userService *service.UserService,
	mutualAidQueryService *service.MutualAidQueryService,
	aidQueryService *service.UserAidQueryService,
	examineAidService *service.ExamineAidService,
	aidManagerService *service.UserAidManagerService,
	logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
				metrics.WithRequests(prom.NewCounter(_metricRequests)),
			),
			logging.Server(logger),
			validate.Validator(),
			ratelimit.Server(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	pb.RegisterUserServiceServer(srv, userService)
	pb.RegisterExamineAidServer(srv, examineAidService)
	pb.RegisterMutualAidQueryServer(srv, mutualAidQueryService)
	pb.RegisterUserAidQueryServer(srv, aidQueryService)
	pb.RegisterUserAidManagerServer(srv, aidManagerService)
	return srv
}
