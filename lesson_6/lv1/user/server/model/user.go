package model

import "gorm.io/gorm"

type UserInfo struct {
	gorm.Model
	Username string
	Password string
}
