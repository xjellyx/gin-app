import { request } from '../request';

/** get constant routes */
export function fetchGetConstantRoutes() {
  return request<Api.Route.UserRoute>({
    url: 'menus/constant/tree',
    method: 'get',
  });
}

/** get user routes */
export function fetchGetUserRoutes(code: string) {
  return request<Api.Route.UserRoute>({
    url: '/user/menus',
    method: 'get',
    params: {
      code:code
    }
  });
}

/**
 * whether the route is exist
 *
 * @param routeName route name
 */
export function fetchIsRouteExist(routeName: string) {
  return request<boolean>({ url: '/menus/route/exist', params: { routeName } });
}
