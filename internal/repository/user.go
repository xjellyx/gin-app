package repository

import (
	"context"
	"gin-app/internal/domain"

	gormgenerics "github.com/xjellyx/gorm-generics"
)

// userRepo 学生存储库实现
type userRepo struct {
	gormgenerics.BasicRepo[domain.User]
	db gormgenerics.Database
}

var _ domain.UserRepo = (*userRepo)(nil)

// NewUserRepo 新建学生存储库
func NewUserRepo(data gormgenerics.Database) domain.UserRepo {
	d := gormgenerics.NewBasicRepository[domain.User](data)
	stu := &userRepo{
		d,
		data,
	}
	return stu
}

func (u *userRepo) GetAllRoles(ctx context.Context, id uint) ([]*domain.Role, error) {
	var user *domain.User
	err := u.db.DB(ctx).Model(&domain.User{}).Preload("Roles").Where("id = ?", id).First(&user).Error
	return user.Roles, err
}
