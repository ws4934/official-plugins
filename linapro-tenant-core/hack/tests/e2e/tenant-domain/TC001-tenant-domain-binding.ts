import { test } from "../../support/linapro-tenant-core";

import { TenantDomainBindingPage } from "../../pages/TenantDomainBindingPage";

test.describe("TC-1 租户编辑弹窗内域名绑定", () => {
  test("TC-1a: 打开租户编辑弹窗显示域名区与已绑定域名（含 i18n 文案）", async ({
    page,
  }) => {
    const domainPage = new TenantDomainBindingPage(page);
    await domainPage.gotoTenants();
    await domainPage.openTenantEditDomains();
  });

  test("TC-1b: 在弹窗内新增域名并出现在列表", async ({ page }) => {
    const domainPage = new TenantDomainBindingPage(page);
    await domainPage.gotoTenants();
    await domainPage.openTenantEditDomains();
    await domainPage.exerciseAddDomain();
  });

  test("TC-1c: 切换验证开关写入目标验证状态", async ({ page }) => {
    const domainPage = new TenantDomainBindingPage(page);
    await domainPage.gotoTenants();
    await domainPage.openTenantEditDomains();
    await domainPage.exerciseVerifyToggle();
  });

  test("TC-1d: 删除域名后对应行消失", async ({ page }) => {
    const domainPage = new TenantDomainBindingPage(page);
    await domainPage.gotoTenants();
    await domainPage.openTenantEditDomains();
    await domainPage.exerciseDeleteDomain();
  });
});
