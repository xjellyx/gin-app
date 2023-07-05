package repository

import (
	"gin-app/internal/domain"
)

// userRepo 学生存储库实现
type userRepo struct {
	domain.BasicRepo[domain.User]
	db domain.Database
}

var _ domain.UserRepo = (*userRepo)(nil)

// NewUserRepo 新建学生存储库
func NewUserRepo(data domain.Database) domain.UserRepo {
	d := NewBasicRepo[domain.User](data)
	stu := &userRepo{
		d,
		data,
	}
	return stu
}
