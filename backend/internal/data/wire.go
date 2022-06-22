package data

import (
	"github.com/google/wire"
	"github.com/ucloud/mutualaid/backend/internal/data/mysql"
	"github.com/ucloud/mutualaid/backend/internal/data/redis"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewAidRepository, mysql.NewData, mysql.NewDB, mysql.NewTransaction, mysql.NewAidRepository, mysql.NewUserDBRepo, NewUserRepo, redis.NewData, redis.NewAidRepository, redis.NewUserRepository)
