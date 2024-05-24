<script setup lang="ts">
import { computed, shallowRef, watch } from 'vue';
import { $t } from '@/locales';
import { fetchAddMenuPerm, fetchGetMenuPerm, fetchSysApiTree } from '@/service/api';

defineOptions({
  name: 'SysApiAuthModal'
});

interface Props {
  /** the menuId */
  menuId: number;
}

const props = defineProps<Props>();

const visible = defineModel<boolean>('visible', {
  default: false
});

function closeModal() {
  visible.value = false;
}

const title = computed(() => $t('common.edit') + $t('page.manage.role.apiAuth'));

const tree = shallowRef<Api.SystemManage.MenuTree[]>([]);

async function getTree() {
  const { error, data } = await fetchSysApiTree();
  if (!error) {
    tree.value = data;
  }
}

const checks = shallowRef<number[]>([]);

async function getChecks() {
  // request
  const { data } = await fetchGetMenuPerm(props.menuId);
  if (!data) {
    return;
  }
  checks.value = data;
}

async function handleSubmit() {
  console.log(checks.value, props.menuId);
  // request
  await fetchAddMenuPerm({
    menuId: props.menuId,
    apiIds: checks.value
  });
  window.$message?.success?.($t('common.modifySuccess'));

  closeModal();
}

function init() {
  getTree();
  getChecks();
}

watch(visible, val => {
  if (val) {
    init();
  }
});
</script>

<template>
  <NModal v-model:show="visible" :title="title" preset="card" class="w-480px">
    <NTree
      v-model:checked-keys="checks"
      :data="tree"
      key-field="id"
      cascade
      checkable
      expand-on-click
      virtual-scroll
      block-line
      class="h-280px"
    />
    <template #footer>
      <NSpace justify="end">
        <NButton size="small" class="mt-16px" @click="closeModal">
          {{ $t('common.cancel') }}
        </NButton>
        <NButton type="primary" size="small" class="mt-16px" @click="handleSubmit">
          {{ $t('common.confirm') }}
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<style scoped></style>
