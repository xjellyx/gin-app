<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import { useFormRules, useNaiveForm } from '@/hooks/common/form';
import { fetchAddUser, fetchEditUser, fetchGetAllRoles } from '@/service/api';
import { $t } from '@/locales';
import { enableStatusOptions, userGenderOptions } from '@/constants/business';

defineOptions({
  name: 'UserOperateDrawer'
});

interface Props {
  /** the type of operation */
  operateType: NaiveUI.TableOperateType;
  /** the edit row data */
  rowData?: Api.SystemManage.User | null;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'submitted'): void;
}

const emit = defineEmits<Emits>();

const visible = defineModel<boolean>('visible', {
  default: false
});

const { formRef, validate, restoreValidation } = useNaiveForm();
const { defaultRequiredRule } = useFormRules();

const title = computed(() => {
  const titles: Record<NaiveUI.TableOperateType, string> = {
    add: $t('page.manage.user.addUser'),
    edit: $t('page.manage.user.editUser')
  };
  return titles[props.operateType];
});

type Model = Pick<Api.SystemManage.User, 'username' | 'gender' | 'nickName' | 'phone' | 'email' | 'roles' | 'status'>;

const model: Model = reactive(createDefaultModel());

function createDefaultModel(): Model {
  return {
    username: '',
    gender: null,
    nickName: '',
    phone: '',
    email: '',
    roles: [],
    status: null
  };
}

type RuleKey = Extract<keyof Model, 'username' | 'status'>;

const rules: Record<RuleKey, App.Global.FormRule> = {
  username: defaultRequiredRule,
  status: defaultRequiredRule
};

/** the enabled role options */
const roleOptions = ref<CommonType.Option<string>[]>([]);

async function getRoleOptions() {
  const { error, data } = await fetchGetAllRoles();

  if (!error) {
    const options = data.map(item => ({
      label: item.name,
      value: item.code
    }));
    roleOptions.value = [...options];
  }
}

function handleInitModel() {
  Object.assign(model, createDefaultModel());

  if (props.operateType === 'edit' && props.rowData) {
    Object.assign(model, props.rowData);
  }
}

function closeDrawer() {
  visible.value = false;
}

async function handleSubmit() {
  await validate();
  // request
  if (props.operateType === 'add') {
    // add
    await fetchAddUser(model);
  } else if (props.operateType === 'edit') {
    // edit
    await fetchEditUser(props.rowData?.id, model);
  }
  window.$message?.success($t('common.updateSuccess'));
  closeDrawer();
  emit('submitted');
}

watch(visible, () => {
  if (visible.value) {
    handleInitModel();
    restoreValidation();
    getRoleOptions();
  }
});
</script>

<template>
  <NDrawer v-model:show="visible" display-directive="show" :width="360">
    <NDrawerContent :title="title" :native-scrollbar="false" closable>
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem :label="$t('page.manage.user.userName')" path="userName">
          <NInput v-model:value="model.username" :placeholder="$t('page.manage.user.form.userName')" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.userGender')" path="userGender">
          <NRadioGroup v-model:value="model.gender">
            <NRadio
              v-for="item in userGenderOptions"
              :key="item.value"
              :value="Number(item.value)"
              :label="$t(item.label)"
            />
          </NRadioGroup>
        </NFormItem>
        <!--        <NFormItem :label="$t('page.manage.user.nickName')" path="nickName">-->
        <!--          <NInput v-model:value="model.nickName" :placeholder="$t('page.manage.user.form.nickName')" />-->
        <!--        </NFormItem>-->
        <NFormItem :label="$t('page.manage.user.userPhone')" path="userPhone">
          <NInput v-model:value="model.phone" :placeholder="$t('page.manage.user.form.userPhone')" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.userEmail')" path="email">
          <NInput v-model:value="model.email" :placeholder="$t('page.manage.user.form.userEmail')" />
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.userStatus')" path="status">
          <NRadioGroup v-model:value="model.status">
            <NRadio
              v-for="item in enableStatusOptions"
              :key="item.value"
              :value="Number(item.value)"
              :label="$t(item.label)"
            />
          </NRadioGroup>
        </NFormItem>
        <NFormItem :label="$t('page.manage.user.userRole')" path="roles">
          <NSelect
            v-model:value="model.roles"
            multiple
            :options="roleOptions"
            :placeholder="$t('page.manage.user.form.userRole')"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace :size="16">
          <NButton @click="closeDrawer">{{ $t('common.cancel') }}</NButton>
          <NButton type="primary" @click="handleSubmit">{{ $t('common.confirm') }}</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped></style>
