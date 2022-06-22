package data

import (
	"context"
	"fmt"

	"github.com/ucloud/mutualaid/backend/internal/biz/aid"
	userBiz "github.com/ucloud/mutualaid/backend/internal/biz/user"
	"github.com/ucloud/mutualaid/backend/internal/data/model"

	"github.com/ucloud/mutualaid/backend/api/mutualaid/uerrors"
	biz "github.com/ucloud/mutualaid/backend/internal/biz/aid"
	"github.com/ucloud/mutualaid/backend/internal/data/mysql"
	"github.com/ucloud/mutualaid/backend/internal/data/redis"
)

type AidRepository struct {
	cache    *redis.AidRepository
	db       *mysql.AidRepository
	udb      *mysql.UserDBRepo
	userRepo userBiz.UserRepo
}

func NewAidRepository(cache *redis.AidRepository, db *mysql.AidRepository, udb *mysql.UserDBRepo, repo userBiz.UserRepo) biz.AidRepository {
	return &AidRepository{cache: cache, db: db, udb: udb, userRepo: repo}
}
func (r *AidRepository) GetExamineList(ctx context.Context, examineStatus int32, examineStatusOrder string, createTimeOrder string, updateTimeOrder string, pageNo int32, pageSize int32, VagueSearch string) ([]*biz.Aid, []biz.ExamineTypeMap, error) {
	// 查找结果集

	list, totalArray, err := r.db.GetExamineList(ctx, examineStatus, []string{examineStatusOrder, updateTimeOrder, createTimeOrder}, int(pageNo*pageSize), int(pageSize), VagueSearch)
	if err != nil {
		return nil, nil, err
	}

	return list, totalArray, nil
}

func (r *AidRepository) Discovery(ctx context.Context, userId uint64, latitude, longitude float64, pageNo, pageSize int32) ([]*biz.Aid, error) {
	var selectedAid []*biz.Aid
	// 查找结果集
	cachedAid, unCachedIds, err := r.cache.Discovery(ctx, userId, latitude, longitude, pageNo, pageSize)
	if err != nil {
		return nil, err
	}
	if len(cachedAid) > 0 {
		selectedAid = cachedAid
	}

	// 缓存无数据,可能被清理，重新加载。
	if len(cachedAid) == 0 && len(unCachedIds) == 0 {
		var list []*biz.Aid

		// 增加默认过滤，examine_status为20未审核成功的在Discovery中展示
		list, err = r.db.Discovery(ctx, nil, []int32{biz.StatusCreated}, userId, []int32{biz.StatusExamineFinish})
		if err != nil {
			return nil, err
		}

		for _, aid := range list {
			if err := r.cache.CreateAid(ctx, *aid); err != nil {
				return nil, err
			}
		}

		cachedAid, unCachedIds, err = r.cache.Discovery(ctx, userId, latitude, longitude, pageNo, pageSize)
		if err != nil {
			return nil, err
		}
		if len(cachedAid) > 0 {
			selectedAid = cachedAid
		}
	}

	// 未缓存的数据，重新加载。
	if len(unCachedIds) > 0 {
		ids := make([]uint64, 0, len(unCachedIds))
		for k, _ := range unCachedIds {
			ids = append(ids, k)
		}

		var list []*biz.Aid

		list, err = r.db.Discovery(ctx, ids, []int32{biz.StatusCreated}, userId, []int32{biz.StatusExamineFinish})
		if err != nil {
			return nil, err
		}

		for _, aid := range list {
			aid.Distance = int64(unCachedIds[aid.ID].Dist)

			// 非登入态数据脱敏
			if userId == 0 {
				aid.SecurityFilter()
			}

			selectedAid = append(selectedAid, aid)

			if err := r.cache.CreateAid(ctx, *aid); err != nil {
				return nil, err
			}
		}
	}

	return selectedAid, nil
}

func (r *AidRepository) CreateAid(ctx context.Context, aid biz.Aid) error {
	ok, err := r.cache.CreateAidLock(ctx, aid.UserID)
	if err != nil {
		return err
	}
	if !ok {
		return uerrors.ErrorRatelimit("Please retry later.")
	}

	// 补充用户信息
	user := &model.User{ID: aid.UserID}
	user, err = r.userRepo.GetUser(ctx, user)
	if err != nil {
		return err
	}

	if user.Status == 3 {
		return uerrors.ErrorBizUserBlock("User Is Block")
	}

	aid.UserInfo = &biz.UserInfo{
		Name: user.Name,
		ICon: user.Icon,
	}

	if err := r.db.CreateAid(ctx, aid); err != nil {
		return err
	}

	// if err := r.cache.CreateAid(ctx, aid); err != nil {
	// 	return err
	// }

	return nil
}

