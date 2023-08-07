import { $t } from '@/locales';
import { transformObjectToOption } from './_shared';

export const loginModuleLabels: Record<UnionKey.LoginModule, string> = {
  'pwd-login': $t('page.login.pwdLogin.title'),
  'code-login': $t('page.login.codeLogin.title'),
  register: $t('page.login.register.title'),
  'reset-pwd': $t('page.login.resetPwd.title'),
  'bind-wechat': $t('page.login.bindWeChat.title')
};

export const userRoleLabels: Record<Auth.RoleType, string> = {
  super: $t('page.login.pwdLogin.superAdmin'),
  admin: $t('page.login.pwdLogin.admin'),
  user: $t('page.login.pwdLogin.user')
};
export const userRoleOptions = transformObjectToOption(userRoleLabels);

/** 用户性别 */
export const genderLabels: Record<UserManagement.GenderKey, string> = {
  1: $t('page.user.gender.unknown'),
  2: $t('page.user.gender.male'),
  3: $t('page.user.gender.female')
};
export const genderOptions = transformObjectToOption(genderLabels);

/** 用户状态 */
export const userStatusLabels: Record<UserManagement.UserStatusKey, string> = {
  1: $t('page.user.status.normal'),
  2: $t('page.user.status.locked'),
  3: $t('page.user.status.freeze'),
  4: $t('page.user.status.deleted')
};
export const userStatusOptions = [
  {
    value: 1,
    label: $t('page.user.status.normal')
  },
  {
    value: 2,
    label: $t('page.user.status.locked')
  },
  {
    value: 3,
    label: $t('page.user.status.freeze')
  },
  {
    value: 4,
    label: $t('page.user.status.deleted')
  }
];
