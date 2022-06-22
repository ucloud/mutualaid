package aid_test

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	biz "github.com/ucloud/mutualaid/backend/internal/biz/aid"
	"github.com/ucloud/mutualaid/backend/internal/conf"
	"github.com/ucloud/mutualaid/backend/internal/data"
	"github.com/ucloud/mutualaid/backend/internal/data/mysql"
	"github.com/ucloud/mutualaid/backend/internal/data/redis"
	"google.golang.org/protobuf/types/known/durationpb"
	"os"
	"testing"
	"time"
)

var (
	_uc  biz.UseCase
	_ctx = context.Background()

	_aid *biz.Aid
)

func init() {
	logger := log.NewStdLogger(os.Stdout)
	cnf := &conf.Data{
		Database: &conf.Data_Database{
			Driver:        "mysql",
			Source:        "root:root@tcp(127.0.0.1:3306)/mutualaid?charset=utf8&parseTime=True&loc=Local",
			SlowThreshold: durationpb.New(time.Second),
		},
		Redis: &conf.Data_Redis{
			Network:      "",
			Addrs:        []string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"},
			UserName:     "",
			Password:     "",
			PoolSize:     5,
			MinIdleConn:  1,
			ReadTimeout:  durationpb.New(time.Millisecond * 100),
			WriteTimeout: durationpb.New(time.Minute * 100),
		},
	}
	db := mysql.NewDB(cnf, logger)
	mysqlData, _, _ := mysql.NewData(db, logger)
	mysqlRepo := mysql.NewAidRepository(logger, mysqlData)
	redisData, _, _ := redis.NewData(cnf, logger)
	redisRepo := redis.NewAidRepository(logger, redisData)
	userDb := mysql.NewUserDBRepo(mysqlData)
	userCache := redis.NewUserRepository(logger, redisData)
	userRepo := data.NewUserRepo(userDb, userCache)
	repo := data.NewAidRepository(redisRepo, mysqlRepo, userRepo)
	tx := mysql.NewTransaction(mysqlData)

	_uc = biz.NewUseCase(logger, repo, tx)
}

func TestUseCase_Discovery(t *testing.T) {
	tests := []struct {
		u    uint64
		l, r float64
		want int
	}{
		{
			u:    2,
			l:    1.001,
			r:    1.001,
			want: 0,
		},
		{
			u:    0,
			l:    1.001,
			r:    1.001,
			want: 50,
		},
	}

	for _, tt := range tests {
		list, err := _uc.Discovery(_ctx, tt.u, tt.l, tt.r, 0, 100)
		assert.NoError(t, err)

		for _, l := range list {
			t.Logf("----> %+v, %v", l, l.UserInfo)
			// 非登入态脱敏
			if tt.u == 0 {
				assert.LessOrEqual(t, len(l.Content), tt.want)
			}
		}
	}
}

func TestUseCase_CreateAid(t *testing.T) {
	tests := []struct {
		p    *biz.Aid
		want int
	}{
		{
			p:    biz.NewAid(1, 10, 10, 10, "TestUseCase_CreateAid", 1.3, 1, "杨浦区1", "隆昌路", "17775208967"),
			want: 0,
		},
		{
			p:    biz.NewAid(2, 10, 10, 10, "TestUseCase_CreateAidTestUseCase_CreateAidTestUseCase_CreateAid", 1.4, 1, "杨浦区2", "隆昌路", "17775208967"),
			want: 0,
		},
		{
			p:    biz.NewAid(3, 10, 10, 10, "TestUseCase_CreateAidTestUseCase_CreateAidTestUseCase_CreateAid", 1.2, 1, "杨浦区2", "隆昌路", "17775208967"),
			want: 0,
		},
	}

	for _, tt := range tests {
		err := _uc.CreateAid(_ctx, *tt.p)
		assert.NoError(t, err)
	}
}

func TestUseCase_CancelAid(t *testing.T) {
	tests := []struct {
		p    *biz.Aid
		want int
	}{
		{
			p:    biz.NewAid(2, 10, 10, 10, "TestUseCase_CancelAid", 1, 1, "杨浦区", "隆昌路", "17775208967"),
			want: 0,
		},
	}

	for _, tt := range tests {
		err := _uc.CreateAid(_ctx, *tt.p)
		assert.NoError(t, err)
		err = _uc.CancelAid(_ctx, tt.p.ID, tt.p.UserID)
		assert.NoError(t, err)
		err = _uc.CancelAid(_ctx, tt.p.ID, tt.p.UserID)
		assert.NoError(t, err)
		err = _uc.FinishAid(_ctx, tt.p.ID, tt.p.UserID, 0)
		assert.Error(t, err)
	}
}

func TestUseCase_FinishAid(t *testing.T) {
	tests := []struct {
		p    *biz.Aid
		want int
	}{
		{
			p:    biz.NewAid(2, 10, 10, 10, "TestUseCase_CancelAid", 1, 1, "杨浦区", "隆昌路", "17775208967"),
			want: 0,
		},
	}

	for _, tt := range tests {
		err := _uc.CreateAid(_ctx, *tt.p)
		assert.NoError(t, err)
		err = _uc.FinishAid(_ctx, tt.p.ID, tt.p.UserID, 0)
		assert.NoError(t, err)
		err = _uc.FinishAid(_ctx, tt.p.ID, tt.p.UserID, 0)
		assert.NoError(t, err)
		err = _uc.CancelAid(_ctx, tt.p.ID, tt.p.UserID)
		assert.Error(t, err)
	}
}

func TestUseCase_CreateAidMessage(t *testing.T) {

	tests := []struct {
		p    *biz.Aid
		pp   *biz.Message
		want int
	}{
		{
			p:    biz.NewAid(2, 10, 10, 10, "TestUseCase_CancelAid", 1, 1, "杨浦区", "隆昌路", "17775208967"),
			pp:   biz.NewMessage(93728292166697930, 10, "17775208967", "简单精炼  不超过110字"),
			want: 0,
		},
	}

	for _, tt := range tests {
		err := _uc.CreateAid(_ctx, *tt.p)
		assert.NoError(t, err)
		tt.pp.AidID = tt.p.ID
		err = _uc.CreateAidMessage(_ctx, *tt.pp)
		assert.NoError(t, err)
		finishMsg := biz.NewMessage(tt.p.ID, 1, "17775208967", "简单精炼  不超过110字")
		err = _uc.CreateAidMessage(_ctx, *finishMsg)
		assert.NoError(t, err)
		err = _uc.FinishAid(_ctx, tt.p.ID, tt.p.UserID, finishMsg.ID)
	}
}

func TestSubString(t *testing.T) {
	a := "123你好👌😄😄"
	s := []rune(a)
	t.Log(len(a), len(s))
}
