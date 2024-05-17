package request

import (
	"gorm.io/gorm/clause"
)

type Page struct {
	PageSize int `json:"pageSize" form:"pageSize,default=10"` // 每页数量
	PageNum  int `json:"pageNum" form:"pageNum,default=1"`    // 页数
}

func (q Page) SetOrmExpression() clause.Limit {
	return clause.Limit{
		Limit:  &q.PageSize,
		Offset: (q.PageNum - 1) * q.PageSize,
	}
}
