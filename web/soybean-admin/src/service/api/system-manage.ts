import { request } from '../request';

/** get role list */
export function fetchGetRoleList(params?: Api.SystemManage.RoleSearchParams) {
  return request<Api.SystemManage.RoleList>({
    url: '/roles',
    method: 'get',
    params
  });
}

/** add role */
export function fetchAddRole(data: any) {
  return request<Api.SystemManage.Role>({
    url: '/roles',
    method: 'post',
    data
  });
}

/** edit role */
export function fetchEditRole(id: any,data: any) {
  return request<Api.SystemManage.Role>({
    url: `/roles/${id}`,
    method: 'put',
    data
  });
}

/** delete role */
export function fetchDeleteRole(id: any) {
  return request<Api.SystemManage.Role>({
    url: `/roles/${id}`,
    method: 'delete'
  });
}

/** batch delete role */
export function fetchBatchDeleteRole(ids: any) {
  return request<Api.SystemManage.Role>({
    url: '/roles',
    method: 'delete',
    data: ids
  });
}

/**
 * get all roles
 *
 * these roles are all enabled
 */
export function fetchGetAllRoles() {
  return request<Api.SystemManage.AllRole[]>({
    url: '/roles/all',
    method: 'get'
  });
}

export function fetchGetUserRoles(id: Number) {
  return request({
    url: `/user/roles`,
    method: 'get'
  })

}

export function fetchAddRoleMenuPerm(data: any){
  return request({
    url: "roles/menu",
    method: "post",
    data
  })
}
export function fetchGetRoleMenuPerm(id: Number) {
  return request({
    url:`/roles/${id}/menu`,
    method: "get"
  })
}

export function fetchGetRolePerm(id: Number) {
  return request({
    url: `/roles/${id}/perm`,
    method: 'get'
  })
}

export function fetchAddRolePerm(data: any){
  return request({
    url: "roles/perm",
    method: "post",
    data
  })
}

export function fetchGetRoleFrontPage(id: Number) {
  return request({
    url: `/roles/${id}/front-page`,
    method: 'get'
  })
}

export function fetchAddRoleFrontPage(data: any){
  return request({
    url: "roles/front-page",
    method: "post",
    data
  })
}

/** get user list */
export function fetchGetUserList(params?: Api.SystemManage.UserSearchParams) {
  return request<Api.SystemManage.UserList>({
    url: '/users',
    method: 'get',
    params
  });
}
/** add user */
export function fetchAddUser(data: any) {
  return request<Api.SystemManage.User>({
    url: '/users',
    method: 'post',
    data
  });
}
/** edit user */
export function fetchEditUser(id: any,data: any) {
  return request<Api.SystemManage.User>({
    url: `/users/${id}`,
    method: 'put',
    data
  });
}

/** delete user */
export function fetchDeleteUser(id: any) {
  return request<Api.SystemManage.User>({
    url: `/users/${id}`,
    method: 'delete'
  });
}
/** batch delete user */
export function fetchBatchDeleteUser(ids: any) {
  return request<Api.SystemManage.User>({
    url: `/users`,
    method: 'delete',
    data:ids
  });
}

/** get menu list */
export function fetchGetMenuList() {
  return request<Api.SystemManage.MenuList>({
    url: '/menus',
    method: 'get'
  });
}

/** get all pages */
export function fetchGetAllPages() {
  return request<string[]>({
    url: '/menus/pages',
    method: 'get'
  });
}

/** get menu tree */
export function fetchGetMenuTree() {
  return request<Api.SystemManage.MenuTree[]>({
    url: '/menus/tree',
    method: 'get'
  });
}

/** get sys api tree */
export function fetchSysApiTree() {
  return request<Api.SystemManage.MenuTree[]>({
    url: '/sys-api/tree',
    method: 'get'
  });
}

export function fetchAddMenu(data: any) {
  return request<Api.SystemManage.Menu>({
    url: '/menus',
    method: 'post',
    data
  });
}

export function fetchEditMenu(id: any,data: any) {
  return request<Api.SystemManage.Menu>({
    url: `/menus/${id}`,
    method: 'put',
    data
  });
}

export function fetchDeleteMenu(id: any) {
  return request<Api.SystemManage.Menu>({
    url: `/menus/${id}`,
    method: 'delete'
  });
}

export function fetchBatchDeleteMenu(ids: any) {
  return request<Api.SystemManage.Menu>({
    url: `/menus`,
    method: 'delete',
    data:ids
  });
}

export function fetchGetMenuPerm(id: any) {
  return request({
    url: `/menus/${id}/perm`,
    method: 'get'
  });
}

export function fetchAddMenuPerm(data: any) {
  return request({
    url: `/menus/perm`,
    method: 'post',
    data
  });
}


