import { adapter } from '@/utils';
import { request } from '../request';
import { adapterOfFetchUserList } from './management.adapter';

/** 获取用户列表 */
export const fetchUserList = async () => {
  const data = await request.get('/api/v1/users');
  return data;
};

export const addUser = async (data: any) => {
  const res = await request.post('/api/v1/users', data);
  return res;
};
