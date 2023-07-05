package domain

import (
	"gin-app/pkg/orm"

	"gorm.io/gorm"
)

// User 学生表模型
type User struct {
	gorm.Model
	Name     string `gorm:"size:255;uniqueIndex;comment:name"`
	Email    string `gorm:"size:255;uniqueIndex;comment:email"`
	Password string `gorm:"size:255;comment:password"`
}

func (s User) TableName() string {
	return "users"
}

var _ orm.Modeler = (*User)(nil)

// UserRepo user repository
type UserRepo interface {
	BasicRepo[User]
}
