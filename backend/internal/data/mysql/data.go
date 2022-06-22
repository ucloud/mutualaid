package mysql

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	biz "github.com/ucloud/mutualaid/backend/internal/biz/aid"
	"github.com/ucloud/mutualaid/backend/internal/conf"
	"github.com/ucloud/mutualaid/backend/tools/micro-mini/log/gormlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Data .
type Data struct {
	db  *gorm.DB
	log *log.Helper
}

type contextTxKey struct{}

func (d *Data) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

// NewTransaction .
func NewTransaction(d *Data) biz.Transaction {
	return d
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "transaction/data"))
	d := &Data{
		db:  db,
		log: l,
	}

	return d, func() {}, nil
}

// NewDB gorm Connecting to a Database
func NewDB(conf *conf.Data, logger log.Logger) *gorm.DB {
	db, err := gorm.Open(mysql.Open(conf.Database.GetSource()), &gorm.Config{
		Logger: gormlog.NewGormLog(log.NewHelper(logger), conf.Database.GetSlowThreshold().AsDuration()),
	})
	if err != nil {
		log.NewHelper(logger).Errorw("err", err)
	}
	return db

}
