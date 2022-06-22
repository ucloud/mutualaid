package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ucloud/mutualaid/backend/api/mutualaid/uerrors"
	biz "github.com/ucloud/mutualaid/backend/internal/biz/aid"
	"github.com/ucloud/mutualaid/backend/tools/ulibgo/utils/json"
	"strconv"
	"time"
)

const (
	nearbyRadiusM     = 100000          // 附近距离半径（米）：100km
	defaultExpireTime = time.Minute * 5 // 默认缓存时间
	aidGEOCacheKey    = "a:geo"         // aid位置集合
	userAidIdSetKey   = "u:%d:a"        // 用户自建或处理过的aid ID集合
	aidCacheKey       = "a:%v"          // aid缓存
	msgCacheKey       = "m:%d"
	aidCreateLockKey  = "a:lock:%d" // aid创建锁
	msgCreateLockKey  = "m:lock:%d" // aid留言创建锁
)

var getAidLocationName = func(id uint64) string { return fmt.Sprintf("%d", id) }
var getAidCacheKey = func(id interface{}) string { return fmt.Sprintf(aidCacheKey, id) }
var getMsgCacheKey = func(id uint64) string { return fmt.Sprintf(msgCacheKey, id) }
var getUserAidIdSetKey = func(id uint64) string { return fmt.Sprintf(userAidIdSetKey, id) }
var getAidCreateLockKey = func(id uint64) string { return fmt.Sprintf(aidCreateLockKey, id) }
var getMsgCreateLockKey = func(id uint64) string { return fmt.Sprintf(msgCreateLockKey, id) }

type AidRepository struct {
	log  *log.Helper
	data *Data
}

func NewAidRepository(logger log.Logger, data *Data) *AidRepository {
	return &AidRepository{log: log.NewHelper(logger), data: data}
}

func (r *AidRepository) lock(ctx context.Context, key string, expire time.Duration) (bool, error) {
	status := r.data.DB(ctx).SetNX(ctx, key, time.Now().Unix(), expire)
	if status.Err() != nil {
		return false, uerrors.ErrorInfraCacheSetError("setnx failed: %d", status.Err())
	}
	return status.Val(), nil
}

func (r *AidRepository) CreateAidLock(ctx context.Context, userId uint64) (bool, error) {
	return r.lock(ctx, getAidCreateLockKey(userId), time.Second*10)
}

func (r *AidRepository) CreateAidMessageLock(ctx context.Context, userId uint64) (bool, error) {
	return r.lock(ctx, getMsgCreateLockKey(userId), time.Second*10)
}

type AidDist struct {
	Id   uint64
	Dist float64
}

func (r *AidRepository) Discovery(ctx context.Context, userId uint64, latitude, longitude float64, pageNo, pageSize int32) ([]*biz.Aid, map[uint64]*AidDist, error) {
	var selectedAidSlice []*AidDist

	// 1. 根据经纬度检索附近的数据
	status := r.data.DB(ctx).GeoRadius(ctx, aidGEOCacheKey, longitude, latitude, &redis.GeoRadiusQuery{
		Radius:   nearbyRadiusM,
		Unit:     "m",
		Sort:     "ASC",
		WithDist: true,
	})
	if status.Err() != nil {
		return nil, nil, uerrors.ErrorInfraCacheSetError("geo compute failed: %s", status.Err())
	}
	nearBySet := status.Val()
	r.log.Debugf("GeoRadius nearby aid: %d", len(nearBySet))

	// 2. 若无附近数据，则返回。
	if len(nearBySet) == 0 {
		return nil, nil, nil
	}

	// 3. 获取需要过滤的数据ID
	membersMap := r.data.DB(ctx).SMembersMap(ctx, getUserAidIdSetKey(userId))
	if membersMap.Err() != nil {
		return nil, nil, uerrors.ErrorInfraCacheGetError("get user aid id cache error: %s", membersMap.Err())
	}
	userAidIdSet, _ := membersMap.Result()
	r.log.Debugf("userAidIdSet: %d", len(userAidIdSet))

	for _, z := range nearBySet {
		// 过滤当前用户数据
		if _, ok := userAidIdSet[z.Name]; ok {
			continue
		}

		id, err := strconv.ParseUint(z.Name, 10, 64)
		if err != nil {
			return nil, nil, uerrors.ErrorInfraCacheGetError("aid cache format error: %s", err.Error())
		}
		p := &AidDist{Id: id, Dist: z.Dist}
		selectedAidSlice = append(selectedAidSlice, p)
	}

	// 分页查询
	start := int(pageNo * pageSize)
	end := int(pageNo*pageSize + pageSize)
	if start >= len(selectedAidSlice) {
		return nil, nil, nil
	}

	if end >= len(selectedAidSlice) {
		end = len(selectedAidSlice)
	}
	r.log.Debugf("page query %d: %d, %d", len(userAidIdSet), start, end)

	// 4. 批量获取缓存的数据
	cachedAid := make([]*biz.Aid, 0)
	uncachedAid := make(map[uint64]*AidDist, 0)
	for _, aidDist := range selectedAidSlice[start:end] {
		result := r.data.DB(ctx).Get(ctx, getAidCacheKey(aidDist.Id))
		if result.Err() != nil && !errors.Is(result.Err(), redis.Nil) {
			return nil, nil, uerrors.ErrorInfraCacheGetError("get aid cache failed: %s", result.Err())
		}

		// 未缓存
		if errors.Is(result.Err(), redis.Nil) {
			uncachedAid[aidDist.Id] = aidDist
			continue
		}

		aid := &biz.Aid{}
		if err := json.FromJson(result.Val(), &aid); err != nil {
			return nil, nil, uerrors.ErrorInfraCacheGetError("aid cache format error: %s", result.Err())
		}

		// 非登入态数据脱敏
		if userId == 0 {
			aid.SecurityFilter()
		}

		aid.Distance = int64(aidDist.Dist)
		cachedAid = append(cachedAid, aid)
	}
	r.log.Debugf("cachedAid %d, uncachedAid %d", len(cachedAid), len(uncachedAid))

	return cachedAid, uncachedAid, nil
}

