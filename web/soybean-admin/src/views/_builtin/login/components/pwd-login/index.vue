<template>
  <n-select v-model:value="singInType" :options="singInOptions" style="bottom: 18px"></n-select>
  <n-form ref="formRef" :model="model" :rules="rules" size="large" :show-label="false">
    <n-form-item v-if="singInType === 'username'" path="username">
      <n-input v-model:value="model.username" :placeholder="$t('page.login.common.userNamePlaceholder')" />
    </n-form-item>
    <n-form-item v-else-if="singInType === 'phone'" path="phone">
      <n-input v-model:value="model.phone" :placeholder="$t('page.login.common.phonePlaceholder')" />
    </n-form-item>
    <n-form-item v-else="singInType === 'email'" path="email">
      <n-input v-model:value="model.email" :placeholder="$t('page.login.common.email')" />
    </n-form-item>
    <n-form-item path="password">
      <n-input
        v-model:value="model.password"
        type="password"
        show-password-on="click"
        :placeholder="$t('page.login.common.passwordPlaceholder')"
      />
    </n-form-item>
    <n-space :vertical="true" :size="24">
      <div class="flex-y-center justify-between">
        <n-checkbox v-model:checked="rememberMe">{{ $t('page.login.pwdLogin.rememberMe') }}</n-checkbox>
        <n-button :text="true" @click="toLoginModule('reset-pwd')">
          {{ $t('page.login.pwdLogin.forgetPassword') }}
        </n-button>
      </div>
      <n-button
        type="primary"
        size="large"
        :block="true"
        :round="true"
        :loading="auth.loginLoading"
        @click="handleSubmit"
      >
        {{ $t('page.login.common.confirm') }}
      </n-button>
      <div class="flex-y-center justify-between">
        <n-button class="flex-1" :block="true" @click="toLoginModule('code-login')">
          {{ loginModuleLabels['code-login'] }}
        </n-button>
        <div class="w-12px"></div>
        <n-button class="flex-1" :block="true" @click="toLoginModule('register')">
          {{ loginModuleLabels.register }}
        </n-button>
      </div>
    </n-space>
<!--    <other-account @login="handleLoginOtherAccount" />-->
  </n-form>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import type { FormInst, FormRules } from 'naive-ui';
import { loginModuleLabels } from '@/constants';
import { useAuthStore } from '@/store';
import { useRouterPush } from '@/composables';
import { formRules } from '@/utils';
import { $t } from '@/locales';
// import { OtherAccount } from './components';
import { SingIn } from '@/typings';

const auth = useAuthStore();
const { login } = useAuthStore();
const { toLoginModule } = useRouterPush();

const formRef = ref<HTMLElement & FormInst>();
const singInType = ref('phone');
const singInOptions = reactive([
  { label: $t('page.login.common.phone'), value: 'phone' },
  { label: $t('page.login.common.email'), value: 'email' },
  { label: $t('page.login.common.username'), value: 'username' }
]);
const model = reactive({
  username: '',
  email: '',
  phone: '13768886999',
  password: 'Qaz12345'
});

const rules: FormRules = {
  password: formRules.password
};

const rememberMe = ref(false);

async function handleSubmit() {
  await formRef.value?.validate();

  await login(model);
}

// function handleLoginOtherAccount(param: SingIn) {
//   const { username, password } = param;
//   login(username, password);
// }
</script>

<style scoped></style>
