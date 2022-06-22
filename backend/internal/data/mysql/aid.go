package mysql

import (
	"context"

	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ucloud/mutualaid/backend/api/mutualaid/uerrors"
	biz "github.com/ucloud/mutualaid/backend/internal/biz/aid"
	"github.com/ucloud/mutualaid/backend/internal/data/model"
	"gorm.io/gorm"
)

type AidRepository struct {
	log  *log.Helper
	data *Data
}

func NewAidRepository(logger log.Logger, data *Data) *AidRepository {
	return &AidRepository{log: log.NewHelper(logger), data: data}
}

func (r *AidRepository) GetExamineList(ctx context.Context, examineStatus int32, orderArr []string, offset int, limit int, vagueSearch string) ([]*biz.Aid, []biz.ExamineTypeMap, error) {
	var aidWithUsers []model.AidWithUser
	var examineTypeMap []biz.ExamineTypeMap

	var cond []string
	var args []interface{}
	countArr := []int64{0, 0, 0, 0}

	if examineStatus != 0 {
		cond = append(cond, "aid.examine_status = ?")
		args = append(args, examineStatus)
	}

	if vagueSearch != "" {
		cond = append(cond, " (aid.id like ?  or  user.name like ?  or  aid.content like ?  or  aid.phone like ? ) ")
		rexString := "%" + vagueSearch + "%"
		args = append(args, rexString, rexString, rexString, rexString)
	}

	where := strings.Join(cond, " and ")

	result := r.data.DB(ctx).Table("aid").
		Select("aid.id, aid.user_id, aid.type, aid.`group`, aid.emergency_level, aid.status,aid.examine_status , aid.finish_time, "+
			"aid.message_count, aid.content, aid.longitude, aid.latitude, aid.phone, aid.district, aid.address, "+
			"aid.create_time, aid.version, user.name, user.icon,aid.update_time").
		Joins("left join user on user.id = aid.user_id").
		Order("aid.examine_status "+orderArr[0]).
		Order("aid.update_time "+orderArr[1]).
		Order("aid.create_time "+orderArr[2]).
		Limit(limit).Offset(offset).
		Where(where, args...).
		Scan(&aidWithUsers)

	if result.Error != nil {
		return nil, examineTypeMap, uerrors.ErrorInfraDbSelectError("select aid failed: %s", result.Error.Error())
	}
	if len(aidWithUsers) == 0 {
		return nil, examineTypeMap, nil
	}
	list := make([]*biz.Aid, 0, len(aidWithUsers))
	for _, aidWithUser := range aidWithUsers {
		aid := model.AidWithUserPO2DO(aidWithUser)
		list = append(list, aid)
	}

	countResult := r.data.DB(ctx).Table("aid").
		Joins("left join user on user.id = aid.user_id").
		Where(where, args...).
		Count(&countArr[0])

	if countResult.Error != nil {
		return nil, examineTypeMap, uerrors.ErrorInfraDbSelectError("select aid failed: %s", result.Error.Error())
	}

	examineStatusMap := r.data.DB(ctx).Table("aid").Select("examine_status  as ExamineStatus,count(1) as Count").Group("examine_status").Find(&examineTypeMap)
	if examineStatusMap.Error != nil {
		return nil, examineTypeMap, uerrors.ErrorInfraDbSelectError("select aid failed: %s", result.Error.Error())
	}
	examineTypeMap = append(examineTypeMap, biz.ExamineTypeMap{ExamineStatus: 0, Count: countArr[0]})

	return list, examineTypeMap, nil
}

