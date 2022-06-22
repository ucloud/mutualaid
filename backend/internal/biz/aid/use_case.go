package aid

import (
	"context"
	"crypto/md5"
	"sort"

	"encoding/hex"
	userutil2 "github.com/ucloud/mutualaid/backend/infra/userutil"

	"fmt"
	"github.com/ucloud/mutualaid/backend/infra/userutil"
	"github.com/ucloud/mutualaid/backend/internal/conf"
	"github.com/ucloud/mutualaid/backend/internal/proxy/wechat"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/ucloud/mutualaid/backend/api/mutualaid/uerrors"
	"github.com/ucloud/mutualaid/backend/internal/biz"
)

type UseCase interface {
	Discovery(ctx context.Context, userId uint64, latitude, longitude float64, pageNo, pageSize int32) ([]*Aid, error)
	ListUserAid(ctx context.Context, userId uint64, offered bool, offset, limit int64) ([]*Aid, int64, error)

	CreateAid(ctx context.Context, aid Aid) error
	CancelAid(ctx context.Context, id uint64, userId uint64) error
	FinishAid(ctx context.Context, id uint64, userId uint64, messageId uint64) error

	// GetAidDetail 返回求助信息详情，以及信息和当前用户的关系(是自己发布的，还是自己帮助过的)
	GetAidDetail(ctx context.Context, id uint64, userId uint64) (a *Aid, isMyAid bool, isMyHelp bool, err error)

	CreateAidMessage(ctx context.Context, message Message) error

	// GetExamineList 审核相关接口
	GetExamineList(ctx context.Context, examineStatus int32, examineStatusOrder string, createTimeOrder string, updateTimeOrder string, pageNo int32, pageSize int32, vagueSearch string) ([]*Aid, []ExamineTypeMap, error)
	ExamineAid(ctx context.Context, id uint64, examineSction string) error
	ExamineLogin(ctx context.Context, userName string, password string) (string, error)
	UpdateAllAidAndMessageExamineStatus(ctx context.Context, userId uint64, examineStatus int32) error
}

type AidRepository interface {
	Discovery(ctx context.Context, userId uint64, latitude, longitude float64, pageNo, pageSize int32) ([]*Aid, error)
	CreateAid(ctx context.Context, aid Aid) error
	UpdateAid(ctx context.Context, aid Aid, isFinal bool, isCreateCache bool) (bool, error)
	ListAid(ctx context.Context, ids []uint64, userIds []uint64, status []int32, userIdNotIn uint64, offset, limit int64) ([]*Aid, int64, error)

	GetAid(ctx context.Context, id uint64) (*Aid, error)

	CreateAidMessage(ctx context.Context, message Message) error
	UpdateAidMessage(ctx context.Context, message Message) (bool, error)
	GetAidMessage(ctx context.Context, id uint64) (*Message, error)
	ListAidMessage(ctx context.Context, aidID uint64, userId uint64, status []int32) ([]*Message, error)

	ExamineLogin(ctx context.Context, userName string, password string) ([]*ExamineUser, error)

	GetExamineList(ctx context.Context, examineStatus int32, examineStatusOrder string, createTimeOrder string, updateTimeOrder string, pageNo int32, pageSize int32, vagueSearch string) ([]*Aid, []ExamineTypeMap, error)
	UpdateAllAidAndMessageExamineStatus(ctx context.Context, userId uint64, examineStatus int32) error
	GetUserName(ctx context.Context, userId uint64) string
	GetUserOpenID(ctx context.Context, userId uint64) string
}

type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}

type useCase struct {
	log  *log.Helper
	repo AidRepository
	tx   Transaction
	tpl  map[string]*conf.TPLArgs
}

func newUseCase(logger log.Logger, repo AidRepository, tpl map[string]*conf.TPLArgs, tx Transaction) *useCase {
	return &useCase{
		log:  log.NewHelper(logger),
		repo: repo,
		tx:   tx,
		tpl:  tpl,
	}
}

func NewUseCase(logger log.Logger, repo AidRepository, tpl map[string]*conf.TPLArgs, tx Transaction) UseCase {
	return newUseCase(logger, repo, tpl, tx)
}

func (u *useCase) Discovery(ctx context.Context, userId uint64, latitude, longitude float64, pageNo, pageSize int32) ([]*Aid, error) {
	var list []*Aid
	// 非登入态判断: 非登入态最多查3页，每页最多查10条
	if userId == 0 {
		// 每页只能查10条
		if pageSize >= 10 {
			pageSize = 10
		}
		// 最大查询3页
		if pageNo >= 3 {
			pageNo = 2
			return list, nil
		}
	}

	list, err := u.repo.Discovery(ctx, userId, latitude, longitude, pageNo, pageSize)
	if err != nil {
		return nil, err
	}

	// 按距离升序、时间降序排序
	sort.Slice(list, func(i, j int) bool {
		if list[i].Distance < list[j].Distance || (list[i].Distance == list[j].Distance && list[i].CreateTime > list[j].CreateTime) {
			return true
		}
		return false
	})

	return list, err
}

