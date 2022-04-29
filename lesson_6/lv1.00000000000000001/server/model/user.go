package model

import "gorm.io/gorm"

type InputUser struct {
	ID       int
	UserName string
	Password string
}

type User struct {
	gorm.Model
	Username string
	Password string
}
