<template>
  <n-select v-model:value="singUpValue" :options="singUpOptions" style="bottom: 18px"></n-select>
  <n-form ref="formRef" :model="model" :rules="rules" size="large" :show-label="false">
    <n-form-item v-if="singUpValue === 'phone'" path="phone">
      <n-input v-model:value="model.phone" :placeholder="$t('page.login.common.phonePlaceholder')" />
    </n-form-item>
    <n-form-item v-else-if="singUpValue === 'email'" path="email">
      <n-input v-model:value="model.email" :placeholder="$t('page.login.common.emailPlaceholder')" />
    </n-form-item>
    <n-form-item path="code">
      <div class="flex-y-center w-full">
        <n-input v-model:value="model.code" :placeholder="$t('page.login.common.codePlaceholder')" />
        <div class="w-18px"></div>
        <n-button size="large" :disabled="isCounting" :loading="smsLoading" @click="handleSmsCode">
          {{ label }}
        </n-button>
      </div>
    </n-form-item>
    <n-form-item path="pwd">
      <n-input
        v-model:value="model.password"
        type="password"
        show-password-on="click"
        :placeholder="$t('page.login.common.passwordPlaceholder')"
      />
    </n-form-item>
    <n-form-item path="confirmPwd">
      <n-input
        v-model:value="model.confirmPwd"
        type="password"
        show-password-on="click"
        :placeholder="$t('page.login.common.confirmPasswordPlaceholder')"
      />
    </n-form-item>
    <n-space :vertical="true" :size="18">
      <login-agreement v-model:value="agreement" />
      <n-button type="primary" size="large" :block="true" :round="true" @click="handleSubmit">
        {{ $t('page.login.common.confirm') }}
      </n-button>
      <n-button size="large" :block="true" :round="true" @click="toLoginModule('pwd-login')">
        {{ $t('page.login.common.back') }}
      </n-button>
    </n-space>
  </n-form>
</template>

<script lang="ts" setup>
import { reactive, ref, toRefs } from 'vue';
import type { FormInst, FormRules } from 'naive-ui';
import { signup } from '@/service';
import { useRouterPush } from '@/composables';
import { useSmsCode } from '@/hooks';
import { formRules, getConfirmPwdRule } from '@/utils';
import { $t } from '@/locales';

const { toLoginModule } = useRouterPush();
const { label, isCounting, loading: smsLoading, start } = useSmsCode();

const formRef = ref<HTMLElement & FormInst>();

const singUpValue = ref<string>('phone');
const singUpOptions = reactive([
  { label: $t('page.register.phone'), value: 'phone' },
  { label: $t('page.register.email'), value: 'email' }
]);

const model = reactive({
  phone: '',
  email: '',
  code: '123456',
  password: '',
  confirmPwd: ''
});

const rules: FormRules = {
  phone: formRules.phone,
  code: formRules.code,
  password: formRules.password,
  confirmPwd: getConfirmPwdRule(toRefs(model).password)
};

const agreement = ref(false);

function handleSmsCode() {
  start();
}

async function handleSubmit() {
  await formRef.value?.validate();
  await signup(model);
  window.$message?.success($t('page.login.common.validateSuccess'));
  toLoginModule('pwd-login');
}
</script>

<style scoped></style>
