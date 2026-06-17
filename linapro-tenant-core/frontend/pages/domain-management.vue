<script lang="ts">
export const pluginPageMeta = {
  capabilities: ['tenant.management'],
  routePath: 'platform/domains',
  title: 'Tenant Domain Management',
};
</script>

<script setup lang="ts">
import type { SwitchProps } from 'ant-design-vue';

import type { PlatformDomain } from './domain-client';

import { Page, useVbenModal } from '@vben/common-ui';

import { Popconfirm, Space, Switch, Tag, message } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { $t } from '#/locales';
import { formatTimestamp } from '#/utils/time';

import DomainModal from './components/domain-modal.vue';
import {
  platformDomainDelete,
  platformDomainList,
  platformDomainSetVerified,
} from './domain-client';

const [DomainModalRef, domainModalApi] = useVbenModal({
  connectedComponent: DomainModal,
});

const statusColorMap: Record<string, string> = {
  active: 'green',
  disabled: 'red',
};

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: {
    schema: [
      {
        component: 'Input',
        componentProps: {
          'data-testid': 'domain-search-domain',
        },
        fieldName: 'domain',
        label: $t('pages.multiTenant.domain.fields.domain'),
      },
      {
        component: 'Input',
        componentProps: {
          'data-testid': 'domain-search-tenant',
        },
        fieldName: 'tenantId',
        label: $t('pages.multiTenant.fields.tenantId'),
      },
      {
        component: 'Select',
        componentProps: {
          'data-testid': 'domain-search-status',
          options: ['active', 'disabled'].map((value) => ({
            label: $t(`pages.multiTenant.domain.status.${value}`),
            value,
          })),
        },
        fieldName: 'status',
        label: $t('pages.common.status'),
      },
    ],
    commonConfig: { componentProps: { allowClear: true }, labelWidth: 90 },
    wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
  },
  gridOptions: {
    columns: [
      {
        field: 'domain',
        minWidth: 220,
        title: $t('pages.multiTenant.domain.fields.domain'),
      },
      {
        field: 'tenantId',
        minWidth: 120,
        title: $t('pages.multiTenant.fields.tenantId'),
      },
      {
        field: 'isVerified',
        slots: { default: 'verified' },
        title: $t('pages.multiTenant.domain.fields.isVerified'),
        width: 160,
      },
      {
        field: 'status',
        slots: { default: 'status' },
        title: $t('pages.common.status'),
        width: 120,
      },
      {
        field: 'createdAt',
        formatter: ({ cellValue }) => formatTimestamp(cellValue),
        title: $t('pages.common.createdAt'),
        width: 180,
      },
      {
        field: 'action',
        fixed: 'right',
        slots: { default: 'action' },
        title: $t('pages.common.actions'),
        width: 200,
      },
    ],
    emptyRender: {
      name: 'Empty',
      props: { description: $t('pages.multiTenant.domain.empty.domains') },
    },
    height: 'auto',
    pagerConfig: {},
    proxyConfig: {
      ajax: {
        query: async (
          { page }: { page: { currentPage: number; pageSize: number } },
          formValues = {},
        ) =>
          await platformDomainList({
            pageNum: page.currentPage,
            pageSize: page.pageSize,
            ...formValues,
          }),
      },
    },
    rowConfig: { keyField: 'id' },
    id: 'platform-tenant-domain-index',
  },
});

function openCreate() {
  domainModalApi.setData({});
  domainModalApi.open();
}

async function handleSetVerified(
  checked: SwitchProps['checked'],
  row: PlatformDomain,
) {
  await platformDomainSetVerified(row.id, Boolean(checked));
  message.success($t('pages.multiTenant.domain.messages.verifiedUpdated'));
  await gridApi.query();
}

async function handleDelete(row: PlatformDomain) {
  await platformDomainDelete(row.id);
  message.success($t('pages.common.deleteSuccess'));
  await gridApi.query();
}
</script>

<template>
  <Page :auto-content-height="true" data-testid="platform-domains-page">
    <Grid :table-title="$t('pages.multiTenant.domain.tableTitle')">
      <template #toolbar-tools>
        <Space>
          <a-button
            type="primary"
            data-testid="domain-create"
            @click="openCreate"
          >
            {{ $t('pages.common.add') }}
          </a-button>
        </Space>
      </template>
      <template #verified="{ row }">
        <Switch
          :checked="row.isVerified"
          :data-testid="`domain-verify-${row.id}`"
          :checked-children="$t('pages.multiTenant.domain.verified.yes')"
          :un-checked-children="$t('pages.multiTenant.domain.verified.no')"
          @change="(checked) => handleSetVerified(checked, row)"
        />
      </template>
      <template #status="{ row }">
        <Tag :color="statusColorMap[row.status] || 'default'">
          {{ $t(`pages.multiTenant.domain.status.${row.status || 'active'}`) }}
        </Tag>
      </template>
      <template #action="{ row }">
        <Space>
          <Popconfirm
            placement="left"
            :title="
              $t('pages.multiTenant.domain.messages.deleteDomainConfirm', {
                domain: row.domain,
              })
            "
            @confirm="handleDelete(row)"
          >
            <ghost-button
              danger
              :data-testid="`domain-delete-${row.id}`"
              @click.stop=""
            >
              {{ $t('pages.common.delete') }}
            </ghost-button>
          </Popconfirm>
        </Space>
      </template>
    </Grid>
    <DomainModalRef @success="gridApi.query()" />
  </Page>
</template>
