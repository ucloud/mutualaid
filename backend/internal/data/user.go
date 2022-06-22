package data

import (
	"context"
	"database/sql"

	"github.com/ucloud/mutualaid/backend/infra/userutil"
	"github.com/ucloud/mutualaid/backend/internal/biz/user"
	"github.com/ucloud/mutualaid/backend/internal/data/model"
	"github.com/ucloud/mutualaid/backend/internal/data/mysql"
	"github.com/ucloud/mutualaid/backend/internal/data/redis"
)

type userRepo struct {
	db    *mysql.UserDBRepo
	cache *redis.UserRepository
}

// NewUserRepo .
func NewUserRepo(db *mysql.UserDBRepo, cache *redis.UserRepository) user.UserRepo {
	return &userRepo{
		db:    db,
		cache: cache,
	}
}

func (r *userRepo) GetUser(ctx context.Context, u *model.User) (*model.User, error) {
	var u1 *model.User
	var err error
	if u.ID > 0 {
		u1, err = r.cache.GetUser(ctx, u.ID)
		if err == nil && u1 != nil {
			return u1, nil
		}
	}

	u1, err = r.db.GetUser(ctx, u)
	if err != nil {
		return nil, err
	}

	r.cache.AddUser(ctx, u1)

	return u1, nil
}

func (r *userRepo) ListUser(ctx context.Context, community string, status int8) ([]*model.User, error) {
	return r.db.ListUser(ctx, community, status)
}

func (r *userRepo) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	u1, err := r.db.GetUserOR(ctx, u)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		u.ID = userutil.GenID()
		err = r.db.CreateUser(ctx, u)
		if err != nil {
			return nil, err
		}
	} else {
		if u.Openid != "" {
			u1.Openid = u.Openid
		}
		if u.Unionid != "" {
			u1.Unionid = u.Unionid
		}
		if u.MpOpenid != "" {
			u1.MpOpenid = u.MpOpenid
		}
		if u.Icon != "" {
			u1.Icon = u.Icon
		}
		if u.Name != "" {
			u1.Name = u.Name
		}
		err = r.UpdateUser(ctx, u1)
		if err != nil {
			return nil, err
		}
	}

	return u1, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, u *model.User) error {
	err := r.db.UpdateUser(ctx, u)
	r.cache.DeleteUser(ctx, u)
	return err
}

func (r *userRepo) DeleteUser(ctx context.Context, u *model.User) error {
	err := r.db.DeleteUser(ctx, u)
	r.cache.DeleteUser(ctx, u)
	return err
}
