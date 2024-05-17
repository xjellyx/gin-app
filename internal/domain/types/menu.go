package types

type MenuStatus int

const (
	MenuStatusEnable  MenuStatus = iota + 1 // 启用
	MenuStatusDisable                       // 禁用
)

type MenuType int

const (
	MenuTypeDir    MenuType = iota + 1 // 目录
	MenuTypeMenu                       // 菜单
	MenuTypeButton                     // 按钮
)

type MenuIconType int

const (
	MenuIconTypeIconify MenuIconType = iota + 1 // iconify
	MenuIconTypeLocal                           // local
)

type MenuQuery struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type MenuButton struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}
