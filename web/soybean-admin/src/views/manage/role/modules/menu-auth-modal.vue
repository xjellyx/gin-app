<script setup lang="ts">
import {computed, shallowRef, watch} from 'vue';
import {$t} from '@/locales';
import {
  fetchAddRoleFrontPage,
  fetchAddRoleMenuPerm,
  fetchGetAllPages,
  fetchGetMenuTree,
  fetchGetRoleFrontPage,
  fetchGetRoleMenuPerm
} from '@/service/api';

defineOptions({
  name: 'MenuAuthModal'
});

interface Props {
  /** the roleId */
  roleId: number;
}

const props = defineProps<Props>();

const visible = defineModel<boolean>('visible', {
  default: false
});

function closeModal() {
  visible.value = false;
}

const title = computed(() => $t('common.edit') + $t('page.manage.role.menuAuth'));

const home = shallowRef('');

async function getHome() {
  console.log(props.roleId);
  const {data}= await fetchGetRoleFrontPage(props.roleId);
  if(data){
    home.value = data
  }
}
function  setTreeLabels (tree: Api.SystemManage.MenuTree[]) {
  tree.forEach(item => {
    item.label  = item.i18nKey ? $t(item.i18nKey) : item.label ;
    if (item.children) {
      setTreeLabels(item.children);
    }
  });
}


async function updateHome(val: string) {
  // request
  home.value = val;
  await fetchAddRoleFrontPage({id: props.roleId,routePath: val})
}

const pages = shallowRef<string[]>([]);

async function getPages() {
  const { error, data } = await fetchGetAllPages();

  if (!error) {
    pages.value = data;
  }
}

const pageSelectOptions = computed(() => {
  const opts: CommonType.Option[] = pages.value.map(page => ({
    label: page,
    value: page
  }));

  return opts;
});

const tree = shallowRef<Api.SystemManage.MenuTree[]>([]);


async function getTree() {
  const { error, data } = await fetchGetMenuTree();
  setTreeLabels(data);
  if (!error) {
    tree.value = data;
  }
}

const checks = shallowRef<number[]>([]);

async function getChecks() {
   const {data} =  await fetchGetRoleMenuPerm(props.roleId)
  // request
  checks.value = data;
}

const handleCheck = (keys) => {
  console.log("32432424323")
  const selectedKeys = new Set(keys)
  const parentKeys = new Set()

  const findParent = (nodes, parentKey = null) => {
    for (const node of nodes) {
      if (selectedKeys.has(node.id) && parentKey) {
        parentKeys.add(parentKey)
      }
      if (node.children) {
        findParent(node.children, node.id)
      }
    }
  }
  findParent(tree.value)
  checks.value = Array.from(new Set([...selectedKeys, ...parentKeys]))
}
async function handleSubmit() {
  // request
  await fetchAddRoleMenuPerm({id: props.roleId,menuIds: checks.value})
  window.$message?.success?.($t('common.modifySuccess'));

  closeModal();
}

function init() {
  getHome();
  getPages();
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
    <div class="flex-y-center gap-16px pb-12px">
      <div>{{ $t('page.manage.menu.home') }}</div>
      <NSelect :value="home" :options="pageSelectOptions" size="small" class="w-160px" @update:value="updateHome" />
    </div>
    <NTree
      v-model:checked-keys="checks"
      :data="tree"
      key-field="id"
      checkable
      @update:checked-keys="handleCheck"
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
