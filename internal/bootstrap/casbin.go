package bootstrap

import (
	"context"
	"gin-app/internal/domain"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	gormgenerics "github.com/xjellyx/gorm-generics"
)

func NewEnforcer(db gormgenerics.Database, modelFile string) (casbin.IEnforcer, error) {
	ap, err := gormadapter.NewAdapterByDBWithCustomTable(db.DB(context.Background()), &domain.CasbinRule{})
	if err != nil {
		return nil, err
	}
	en, err := casbin.NewEnforcer(modelFile, ap)
	if err != nil {
		return nil, err
	}
	return en, nil
}
