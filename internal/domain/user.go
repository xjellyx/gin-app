package domain

import (
	"context"
	"strings"

	"gin-app/internal/domain/types"
	"gin-app/pkg/serror"

	"github.com/jackc/pgx/v5/pgconn"
	gormgenerics "github.com/xjellyx/gorm-generics"
	"gorm.io/gorm"
)

// User 用户表模型
type User struct {
	gorm.Model
	Uuid     string           `gorm:"size:36;uniqueIndex;default:null;comment:用户uuid"`
	Username string           `gorm:"size:255;uniqueIndex;default:null;comment:name"`
	Email    string           `gorm:"size:255;uniqueIndex;default:null;comment:email"`
	Password string           `gorm:"size:255;comment:password"`
	Gender   types.UserGender `gorm:"default:1;comment:gender,1:male,2:female,3:unknown"`
	Status   types.UserStatus `gorm:"default:1;comment:status,1:normal,2:禁用"`
	Phone    string           `gorm:"size:255;default:null;uniqueIndex;comment:phone"`
	Roles    []*Role          `gorm:"many2many:user_roles;"`
}

// UserRepo user repository
type UserRepo interface {
	gormgenerics.BasicRepo[User]
	GetAllRoles(ctx context.Context, id uint) ([]*Role, error)
}

func TranslateUserDBErr(errV *pgconn.PgError, lan string) error {
	switch {
	case strings.Contains(errV.Detail, "email"):
		return serror.Error(serror.ErrEmailAlredayInUse, lan)
	case strings.Contains(errV.Detail, "name"):
		return serror.Error(serror.ErrUsernameAlredayInUse, lan)
	case strings.Contains(errV.Detail, "phone"):
		return serror.Error(serror.ErrPhoneNumberAlredayInUse, lan)

	}
	return errV
}
