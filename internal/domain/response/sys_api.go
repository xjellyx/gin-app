package response

import "gin-app/pkg/str"

type SysAPITreeResp struct {
	ID       any                `json:"id"`
	Label    string             `json:"label"`
	PId      int                `json:"pId"`
	Children SysAPITreeRespList `json:"children"`
}

type SysAPITreeRespList []*SysAPITreeResp

func (c SysAPITreeRespList) Len() int {
	return len(c)
}

func (c SysAPITreeRespList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c SysAPITreeRespList) Less(i, j int) bool {
	a, _ := str.UTF82GBK(c[i].Label)
	b, _ := str.UTF82GBK(c[j].Label)
	bLen := len(b)
	for i, chr := range a {
		if i > bLen-1 {
			return false
		}
		if chr != b[i] {
			return chr < b[i]
		}
	}
	return true
}
