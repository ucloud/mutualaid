package main

import (
	"flag"
	"os"

	"github.com/ucloud/mutualaid/backend/tools/micro-mini/log/zap"
	"github.com/ucloud/mutualaid/backend/tools/micro-mini/middleware/tracing"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/ucloud/mutualaid/backend/infra/logx"
	"github.com/ucloud/mutualaid/backend/internal/conf"
	"github.com/ucloud/mutualaid/backend/internal/proxy/wechat"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
	// path of log file
	logPath string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs/config.yaml", "config path, eg: -conf config.yaml")
	flag.StringVar(&logPath, "log", "./api.log", "log file path, eg: -log /var/log/mutualaid/api.log")
}

func newApp(
	logger log.Logger,
	hs *http.Server,
	gs *grpc.Server,
	//ps *producer.Server,
	//cs []*consumer.Server,
) *kratos.App {
	servers := make([]transport.Server, 0, 20)
	servers = append(servers, hs)
	servers = append(servers, gs)
	//servers = append(servers, ps)

	//for _, c := range cs {
	//	servers = append(servers, c)
	//}

	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(servers...),
	)
}

func main() {
	flag.Parse()

	zlog := zap.NewLogger(logx.GetLogger(logPath))
	defer zlog.Sync()

	logger := log.With(zlog,
		"caller", zap.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	wechat.InitWXBackend(bc.Proxy, bc.Data, logger)

	app, cleanup, err := initApp(bc.Server, bc.Data, bc.Proxy, bc.BizConfig, bc.Proxy.Wxkey.MsgTpl, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