func (r *AidRepository) CreateAid(ctx context.Context, aid biz.Aid) error {
	status := r.data.DB(ctx).GeoAdd(ctx, aidGEOCacheKey, &redis.GeoLocation{
		Name:      getAidLocationName(aid.ID),
		Longitude: aid.Longitude,
		Latitude:  aid.Latitude,
	})
	if status.Err() != nil {
		return uerrors.ErrorInfraCacheSetError("geoadd failed: %d", status.Err())
	}

	if status := r.data.DB(ctx).Set(ctx, getAidCacheKey(aid.ID), json.ToJson(aid), defaultExpireTime); status.Err() != nil {
		return uerrors.ErrorInfraCacheSetError("set aid cache failed: %d", status.Err())
	}

	return nil
}

func (r *AidRepository) DeleteAid(ctx context.Context, aid biz.Aid) (bool, error) {
	status := r.data.DB(ctx).ZRem(ctx, aidGEOCacheKey, getAidLocationName(aid.ID))
	if status.Err() != nil {
		return false, uerrors.ErrorInfraCacheDeleteError("zrem failed: %d", status.Err())
	}

	if status := r.data.DB(ctx).SRem(ctx, getUserAidIdSetKey(aid.UserID), aid.ID, 0); status.Err() != nil {
		return false, uerrors.ErrorInfraCacheDeleteError("srem failed: %d", status.Err())
	}

	return true, nil
}

func (r *AidRepository) ListAid(ctx context.Context, ids []uint64, userIds []uint64, status []int32, filter func(aid *biz.Aid) bool) ([]*biz.Aid, error) {

	return nil, nil
}

func (r *AidRepository) GetAid(ctx context.Context, id uint64) (*biz.Aid, error) {
	panic("implement me")
}

func (r *AidRepository) CreateAidMessage(ctx context.Context, message biz.Message) error {
	if status := r.data.DB(ctx).Set(ctx, getMsgCacheKey(message.ID), json.ToJson(message), defaultExpireTime); status.Err() != nil {
		return uerrors.ErrorInfraCacheSetError("set message cache failed: %d", status.Err())
	}

	// 记录用户已处理过的数据
	if status := r.data.DB(ctx).SAdd(ctx, getUserAidIdSetKey(message.UserID), message.AidID, 0); status.Err() != nil {
		return uerrors.ErrorInfraCacheSetError("set user aid cache failed: %d", status.Err())
	}

	return nil
}

func (r *AidRepository) UpdateAidMessage(ctx context.Context, message biz.Message) (bool, error) {
	panic("implement me")
}

func (r *AidRepository) GetAidMessage(ctx context.Context, id uint64) (*biz.Message, error) {
	panic("implement me")
}

func (r *AidRepository) ListAidMessage(ctx context.Context, id uint64) ([]*biz.Message, error) {
	panic("implement me")
}
