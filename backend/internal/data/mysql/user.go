package mysql

import (
	"context"
	"database/sql"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/ucloud/mutualaid/backend/internal/data/model"
	"gorm.io/gorm"
)

type UserDBRepo struct {
	*Data
}

// NewUserDBRepo .
func NewUserDBRepo(data *Data) *UserDBRepo {
	return &UserDBRepo{
		Data: data,
	}
}

func (r *UserDBRepo) GetUser(ctx context.Context, u *model.User) (*model.User, error) {
	gdb := r.db.WithContext(ctx)
	if u.ID > 0 {
		gdb = gdb.Where("id = ?", u.ID)
	} else if u.Openid != "" {
		gdb = gdb.Where("openid = ?", u.Openid)
	} else if u.MpOpenid != "" {
		gdb = gdb.Where("mp_openid = ?", u.MpOpenid)
	} else if u.Unionid != "" {
		gdb = gdb.Where("unionid = ?", u.Unionid)
	} else if u.Phone != "" {
		gdb = gdb.Where("phone = ?", u.Phone)
	}

	var u1 model.User
	result := gdb.First(&u1)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, sql.ErrNoRows
		}

		return nil, result.Error
	}

	return &u1, nil
}

func (r *UserDBRepo) GetUserOR(ctx context.Context, u *model.User) (*model.User, error) {
	gdb := r.db.WithContext(ctx)
	var u1 model.User

	szselect := `select id,phone,name,icon,openid,mp_openid,unionid,addr,community,status,create_time from user where `
	var result *gorm.DB
	if u.Openid != "" {
		result = gdb.Raw(szselect+" openid=? or unionid=?", u.Openid, u.Unionid).Scan(&u1)
	} else {
		result = gdb.Raw(szselect+" mp_openid=? or unionid=?", u.MpOpenid, u.Unionid).Scan(&u1)
	}
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, sql.ErrNoRows
		}

		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return &u1, nil
}

func (r *UserDBRepo) ListUser(ctx context.Context, community string, status int8) ([]*model.User, error) {
	var ulist []*model.User
	var cond []string
	var args []interface{}

	if community != "" {
		cond = append(cond, "community = ?")
		args = append(args, community)
	}

	if status > 0 {
		cond = append(cond, "status = ?")
		args = append(args, status)
	}
	where := strings.Join(cond, " and ")

	result := r.db.WithContext(ctx).Where(where, args...).Find(&ulist)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return ulist, nil
}

func (r *UserDBRepo) CreateUser(ctx context.Context, u *model.User) error {
	u.Status = 2
	result := r.db.WithContext(ctx).Create(u)
	return result.Error
}

func (r *UserDBRepo) UpdateUser(ctx context.Context, u *model.User) error {
	result := r.db.WithContext(ctx).Save(u)
	return result.Error
}

func (r *UserDBRepo) DeleteUser(ctx context.Context, u *model.User) error {
	result := r.db.WithContext(ctx).Delete(u)
	return result.Error
}
