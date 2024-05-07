package request

type QueryReq struct {
	PageSize uint `query:"pageSize"` // 每页数量
	PageNum  uint `query:"pageNum"`  // 页数
}
