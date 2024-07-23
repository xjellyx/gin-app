package response

// Response 返回体
type Response struct {
	// Code 0为成功，其他都是失败
	Code string `json:"code"`
	// Msg 返回提示
	Msg string `json:"msg"`
	// InternalError 错误信息
	InternalError string `json:"internalError,omitempty"`
	// Data 返回数据
	Data any `json:"data"`
}

// Pagination 页数信息
type Pagination struct {
	Total    int64 `json:"total"`
	PageSize int   `json:"pageSize"`
	PageNum  int   `json:"pageNum"`
}
