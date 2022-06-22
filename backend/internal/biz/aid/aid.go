package aid

import (
	"github.com/ucloud/mutualaid/backend/api/mutualaid/uerrors"
	"github.com/ucloud/mutualaid/backend/tools/ulibgo/utils/serials"
	"time"
)

const (
	StatusCreated = (iota + 1) * 10
	StatusFinished
	StatusCanceled      = 15
	StatusExamineWait   = 10
	StatusExamineBlock  = 15
	StatusExamineFinish = 20
	SaltString          = "VW?2y2bT2nr"
	UserIsBlock         = 3
)

type Aid struct {
	ID             uint64  // 求助ID，非自增
	UserID         uint64  // 用户ID
	Type           int32   // 求助类型: 10-食品生活物资，15-就医，20-求药，25-防疫物资，30-隔离求助，35-心理援助，40-其他
	Group          int32   // 求助人群: 10-重症患者，15-儿童婴儿，20-孕妇，25-老人，30-残障，35-外来务工人员，40-滞留人员，45-新冠阳性，50-医护工作者，55-街道社区，60-外籍人士
	EmergencyLevel int32   // 紧急程度: 1-威胁生命，2-威胁健康，3-处境困难，4-暂无危险
	Status         int32   // 求助状态: 10-已创建，15-已取消，20-已完成
	ExamineStatus  int32   // 求助状态: 10-未审核，15-审核不通过，20-审核通过
	FinishUserID   uint64  // 完成用户ID
	FinishTime     int64   // 完成时间，用unix时间戳表示
	MessageCount   int32   // 留言数量，增加留言时更新
	Content        string  // 描述
	Longitude      float64 // 坐标：经度
	Latitude       float64 // 坐标：纬度
	Phone          string  // 联系电话
	District       string  // 区县，使用国家民政局标准
	Address        string  // 地址
	CreateTime     int64   // 创建时间，用unix时间戳表示
	UpdateTime     int64   // 更新时间，用unix时间戳表示
	Version        int32   // 版本，用于乐观锁控制
	Messages       []*Message
	Distance       int64
	UserInfo       *UserInfo
}

type UserInfo struct {
	Name string
	ICon string
}

type ExamineUser struct {
	ID       uint64 // 用户ID,自增
	NameCn   string // 中文名
	UserName string // 用户名
	Password string // 密码加盐
}

type ExamineTypeMap struct {
	ExamineStatus int
	Count         int64
}

func (a *Aid) IsCanceled() bool {
	return a.Status == StatusCanceled
}

func (a *Aid) Cancel() error {
	if a.Status != StatusCreated {
		return uerrors.ErrorBizInvalidParam("aid can not be canceled")
	}
	a.Status = StatusCanceled
	a.FinishTime = time.Now().Unix()
	return nil
}

func (a *Aid) IsCreated() bool {
	return a.Status == StatusCreated
}

func (a *Aid) IsFinished() bool {
	return a.Status == StatusFinished
}

func (a *Aid) Finish() error {
	if a.Status != StatusCreated {
		return uerrors.ErrorBizInvalidParam("aid can not be finished")
	}
	for _, m := range a.Messages {
		m.Finished()
	}
	a.Status = StatusFinished
	a.FinishTime = time.Now().Unix()
	return nil
}

func (a *Aid) AddMessage(msg *Message) {
	a.Messages = append(a.Messages, msg)
	a.MessageCount += 1
	return
}

//SecurityFilter 数据脱敏：仅显示50个字
func (a *Aid) SecurityFilter() {
	maxLen := 50
	contentLen := len([]rune(a.Content))
	if contentLen < maxLen {
		maxLen = contentLen
	}

	a.Content = string([]rune(a.Content)[:maxLen])
}

func NewAid(userID uint64, typ int32, group int32, emergencyLevel int32, content string, longitude, latitude float64, district string, address, phone string) *Aid {
	return &Aid{
		ID:             serials.GenSafeID(),
		UserID:         userID,
		Type:           typ,
		Group:          group,
		EmergencyLevel: emergencyLevel,
		Status:         StatusCreated,
		Content:        content,
		Longitude:      longitude,
		Latitude:       latitude,
		Phone:          phone,
		District:       district,
		Address:        address,
		CreateTime:     time.Now().Unix(),
	}
}

type Message struct {
	ID            uint64 // 帮助ID
	AidID         uint64 // 求助ID
	Status        int32  // 状态: 10-已创建，15-已取消，20-已完成
	UserID        uint64 // 用户ID
	UserPhone     string // 联系电话
	Content       string // 帮助说明
	ExamineStatus int32  // 求助状态: 10-未审核，15-审核通过，20-审核不通过
	CreateTime    int64  // 创建时间，用unix时间戳表示
	Version       int32  // 版本，用于乐观锁控制
}

type BlocUserList struct {
	UserID     uint64 // 用户ID
	UserPhone  string // 联系电话
	Name       string // 帮助说明
	CreateTime int64  // 创建时间，用unix时间戳表示
	Version    int32  // 版本，用于乐观锁控制
}

func (m *Message) Finished() {
	m.Status = StatusFinished
}

func NewMessage(aidID uint64, userID uint64, userPhone string, content string) *Message {
	return &Message{
		ID:         serials.GenSafeID(),
		AidID:      aidID,
		Status:     StatusCreated,
		UserID:     userID,
		UserPhone:  userPhone,
		Content:    content,
		CreateTime: time.Now().Unix(),
	}
}
