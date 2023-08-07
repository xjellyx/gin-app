import {request} from '../request';

/**
 * 注册
 */
export function signup(data: any) {
	return request.post('/api/v1/signup', data);
}
