<script setup lang="ts">
import type { PlatformDomain } from '../domain-client';

import { ref, watch } from 'vue';

import { Popconfirm, Space, Switch, message } from 'ant-design-vue';

import { $t } from '#/locales';

import {
  platformDomainCreate,
  platformDomainDelete,
  platformDomainList,
  platformDomainSetVerified,
} from '../domain-client';

const props = defineProps<{ tenantId: number }>();

const domains = ref<PlatformDomain[]>([]);
const newDomain = ref('');

async function refresh() {
  if (!props.tenantId) {
    domains.value = [];
    return;
  }
  const { items } = await platformDomainList({ tenantId: props.tenantId });
  domains.value = items;
}

async function handleAdd() {
  const domain = newDomain.value.trim();
  if (!domain) {
    return;
  }
  await platformDomainCreate({ tenantId: props.tenantId, domain });
  message.success($t('pages.common.createSuccess'));
  newDomain.value = '';
  await refresh();
}

async function handleSetVerified(checked: boolean, row: PlatformDomain) {
  await platformDomainSetVerified(row.id, checked);
  message.success($t('pages.multiTenant.domain.messages.verifiedUpdated'));
  await refresh();
}

async function handleDelete(row: PlatformDomain) {
  await platformDomainDelete(row.id);
  message.success($t('pages.common.deleteSuccess'));
  await refresh();
}

watch(() => props.tenantId, refresh, { immediate: true });
</script>

<template>
  <div class="flex flex-col gap-2" data-testid="tenant-domains-section">
    <div
      v-for="row in domains"
      :key="row.id"
      class="flex items-center justify-between gap-2"
      data-testid="tenant-domain-row"
    >
      <span class="truncate">{{ row.domain }}</span>
      <Space>
        <Switch
          :checked="row.isVerified"
          :checked-children="$t('pages.multiTenant.domain.verified.yes')"
          :data-testid="`tenant-domain-verify-${row.id}`"
          :un-checked-children="$t('pages.multiTenant.domain.verified.no')"
          @change="(checked) => handleSetVerified(Boolean(checked), row)"
        />
        <Popconfirm
          placement="left"
          :title="
            $t('pages.multiTenant.domain.messages.deleteDomainConfirm', {
              domain: row.domain,
            })
          "
          @confirm="handleDelete(row)"
        >
          <a-button
            danger
            size="small"
            :data-testid="`tenant-domain-delete-${row.id}`"
          >
            {{ $t('pages.common.delete') }}
          </a-button>
        </Popconfirm>
      </Space>
    </div>
    <Space>
      <a-input
        v-model:value="newDomain"
        allow-clear
        data-testid="tenant-domain-input"
        :placeholder="$t('pages.multiTenant.domain.placeholder')"
      />
      <a-button
        data-testid="tenant-domain-add"
        type="primary"
        @click="handleAdd"
      >
        {{ $t('pages.common.add') }}
      </a-button>
    </Space>
  </div>
</template>
