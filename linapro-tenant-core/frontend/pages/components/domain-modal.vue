<script setup lang="ts">
import { useVbenModal } from '@vben/common-ui';

import { message } from 'ant-design-vue';

import { useVbenForm, z } from '#/adapter/form';
import { $t } from '#/locales';

import { platformDomainCreate } from '../domain-client';

const emit = defineEmits<{ success: [] }>();

interface DomainFormValues {
  tenantId: number;
  domain: string;
  isPrimary: boolean;
}

const formValues: DomainFormValues = {
  tenantId: 0,
  domain: '',
  isPrimary: false,
};

const [DomainForm, formApi] = useVbenForm({
  commonConfig: {
    componentProps: {
      class: 'w-full',
    },
    labelWidth: 90,
  },
  schema: [
    {
      component: 'InputNumber',
      componentProps: {
        'data-testid': 'domain-tenant-input',
        class: 'w-full',
        min: 1,
        precision: 0,
      },
      fieldName: 'tenantId',
      label: $t('pages.multiTenant.fields.tenantId'),
      rules: z
        .number({ message: $t('pages.multiTenant.domain.messages.tenantRequired') })
        .int()
        .min(1, { message: $t('pages.multiTenant.domain.messages.tenantRequired') }),
    },
    {
      component: 'Input',
      componentProps: {
        'data-testid': 'domain-input',
      },
      fieldName: 'domain',
      label: $t('pages.multiTenant.domain.fields.domain'),
      rules: z
        .string()
        .min(1, { message: $t('pages.multiTenant.domain.messages.domainRequired') })
        .max(255, { message: $t('pages.multiTenant.domain.messages.domainRule') }),
    },
    {
      component: 'Switch',
      componentProps: {
        'data-testid': 'domain-primary-input',
      },
      fieldName: 'isPrimary',
      label: $t('pages.multiTenant.domain.fields.isPrimary'),
    },
  ],
  showDefaultActions: false,
});

const [Modal, modalApi] = useVbenModal({
  fullscreenButton: false,
  onClosed: handleClosed,
  onConfirm: handleConfirm,
  onOpenChange: handleOpenChange,
});

async function handleConfirm() {
  try {
    modalApi.lock(true);
    const { valid } = await formApi.validate();
    if (!valid) {
      return;
    }
    const values = await formApi.getValues<DomainFormValues>();
    await platformDomainCreate({
      tenantId: Number(values.tenantId),
      domain: values.domain,
      isPrimary: Boolean(values.isPrimary),
    });
    message.success($t('pages.common.createSuccess'));
    emit('success');
    modalApi.close();
  } finally {
    modalApi.lock(false);
  }
}

async function handleOpenChange(open: boolean) {
  if (!open) {
    return;
  }
  await formApi.setValues(formValues);
}

async function handleClosed() {
  await formApi.resetForm();
}
</script>

<template>
  <Modal :title="$t('pages.multiTenant.domain.actions.create')">
    <div data-testid="domain-form">
      <DomainForm />
    </div>
  </Modal>
</template>
