package model

import "github.com/lejianwen/rustdesk-api/v2/model/custom_types"

type UserDeviceStatus int

const (
	UserDeviceStatusBound   UserDeviceStatus = 1
	UserDeviceStatusUnbound UserDeviceStatus = 2
)

type UserDevice struct {
	IdModel
	UserId       uint             `json:"user_id" gorm:"default:0;not null;index"`
	Uuid         string           `json:"uuid" gorm:"default:'';not null;index:idx_user_uuid,unique"`
	DeviceId     string           `json:"device_id" gorm:"default:'';not null;index"`
	Platform     string           `json:"platform" gorm:"default:'';not null;"`
	Ip           string           `json:"ip" gorm:"default:'';not null;"`
	Status       UserDeviceStatus `json:"status" gorm:"default:1;not null;index"`
	FirstLoginAt custom_types.AutoTime `json:"first_login_at" gorm:"type:timestamp;"`
	LastLoginAt  custom_types.AutoTime `json:"last_login_at" gorm:"type:timestamp;"`
	TimeModel
}

type UserDeviceList struct {
	List []*UserDevice `json:"list"`
	Pagination
}
