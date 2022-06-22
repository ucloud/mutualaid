// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameAidMessage = "aid_messages"

// AidMessage mapped from table <aid_messages>
type AidMessage struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement:false" json:"id"`        // 帮助ID
	AidID         uint64 `gorm:"column:aid_id;not null" json:"aid_id"`            // 求助ID
	Status        int32  `gorm:"column:status;not null;default:10" json:"status"` // 状态: 10-已创建，15-已取消，20-已完成
	UserID        uint64 `gorm:"column:user_id;not null" json:"user_id"`          // 用户ID
	UserPhone     string `gorm:"column:user_phone;not null" json:"user_phone"`    // 联系电话
	Content       string `gorm:"column:content;not null" json:"content"`
	CreateTime    int64  `gorm:"autoCreateTime" json:"create_time"`                // 创建时间，用unix时间戳表示
	UpdateTime    int64  `gorm:"autoUpdateTime" json:"update_time"`                // 更新时间，用unix时间戳表示
	Version       int32  `gorm:"column:version;not null;default:1" json:"version"` // 版本，用于乐观锁控制
	ExamineStatus int32  `gorm:"column:examine_status;not null;default:20" json:"examine_status"`
}

// TableName AidMessage's table name
func (*AidMessage) TableName() string {
	return TableNameAidMessage
}
