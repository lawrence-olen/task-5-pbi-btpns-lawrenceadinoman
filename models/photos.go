package models

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	Title    string  `json:"title"`
	Caption  string  `json:"caption"`
	PhotoUrl string  `json:"photo_url"`
	UserID   float64 `json:"user_id"`
	User     User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