func (r *AidRepository) UpdateAid(ctx context.Context, aid biz.Aid, isFinal bool, isCreateCache bool) (bool, error) {
	ok, err := r.db.UpdateAid(ctx, aid)
	if err != nil {
		return false, err
	}

	if isCreateCache {
		// 创建缓存时，客户可能没有用户的基础信息，在此补上
		u, err := r.udb.GetUser(ctx, &model.User{ID: aid.UserID})

		if err != nil {
			return false, err
		}
		fmt.Println(u, aid.UserID)
		aid.UserInfo = &biz.UserInfo{
			ICon: u.Icon,
			Name: u.Name,
		}

		if err = r.cache.CreateAid(ctx, aid); err != nil {
			return false, err
		}
	}

	if isFinal {
		ok, err = r.cache.DeleteAid(ctx, aid)
		if err != nil {
			return false, err
		}
	}

	return ok, nil
}

func (r *AidRepository) ListAid(ctx context.Context, ids []uint64, userIds []uint64, status []int32, userIdNotIn uint64, offset, limit int64) ([]*biz.Aid, int64, error) {
	return r.db.ListAid(ctx, ids, userIds, status, userIdNotIn, offset, limit)
}

func (r *AidRepository) ExamineLogin(ctx context.Context, userName string, password string) ([]*aid.ExamineUser, error) {
	row, err := r.db.ListExamineUser(ctx, 0, userName, password)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, err
	}
	// 转换
	getFormatList := make([]*aid.ExamineUser, 0, len(row))
	for _, a := range row {
		thisUser := &aid.ExamineUser{
			ID:       a.ID,
			NameCn:   a.NameCn,
			UserName: a.UserName,
		}
		getFormatList = append(getFormatList, thisUser)
	}
	return getFormatList, nil
}

func (r *AidRepository) GetAid(ctx context.Context, id uint64) (*biz.Aid, error) {
	return r.db.GetAid(ctx, id)
}

func (r *AidRepository) CreateAidMessage(ctx context.Context, message biz.Message) error {
	ok, err := r.cache.CreateAidMessageLock(ctx, message.UserID)
	if err != nil {
		return err
	}
	if !ok {
		return uerrors.ErrorRatelimit("Please retry later.")
	}

	// 如果用户Id
	user := &model.User{ID: message.UserID}
	user, err = r.userRepo.GetUser(ctx, user)
	if err != nil {
		return err
	}

	if user.Status == 3 {
		return uerrors.ErrorBizUserBlock("User Is Block")
	}

	if err := r.db.CreateAidMessage(ctx, message); err != nil {
		return err
	}

	if err := r.cache.CreateAidMessage(ctx, message); err != nil {
		return err
	}

	return nil
}

func (r *AidRepository) UpdateAidMessage(ctx context.Context, message biz.Message) (bool, error) {
	ok, err := r.db.UpdateAidMessage(ctx, message)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (r *AidRepository) UpdateAllAidAndMessageExamineStatus(ctx context.Context, userId uint64, examineStatus int32) error {
	err := r.db.UpdateAllAidAndMessageExamineStatus(ctx, userId, examineStatus)
	if err != nil {
		return err
	}

	return nil
}

func (r *AidRepository) GetAidMessage(ctx context.Context, id uint64) (*biz.Message, error) {
	return r.db.GetAidMessage(ctx, id)
}

func (r *AidRepository) ListAidMessage(ctx context.Context, id uint64, userId uint64, status []int32) ([]*biz.Message, error) {
	return r.db.ListAidMessage(ctx, id, userId, status, []int32{biz.StatusExamineFinish})
}
func (r *AidRepository) GetUserName(ctx context.Context, userId uint64) string {
	u, err := r.udb.GetUser(ctx, &model.User{ID: userId})
	if err != nil {
		return ""
	}
	return u.Name
}

func (r *AidRepository) GetUserOpenID(ctx context.Context, userId uint64) string {
	u, err := r.udb.GetUser(ctx, &model.User{ID: userId})
	if err != nil {
		return ""
	}
	return u.Openid
}
