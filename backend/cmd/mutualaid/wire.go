//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	aidBiz "github.com/ucloud/mutualaid/backend/internal/biz/aid"
	userBiz "github.com/ucloud/mutualaid/backend/internal/biz/user"
	"github.com/ucloud/mutualaid/backend/internal/data"
	//wechatProxy "github.com/ucloud/mutualaid/backend/internal/proxy/wechat"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/ucloud/mutualaid/backend/internal/conf"
	"github.com/ucloud/mutualaid/backend/internal/server"
	"github.com/ucloud/mutualaid/backend/internal/service"
)

// initApp init kratos application.
func initApp(
	*conf.Server,
	*conf.Data,
	*conf.Proxy,
	*conf.BizConfig,
	map[string]*conf.TPLArgs,
	log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		service.ProviderSet,
		data.ProviderSet,
		//wechatProxy.ProviderSet,
		userBiz.ProviderSet,
		aidBiz.ProviderSet,
		newApp))
}
