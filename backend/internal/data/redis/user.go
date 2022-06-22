package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ucloud/mutualaid/backend/internal/data/model"
	"github.com/ucloud/mutualaid/backend/tools/ulibgo/utils/json"
)

type UserRepository struct {
	log  *log.Helper
	data *Data
}

const (
	expireTime  = time.Minute * 15 // 默认缓存时间
	userInfoKey = "u:i:%d"         // 用户aid ID集合
)

func NewUserRepository(logger log.Logger, data *Data) *UserRepository {
	return &UserRepository{log: log.NewHelper(logger), data: data}
}

func (r *UserRepository) GetUser(ctx context.Context, userID uint64) (*model.User, error) {
	result := r.data.DB(ctx).Get(ctx, fmt.Sprintf(userInfoKey, userID))
	if result.Err() != nil && !errors.Is(result.Err(), redis.Nil) {
		return nil, fmt.Errorf("get user(%d) info redis failed: %s", userID, result.Err())
	}

	var u model.User
	if err := json.FromJson(result.Val(), &u); err != nil {
		return nil, fmt.Errorf("user(%d) info in redis is invalid: %+v", userID, err)
	}

	return &u, nil
}

func (r *UserRepository) AddUser(ctx context.Context, u *model.User) error {
	result := r.data.DB(ctx).Set(ctx, fmt.Sprintf(userInfoKey, u.ID), json.ToJson(u), expireTime)
	if result.Err() != nil && !errors.Is(result.Err(), redis.Nil) {
		return fmt.Errorf("set user(%d) info into redis failed: %s", u.ID, result.Err())
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, u *model.User) error {
	result := r.data.DB(ctx).Del(ctx, fmt.Sprintf(userInfoKey, u.ID), json.ToJson(u))
	if result.Err() != nil && !errors.Is(result.Err(), redis.Nil) {
		return fmt.Errorf("delete user(%d) info from redis failed: %s", u.ID, result.Err())
	}

	return nil
}