func (r *AidRepository) Discovery(ctx context.Context, ids []uint64, status []int32, userIdNotIn uint64, examineStatus []int32) ([]*biz.Aid, error) {
	var aidWithUsers []model.AidWithUser
	var cond []string
	var args []interface{}

	if len(ids) > 0 {
		cond = append(cond, "aid.id in ?")
		args = append(args, ids)
	}
	if len(examineStatus) > 0 {
		cond = append(cond, "aid.examine_status in ?")
		args = append(args, examineStatus)
	}

	if len(status) > 0 {
		cond = append(cond, "aid.status in ?")
		args = append(args, status)
	}
	if userIdNotIn > 0 {
		cond = append(cond, "aid.user_id <> ? and aid.finish_user_id <> ?")
		args = append(args, userIdNotIn, userIdNotIn)
	}
	where := strings.Join(cond, " and ")

	result := r.data.DB(ctx).Table("aid").
		Select("aid.id, aid.user_id, aid.type, aid.`group`, aid.emergency_level, aid.status,aid.examine_status, aid.finish_time, "+
			"aid.message_count, aid.content, aid.longitude, aid.latitude, aid.phone, aid.district, aid.address, "+
			"aid.create_time, aid.version, user.name, user.icon").
		Joins("left join user on user.id = aid.user_id").
		Where(where, args...).
		Scan(&aidWithUsers)
	if result.Error != nil {
		return nil, uerrors.ErrorInfraDbSelectError("select aid failed: %s", result.Error.Error())
	}
	if len(aidWithUsers) == 0 {
		return nil, nil
	}

	list := make([]*biz.Aid, 0, len(aidWithUsers))
	for _, aidWithUser := range aidWithUsers {
		aid := model.AidWithUserPO2DO(aidWithUser)
		list = append(list, aid)
	}

	return list, nil
}

func (r *AidRepository) CreateAid(ctx context.Context, aid biz.Aid) error {
	Aid := model.AidDO2PO(aid)

	return r.data.DB(ctx).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Create(Aid)
		return nil
	})
}

func (r *AidRepository) UpdateAid(ctx context.Context, aid biz.Aid) (bool, error) {
	po := model.AidDO2PO(aid)

	result := r.data.DB(ctx).Model(po).Where("version = ?", po.Version).Updates(map[string]interface{}{
		"message_count":  po.MessageCount,
		"finish_user_id": po.FinishUserID,
		"finish_time":    po.FinishTime,
		"examine_status": po.ExamineStatus,
		"status":         po.Status,
		"version":        gorm.Expr("version + ?", 1),
	})
	if result.Error != nil {
		return false, uerrors.ErrorInfraDbUpdateError("update aid failed: %s", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r *AidRepository) ListAid(ctx context.Context, ids []uint64, userIds []uint64, status []int32, userIdNotIn uint64, offset, limit int64) ([]*biz.Aid, int64, error) {
	var list []*model.Aid

	prepareDB := func() *gorm.DB {
		db := r.data.DB(ctx).WithContext(ctx)
		if len(ids) > 0 {
			db = db.Where("id IN ?", ids)
		}

		if len(userIds) > 0 {
			db = db.Where("user_id IN ?", userIds)
		}

		if userIdNotIn > 0 {
			db = db.Where("user_id <> ?", userIdNotIn)
		}

		if len(status) > 0 {
			db = db.Where("status IN ?", status)
		}

		return db
	}

	db := prepareDB()
	if limit > 0 {
		// get count
		db = db.Offset(int(offset)).Limit(int(limit))
	}
	result := db.Order("create_time DESC").
		Find(&list)
	if result.Error != nil {
		return nil, 0, uerrors.ErrorInfraDbSelectError("select aid failed: %s", result.Error.Error())
	}

	var bizList []*biz.Aid
	for _, a := range list {
		bizList = append(bizList, model.AidPO2DO(*a))
	}

	total := int64(len(bizList))
	if limit > 0 {
		// get count
		db := prepareDB()
		db.Model(&model.Aid{}).Count(&total)
	}

	return bizList, total, nil
}

func (r *AidRepository) GetAid(ctx context.Context, id uint64) (*biz.Aid, error) {
	var po model.Aid
	result := r.data.DB(ctx).Where("id = ?", id).First(&po)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, uerrors.ErrorInfraDbSelectError("select aid failed: %s", result.Error.Error())
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return model.AidPO2DO(po), nil
}

func (r *AidRepository) CreateAidMessage(ctx context.Context, message biz.Message) error {
	msg := model.AidMessageDO2PO(message)
	err := r.data.DB(ctx).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Create(msg)
		return result.Error
	})
	return err
}

