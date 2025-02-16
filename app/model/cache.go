package model

import "time"

type Status string

const (
	StatusCreate  Status = "CREATE"  // 创建中
	StatusRunning Status = "RUNNING" // 运行中
	StatusDone    Status = "DONE"    // 完成
	StatusFailed  Status = "FAILED"  // 完成
)

type Cache struct {
	ID        string    `json:"id" gorm:"primary_key;type:varchar(64);not null;"`
	Value     string    `json:"value" gorm:"type:blob"`
	Status    Status    `json:"status" gorm:"type:varchar(64);not null;default:CREATE"`
	CreatedAt time.Time `gorm:"column:created_at;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;comment:更新时间" json:"updated_at"`
}