func (u *useCase) GetExamineList(ctx context.Context, examineStatus int32, examineStatusOrder string, createTimeOrder string, updateTimeOrder string, pageNo, pageSize int32, vagueSearch string) ([]*Aid, []ExamineTypeMap, error) {
	var list []*Aid
	list, countArr, err := u.repo.GetExamineList(ctx, examineStatus, examineStatusOrder, createTimeOrder, updateTimeOrder, pageNo, pageSize, vagueSearch)

	// 增加对其它计数的统计

	return list, countArr, err
}

func (u *useCase) UpdateAllAidAndMessageExamineStatus(ctx context.Context, userId uint64, examineStatus int32) error {
	// 增加对其它计数的统计
	err := u.repo.UpdateAllAidAndMessageExamineStatus(ctx, userId, examineStatus)

	if err != nil {
		return err
	}

	return nil
}

func (u *useCase) ExamineLogin(ctx context.Context, userName string, password string) (string, error) {
	// 确定用户名密码是否正确
	password = SaltString + password

	data := []byte(password)
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)
	examineUserInfo, err := u.repo.ExamineLogin(ctx, userName, hex.EncodeToString(cipherStr))
	if err != nil {
		return "", err
	}

	if len(examineUserInfo) == 0 {
		return "", uerrors.ErrorBizLoginFail("verification user password fail ")
	}

	// 得到JWTToken
	jwtToken, err := userutil2.NewJWT().Auth(examineUserInfo[0].ID, examineUserInfo[0].UserName)

	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (u *useCase) ExamineAid(ctx context.Context, id uint64, examineSction string) error {
	// 取消息
	aid, err := u.repo.GetAid(ctx, id)

	if err != nil {
		return err
	}

	if aid == nil {
		return uerrors.ErrorBadRequest("aid not found: %d", id)
	}

	// 得到对应的审核值
	// 待审核：通过、不通过
	// 审核通过,无可操作状态
	// 审核不通过：撤回到待审核
	// ['PASS','BLOCK','UNBLOCK']
	waitIngNextStepMap := map[string]int32{
		"PASS":  StatusExamineFinish,
		"BLOCK": StatusExamineBlock,
	}
	blockNextStepMap := map[string]int32{
		"UNBLOCK": StatusExamineWait,
	}

	nextStepMap := make(map[int]map[string]int32)
	nextStepMap[StatusExamineWait] = waitIngNextStepMap
	nextStepMap[StatusExamineBlock] = blockNextStepMap

	nextExamineStatus, ok := nextStepMap[int(aid.ExamineStatus)][examineSction]
	aid.ExamineStatus = nextExamineStatus

	if !ok {
		return uerrors.ErrorBizErrorExamineStep("Error Next Step For %s", examineSction)
	}

	// 执行审核值的更新,不删除cache，如果是审核完成，就生成缓存
	_, err = u.repo.UpdateAid(ctx, *aid, false, nextExamineStatus == StatusExamineFinish)
	if err != nil {
		return err
	}

	return nil
}

func (u *useCase) ListUserAid(ctx context.Context, userId uint64, offered bool, offset, limit int64) ([]*Aid, int64, error) {
	if !offered {
		// 我的求助
		list, total, err := u.repo.ListAid(ctx, nil, []uint64{userId}, nil, 0, offset, limit)
		return list, total, err
	}

	// 我的帮助
	msgList, err := u.repo.ListAidMessage(ctx, 0, userId, nil)
	if err != nil {
		return nil, 0, err
	}

	// 因为一个人可能在一个求助下多次回复, 所以这里进行一次去重
	idsOfAid := make(map[uint64]struct{})
	for _, m := range msgList {
		idsOfAid[m.AidID] = struct{}{}
	}
	ids := make([]uint64, 0, len(idsOfAid))
	for id := range idsOfAid {
		ids = append(ids, id)
	}

	// 手工处理分页; 因为前一步已经取出了所有数据
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	begin := offset
	if begin > int64(len(ids)) {
		begin = int64(len(ids))
	}
	end := offset + limit
	if end > int64(len(ids)) {
		end = int64(len(ids))
	}

	total := int64(len(idsOfAid))
	aidIDList := ids[begin:end]
	if len(aidIDList) == 0 {
		return nil, 0, nil
	}
	list, _, err := u.repo.ListAid(ctx, ids[begin:end], nil, nil, 0, 0, 0)

	return list, total, err
}

func (u *useCase) CreateAid(ctx context.Context, aid Aid) error {
	// 对发表的内容进行审核
	if err := wechat.MessageSecCheck(userutil.ExtractOpenID(ctx), aid.Content); err != nil {
		return err
	}

	return u.repo.CreateAid(ctx, aid)
}

func (u *useCase) CancelAid(ctx context.Context, id uint64, userId uint64) error {
	aid, err := u.repo.GetAid(ctx, id)
	if err != nil {
		return err
	}
	if aid == nil {
		return uerrors.ErrorBadRequest("aid not found: %d", id)
	}

	if userId != aid.UserID {
		return uerrors.ErrorPermissionError("user forbidden")
	}

	if aid.IsCanceled() {
		return nil
	}

	if err := aid.Cancel(); err != nil {
		return err
	}

	_, err = u.repo.UpdateAid(ctx, *aid, true, false)
	if err != nil {
		return err
	}

	return nil
}

