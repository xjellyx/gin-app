package request

type AddAPIPermissionReq struct {
	ID     int   `json:"id"`
	APIIds []any `json:"apiIds"`
}

type AddRoleMenuPermissionReq struct {
	ID      int   `json:"id"`
	MenuIds []int `json:"menuIds"`
}

type SetRoleRouteFrontPageReq struct {
	ID        int    `json:"id"`
	RoutePath string `json:"routePath"`
}
