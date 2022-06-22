package redis

import (
	"context"
	"github.com/go-redis/redis/v8"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/ucloud/mutualaid/backend/api/mutualaid/uerrors"
	"github.com/ucloud/mutualaid/backend/internal/conf"
	"time"
)

// common variables
const (
	defaultPoolSize    = 60
	defaultMinIdleConn = 20
	defaultMaxRetries  = 3
	defaultTimeOut     = 1 * time.Second
	defaultIdleTimeout = 5 * time.Minute
)

type Data struct {
	db  *redis.ClusterClient
	log *log.Helper
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	poolSize := c.Redis.GetPoolSize()
	if poolSize == 0 {
		poolSize = defaultPoolSize
	}
	minIdleConn := c.Redis.MinIdleConn
	if minIdleConn <= 0 {
		minIdleConn = defaultMinIdleConn
	}

	cache := redis.NewClusterClient(&redis.ClusterOptions{
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
		return nil, nil, uerrors.ErrorInfraCacheOpenError("redis open error: %v", status.Err())
	}

	cleanup := func() {
		cache.Close()
		log.NewHelper(logger).Info("redis closed")
	}

	return &Data{db: cache, log: log.NewHelper(logger)}, cleanup, nil
}

func (d *Data) DB(ctx context.Context) *redis.ClusterClient {
	return d.db.WithContext(ctx)
}