func (u *useCase) FinishAid(ctx context.Context, id uint64, userId uint64, messageId uint64) error {
	aid, err := u.repo.GetAid(ctx, id)
	if err != nil {
		return err
	}
	if aid == nil {
		return uerrors.ErrorBadRequest("aid not found: %d", id)
	}

	if userId != aid.UserID {
		return uerrors.ErrorPermissionError("user forbidden")
	}

	if aid.IsFinished() {
		return nil
	}

	if messageId > 0 {
		msg, err := u.repo.GetAidMessage(ctx, messageId)
		if err != nil {
			return err
		}
		aid.Messages = append(aid.Messages, msg)
		aid.FinishUserID = msg.UserID
	}

	if err := aid.Finish(); err != nil {
		return err
	}

	return u.tx.ExecTx(ctx, func(ctx context.Context) error {
		ok, err := u.repo.UpdateAid(ctx, *aid, true, false)
		if err != nil {
			return err
		}
		if !ok {
			return uerrors.ErrorBizInvalidParam("aid updated")
		}

		var mcontent string
		var muid uint64
		for _, m := range aid.Messages {
			if m != nil {
				_, err := u.repo.UpdateAidMessage(ctx, *m)
				if err != nil {
					return err
				}
				mcontent = m.Content
				muid = m.UserID
			}
		}

		tplName := "帮助被采纳提醒"
		tpargs := u.tpl[tplName].Args
		args := map[string]wechat.TplItem{
			"求助编号": {tpargs["求助编号"], fmt.Sprint(aid.ID)},
			"求助内容": {tpargs["求助内容"], aid.Content},
			"回复内容": {tpargs["回复内容"], mcontent},
			"采纳时间": {tpargs["采纳时间"], time.Now().Format("2006年01月02日 15:04")},
		}

		go func() {
			defer func() {
				if err := recover(); err != nil {
					u.log.Error(err)
				}
			}()
			wechat.SubscribeMsgSend(u.repo.GetUserOpenID(ctx, muid), tplName, u.tpl[tplName].Id, args)
		}()

		return nil
	})

}

func (u *useCase) GetAidDetail(ctx context.Context, id uint64, userId uint64) (*Aid, bool, bool, error) {
	a, err := u.repo.GetAid(ctx, id)
	if err != nil {
		return nil, false, false, err
	}

	isMyAid := a.UserID == userId
	isMyHelp := false

	msgList, err := u.repo.ListAidMessage(ctx, a.ID, 0, nil)
	if err != nil {
		return nil, false, false, err
	}

	myhelpcount := 0
	for _, m := range msgList {
		if !isMyAid && m.UserID != userId {
			m.UserPhone = biz.MaskPhone(m.UserPhone)
		}
		a.Messages = append(a.Messages, m)
		if m.UserID == userId {
			myhelpcount++
		}
	}
	isMyHelp = myhelpcount > 0

	return a, isMyAid, isMyHelp, nil
}

func (u *useCase) CreateAidMessage(ctx context.Context, message Message) error {
	// 对发表的内容进行审核
	if err := wechat.MessageSecCheck(userutil.ExtractOpenID(ctx), message.Content); err != nil {
		return err
	}

	aid, err := u.repo.GetAid(ctx, message.AidID)
	if err != nil {
		return err
	}
	if aid == nil {
		return uerrors.ErrorBadRequest("aid not found: %d", message.AidID)
	}
	if !aid.IsCreated() {
		return uerrors.ErrorBadRequest("aid's state is err, state: %d", aid.Status)
	}
	aid.AddMessage(&message)

	return u.tx.ExecTx(ctx, func(ctx context.Context) error {
		err = u.repo.CreateAidMessage(ctx, message)
		if err != nil {
			return err
		}
		ok, err := u.repo.UpdateAid(ctx, *aid, false, false)
		if err != nil {
			return err
		}
		if !ok {
			return uerrors.ErrorBizInvalidParam("aid updated")
		}

		tplName := "收到帮助提醒"
		tpargs := u.tpl[tplName].Args
		args := map[string]wechat.TplItem{
			"求助编号": {tpargs["求助编号"], fmt.Sprint(message.AidID)},
			"求助内容": {tpargs["求助内容"], aid.Content},
			"回复人":  {tpargs["回复人"], u.repo.GetUserName(ctx, message.UserID)},
			"回复内容": {tpargs["回复内容"], message.Content},
			"回复时间": {tpargs["回复时间"], time.Now().Format("2006年01月02日 15:04")},
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					u.log.Error(err)
				}
			}()
			wechat.SubscribeMsgSend(u.repo.GetUserOpenID(ctx, aid.UserID), tplName, u.tpl[tplName].Id, args)
		}()

		return nil
	})
}
