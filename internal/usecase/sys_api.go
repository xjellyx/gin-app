package usecase

import (
	"context"
	"sort"
	"time"

	"gin-app/internal/domain"
	"gin-app/internal/domain/response"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type SysAPIConfig struct {
	Repo    domain.SysAPIRepo
	Timeout time.Duration
}

type sysApiUsecase struct {
	cfg SysAPIConfig
}

func NewSysAPIUsecase(cfg SysAPIConfig) domain.SysAPIUsecase {
	return &sysApiUsecase{
		cfg: cfg,
	}
}

func (u *sysApiUsecase) GetSysAPITree(ctx context.Context) (response.SysAPITreeRespList, error) {
	sysAPIs, err := u.cfg.Repo.Find(ctx, clause.Eq{Column: "white", Value: false})
	if err != nil {
		return nil, err
	}
	var apiMap = make(map[string][]*response.SysAPITreeResp)
	for _, v := range sysAPIs {
		apiMap[v.Tag] = append(apiMap[v.Tag], &response.SysAPITreeResp{
			ID:       int(v.ID),
			Label:    v.Summary,
			Children: nil,
		})
	}
	var ret response.SysAPITreeRespList
	for k, v := range apiMap {
		d := &response.SysAPITreeResp{
			ID:       uuid.New().String(),
			Label:    k,
			Children: v,
		}
		sort.Sort(d.Children)
		ret = append(ret, d)
	}
	sort.Sort(ret)
	return ret, nil
}
