package model

type AidWithUser struct {
	ID             uint64  `gorm:"primaryKey;autoIncrement:false" json:"id"`                        // 求助ID，非自增
	UserID         uint64  `gorm:"column:user_id;not null" json:"user_id"`                          // 用户ID
	Type           int32   `gorm:"column:type;not null" json:"type"`                                // 求助类型: 10-食品生活物资，15-就医，20-求药，25-防疫物资，30-隔离求助，35-心理援助，40-其他
	Group          int32   `gorm:"column:group;not null" json:"group"`                              // 求助人群: 10-重症患者，15-儿童婴儿，20-孕妇，25-老人，30-残障，35-外来务工人员，40-滞留人员，45-新冠阳性，50-医护工作者，55-街道社区，60-外籍人士
	EmergencyLevel int32   `gorm:"column:emergency_level;not null" json:"emergency_level"`          // 紧急程度: 1-威胁生命，2-威胁健康，3-处境困难，4-暂无危险
	ExamineStatus  int32   `gorm:"column:examine_status;not null;default:10" json:"examine_status"` // '审核状态: 10-待审核，15-审核不通过，20-审核通过',
	Status         int32   `gorm:"column:status;not null;default:10" json:"status"`                 // 求助状态: 10-已创建，15-已取消，20-已完成
	FinishUserID   uint64  `gorm:"column:finish_user_id;not null;default:0" json:"finish_user_id"`  // 完成用户ID
	FinishTime     int64   `gorm:"column:finish_time;not null;default:0" json:"finish_time"`        // 完成时间，用unix时间戳表示
	MessageCount   int32   `gorm:"column:message_count;not null;default:0" json:"message_count"`    // 留言数量，增加留言时更新
	Content        string  `gorm:"column:content;not null" json:"content"`                          // 描述
	Longitude      float64 `gorm:"column:longitude;not null;default:0.00000000" json:"longitude"`   // 坐标：经度
	Latitude       float64 `gorm:"column:latitude;not null;default:0.00000000" json:"latitude"`     // 坐标：纬度
	Phone          string  `gorm:"column:phone;not null" json:"phone"`                              // 联系电话
	District       string  `gorm:"column:district;not null;default:''" json:"district"`             // 区县，使用国家民政局标准
	Address        string  `gorm:"column:address;not null;default:''" json:"address"`               // 地址
	CreateTime     int64   `gorm:"autoCreateTime" json:"create_time"`                               // 创建时间，用unix时间戳表示
	UpdateTime     int64   `gorm:"autoUpdateTime" json:"update_time"`                               // 更新时间，用unix时间戳表示
	Version        int32   `gorm:"column:version;not null;default:1" json:"version"`                // 版本，用于乐观锁控制
	Name           string  `gorm:"column:name;not null" json:"name"`                                // 姓名
	Icon           string  `gorm:"column:icon;not null;default:''" json:"icon"`                     // 头像
}