func (r *AidRepository) UpdateAidMessage(ctx context.Context, message biz.Message) (bool, error) {
	po := model.AidMessageDO2PO(message)
	result := r.data.DB(ctx).Model(po).Where("version = ?", po.Version).Updates(map[string]interface{}{
		"status":  po.Status,
		"version": gorm.Expr("version + ?", 1),
	})
	if result.Error != nil {
		return false, uerrors.ErrorInfraDbUpdateError("update aid message failed: %s", result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r *AidRepository) GetAidMessage(ctx context.Context, id uint64) (*biz.Message, error) {
	var po model.AidMessage
	result := r.data.DB(ctx).Where("id = ?", id).First(&po)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, uerrors.ErrorInfraDbSelectError("select aid message failed: %s", result.Error.Error())
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return model.AidMessagePO2DO(po), nil
}

func (r *AidRepository) ListAidMessage(ctx context.Context, aidID uint64, userId uint64, status []int32, examineStatus []int32) ([]*biz.Message, error) {
	var list []*model.AidMessage

	db := r.data.DB(ctx).WithContext(ctx)

	if aidID > 0 {
		db = db.Where("aid_id=?", aidID)
	}

	if userId > 0 {
		db = db.Where("user_id=?", userId)
	}

	if len(status) > 0 {
		db = db.Where("status IN ?", status)
	}

	if len(examineStatus) > 0 {
		db = db.Where("examine_status IN ?", examineStatus)
	}

	result := db.Find(&list)
	if result.Error != nil {
		return nil, uerrors.ErrorInfraDbSelectError("select aid_messages failed: %s", result.Error.Error())
	}

	var bizList []*biz.Message
	for _, m := range list {
		bizList = append(bizList, model.AidMessagePO2DO(*m))
	}
	return bizList, nil
}

func (r *AidRepository) ListExamineUser(ctx context.Context, uId uint64, userName string, password string) ([]*model.ExamineUser, error) {
	var uList []*model.ExamineUser

	db := r.data.DB(ctx).WithContext(ctx)

	if uId > 0 {
		db = db.Where("id=?", uId)
	}

	if password != "" {
		db = db.Where("user_name=?", userName)
		db = db.Where("password = ?", password)
	}

	result := db.Find(&uList)
	if result.Error != nil {
		return nil, uerrors.ErrorInfraDbSelectError("select examine user failed: %s", result.Error.Error())
	}

	return uList, nil
}

func (r *AidRepository) UpdateAllAidAndMessageExamineStatus(ctx context.Context, userId uint64, examineStatus int32) error {
	messagePo := &model.AidMessage{}
	aidPo := &model.Aid{}

	// 如果后续用户恢复，旧消息不恢复
	if examineStatus != biz.StatusExamineWait {
		messageResult := r.data.DB(ctx).Model(messagePo).Where(" user_id = ?", userId).Updates(map[string]interface{}{
			"examine_status": examineStatus,
			"version":        gorm.Expr("version + ?", 1),
		})
		if messageResult.Error != nil {
			return uerrors.ErrorInfraDbUpdateError("update aid message failed: %s", messageResult.Error.Error())
		}
	}

	aidResult := r.data.DB(ctx).Model(aidPo).Where(" user_id = ?", userId).Updates(map[string]interface{}{
		"examine_status": examineStatus,
		"version":        gorm.Expr("version + ?", 1),
	})
	if aidResult.Error != nil {
		return uerrors.ErrorInfraDbUpdateError("update aid message failed: %s", aidResult.Error.Error())
	}

	return nil
}
