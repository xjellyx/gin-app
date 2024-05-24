<script setup lang="ts">
import { computed, ref } from 'vue';
import type { VNode } from 'vue';
import { useAuthStore } from '@/store/modules/auth';
import { useRouteStore } from '@/store/modules/route';
import { useRouterPush } from '@/hooks/common/router';
import { useSvgIcon } from '@/hooks/common/icon';
import { $t } from '@/locales';
import { fetchGetUserRoles } from '@/service/api';
import { localStg } from '@/utils/storage';

defineOptions({
  name: 'UserAvatar'
});

const authStore = useAuthStore();
const routeStore = useRouteStore();
const { routerPushByKey, toLogin } = useRouterPush();
const { SvgIconVNode } = useSvgIcon();
function loginOrRegister() {
  toLogin();
}
const showDrawer = ref(false);
type DropdownKey = 'user-center' | 'logout';
const roleCode = ref(localStg.get('currentRole'));
type DropdownOption =
  | {
      key: DropdownKey;
      label: string;
      icon?: () => VNode;
    }
  | {
      type: 'divider';
      key: string;
    };

const roleOptions = ref<CommonType.Option<string>[]>([]);

async function getRoleOptions() {
  const { error, data } = await fetchGetUserRoles();

  if (!error) {
    const options = data.map(item => ({
      label: item.name,
      value: item.code
    }));
    roleOptions.value = [...options];
  }
}
const options = computed(() => {
  const opts: DropdownOption[] = [
    {
      label: '切换角色',
      icon: SvgIconVNode({ icon: 'ph:user-switch', fontSize: 18 }),
      key: 'switch-role'
    },
    {
      label: $t('common.userCenter'),
      key: 'user-center',
      icon: SvgIconVNode({ icon: 'ph:user-circle', fontSize: 18 })
    },
    {
      type: 'divider',
      key: 'divider'
    },
    {
      label: $t('common.logout'),
      key: 'logout',
      icon: SvgIconVNode({ icon: 'ph:sign-out', fontSize: 18 })
    }
  ];

  return opts;
});

function logout() {
  window.$dialog?.info({
    title: $t('common.tip'),
    content: $t('common.logoutConfirm'),
    positiveText: $t('common.confirm'),
    negativeText: $t('common.cancel'),
    onPositiveClick: () => {
      authStore.resetStore();
    }
  });
}

function handleDropdown(key: DropdownKey) {
  console.log('key', key);
  if (key === 'logout') {
    logout();
  } else if (key === 'switch-role') {
    showDrawer.value = true;
    getRoleOptions();
  } else {
    routerPushByKey(key);
  }
}
// 切换角色，刷新整个应用
async function handleRoleChange() {
  const currentRole = localStg.get('currentRole');
  if (currentRole === roleCode.value) {
    showDrawer.value = false;
    return;
  }
  localStg.set('currentRole', roleCode.value);
  await routeStore.initAuthRoute();
  showDrawer.value = false;
  window.location.reload();
}
</script>

<template>
  <NButton v-if="!authStore.isLogin" quaternary @click="loginOrRegister">
    {{ $t('page.login.common.loginOrRegister') }}
  </NButton>
  <NDropdown v-else placement="bottom" trigger="click" :options="options" @select="handleDropdown">
    <div>
      <ButtonIcon>
        <SvgIcon icon="ph:user-circle" class="text-icon-large" />
        <span class="text-16px font-medium">{{ authStore.userInfo.username }}</span>
      </ButtonIcon>
    </div>
  </NDropdown>
  <NDrawer v-model:show="showDrawer" :width="150">
    <NSelect
      v-model:value="roleCode"
      :options="roleOptions"
      :placeholder="$t('page.manage.user.form.userRole')"
      @update:value="handleRoleChange"
    />
  </NDrawer>
</template>

<style scoped></style>
