import { test } from "../../support/linapro-tenant-core";

import { DomainManagementPage } from "../../pages/DomainManagementPage";

test.describe("TC-1 租户域名管理", () => {
  test("TC-1a: 列表渲染已映射域名并显示翻译后的列标题与按钮文案", async ({
    page,
  }) => {
    const domainPage = new DomainManagementPage(page);
    await domainPage.goto();
    await domainPage.expectDomainListWithTranslations();
  });

  test("TC-1b: 创建新域名映射后出现在列表", async ({ page }) => {
    const domainPage = new DomainManagementPage(page);
    await domainPage.goto();
    await domainPage.exerciseCreate();
  });

  test("TC-1c: 重复域名被唯一约束拒绝且不新增行", async ({ page }) => {
    const domainPage = new DomainManagementPage(page);
    await domainPage.goto();
    await domainPage.exerciseDuplicateRejected();
  });

  test("TC-1d: 切换验证开关写入目标验证状态", async ({ page }) => {
    const domainPage = new DomainManagementPage(page);
    await domainPage.goto();
    await domainPage.exerciseVerifyToggle();
  });

  test("TC-1e: 删除域名映射后对应行消失", async ({ page }) => {
    const domainPage = new DomainManagementPage(page);
    await domainPage.goto();
    await domainPage.exerciseDelete();
  });
});
