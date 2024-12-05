package model

import "gorm.io/gorm"

type Login struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	gorm.Model
	User     string
	Password string
}

func (User) TableName() string {
	return "user"
}
