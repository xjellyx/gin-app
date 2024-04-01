package domain

import (
	"strings"

	"gin-app/pkg/serror"

	"github.com/jackc/pgx/v5/pgconn"
	gormgenerics "github.com/olongfen/gorm-generics"
	"gorm.io/gorm"
)

type UserGender uint // user gender
const (
	GenderUnknown UserGender = iota + 1
	GenderMale               // male
	GenderFemale             // female
)

type UserStatus uint // user status
const (
	StatusNormal  UserStatus = iota + 1 // 正常
	StatusLocked                        // 锁定
	StatusFreeze                        // 冻结
	StatusDeleted                       // deleted

)

// User 学生表模型
type User struct {
	gorm.Model
	Uuid     string     `gorm:"size:36;uniqueIndex;default:null;comment:用户uuid"`
	Username string     `gorm:"size:255;uniqueIndex;default:null;comment:name"`
	Email    string     `gorm:"size:255;uniqueIndex;default:null;comment:email"`
	Password string     `gorm:"size:255;comment:password"`
	Gender   UserGender `gorm:"default:1;comment:gender,1:unknown,2:male,3:female"`
	Status   UserStatus `gorm:"default:1;comment:status,1:normal,2:锁定,3:冻结,4:deleted"`
	Phone    string     `gorm:"size:255;default:null;uniqueIndex;comment:phone"`
}

// UserRepo user repository
type UserRepo interface {
	gormgenerics.BasicRepo[User]
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
