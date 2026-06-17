import type { Page } from "@host-tests/support/playwright";

import { expect } from "@host-tests/support/playwright";

import { waitForRouteReady } from "@host-tests/support/ui";

interface DomainRow {
  id: number;
  tenantId: number;
  domain: string;
  isPrimary: boolean;
  isVerified: boolean;
  status: string;
  createdAt: number;
}

const initialDomains: DomainRow[] = [
  {
    id: 1,
    tenantId: 101,
    domain: "shop.alpha.com",
    isPrimary: true,
    isVerified: true,
    status: "active",
    createdAt: 1_717_000_000_000,
  },
  {
    id: 2,
    tenantId: 102,
    domain: "shop.beta.com",
    isPrimary: false,
    isVerified: false,
    status: "active",
    createdAt: 1_717_100_000_000,
  },
];

function ok(data: unknown) {
  return {
    contentType: "application/json",
    body: JSON.stringify({ code: 0, data }),
  };
}

function fail(errorCode: string, message: string) {
  return {
    status: 400,
    contentType: "application/json",
    body: JSON.stringify({ code: errorCode, message }),
  };
}

function routePath(url: string) {
  return new URL(url).pathname
    .replace(/^\/x\/linapro-tenant-core\/api\/v1/, "")
    .replace(/^\/api\/v1/, "")
    .replace(/^\/api/, "");
}

// DomainManagementPage drives the tenant domain management page with mocked
// workbench shell and domain CRUD APIs, so browser tests verify list, create,
// verification toggle, delete, the unique-domain failure path, and translated
// i18n text without a live backend.
export class DomainManagementPage {
  private rows: DomainRow[] = initialDomains.map((row) => ({ ...row }));
  private lastCreatePayload: Record<string, any> | null = null;
  private lastVerifyPayload: Record<string, any> | null = null;
  private lastVerifyId = 0;
  private lastDeleteId = 0;

  constructor(private page: Page) {}

  private menuTree() {
    return [
      {
        children: [
          {
            children: [],
            meta: { order: 1 },
            name: "Analytics",
            path: "/dashboard/analytics",
          },
        ],
        meta: { order: -1 },
        name: "Dashboard",
        path: "/dashboard",
      },
      {
        children: [
          {
            children: [],
            meta: { order: 20 },
            name: "PlatformDomains",
            path: "/platform/domains",
          },
        ],
        meta: { order: 15 },
        name: "Platform",
        path: "/platform",
      },
    ];
  }

  private routeTree() {
    return [
      {
        children: [
          {
            children: [],
            component: "#/views/dashboard/analytics/index.vue",
            meta: { icon: "lucide:circle", order: 1, title: "page.dashboard.analytics" },
            name: "Analytics",
            path: "/dashboard/analytics",
          },
        ],
        component: "",
        meta: { icon: "lucide:circle", order: -1, title: "page.dashboard.title" },
        name: "Dashboard",
        path: "/dashboard",
      },
      {
        children: [
          {
            children: [],
            component: "#/views/system/plugin/dynamic-page",
            meta: {
              icon: "lucide:globe",
              order: 20,
              title: "plugin.linapro-tenant-core.menu.platform.domains.name",
            },
            name: "PluginMultiTenantPlatformDomains",
            path: "/platform/domains",
          },
        ],
        component: "",
        meta: {
          icon: "lucide:circle",
          order: 15,
          tenantAccessMode: "platform",
          title: "menu.platform.title",
        },
        name: "Platform",
        path: "/platform",
      },
    ];
  }

  private userInfo() {
    return {
      userId: 1,
      username: "admin",
      realName: "Administrator",
      avatar: "",
      email: "admin@example.com",
      permissions: ["*"],
      roles: ["admin"],
      homePath: "/platform/domains",
      menus: this.menuTree(),
    };
  }

  private async installAuthState() {
    await this.page.addInitScript(() => {
      localStorage.clear();
      localStorage.setItem(
        "lina-web-antd-5.6.0-dev-core-access",
        JSON.stringify({
          accessToken: "platform-token",
          accessCodes: ["*"],
          refreshToken: null,
          isLockScreen: false,
        }),
      );
      localStorage.setItem(
        "linapro:tenant-state",
        JSON.stringify({
          enabled: true,
          currentTenant: null,
          tenants: [],
          impersonation: { active: false },
        }),
      );
    });
  }

  private async mockShell() {
    await this.page.route(
      /\/api(?:\/v1)?\/config\/public\/frontend$/,
      (route) => route.fulfill(ok({})),
    );
    await this.page.route(
      /\/api(?:\/v1)?\/i18n\/runtime\/locales(?:\?.*)?$/,
      (route) =>
        route.fulfill(
          ok({
            enabled: true,
            locale: "zh-CN",
            items: [
              { locale: "zh-CN", name: "简体中文", nativeName: "简体中文", isDefault: true },
              { locale: "en-US", name: "English", nativeName: "English" },
            ],
          }),
        ),
    );
    await this.page.route(
      /\/api(?:\/v1)?\/i18n\/runtime\/messages(?:\?.*)?$/,
      (route) => route.fulfill(ok({ messages: {} })),
    );
    await this.page.route(
      /\/api(?:\/v1)?\/plugins\/dynamic(?:\?.*)?$/,
      (route) =>
        route.fulfill(
          ok({
            list: [
              {
                enabled: 1,
                generation: 1,
                id: "linapro-tenant-core",
                installed: 1,
                statusKey: "sys_plugin.status:linapro-tenant-core",
                version: "v0.1.0",
              },
            ],
          }),
        ),
    );
    await this.page.route(
      /\/api(?:\/v1)?\/user\/message\/count(?:\?.*)?$/,
      (route) => route.fulfill(ok({ count: 0 })),
    );
    await this.page.route(
      /\/api(?:\/v1)?\/user\/message(?:\?.*)?$/,
      (route) => route.fulfill(ok({ list: [], total: 0 })),
    );
    await this.page.route(/\/api(?:\/v1)?\/user\/info$/, (route) =>
      route.fulfill(ok(this.userInfo())),
    );
    await this.page.route(/\/api(?:\/v1)?\/menus\/all(?:\?.*)?$/, (route) =>
      route.fulfill(ok({ list: this.routeTree() })),
    );
  }

  private async mockDomainApis() {
    await this.page.route(
      /(?:\/api(?:\/v1)?|\/x\/linapro-tenant-core\/api\/v1)\/platform\/domains(?:\/\d+(?:\/verification)?)?(?:\?.*)?$/,
      async (route) => {
        const method = route.request().method();
        const path = routePath(route.request().url());
        const segments = path.split("/").filter(Boolean);
        if (method === "GET") {
          await route.fulfill(ok({ list: this.rows, total: this.rows.length }));
          return;
        }
        if (method === "POST") {
          const body = route.request().postDataJSON();
          this.lastCreatePayload = body;
          const normalized = String(body.domain ?? "").trim().toLowerCase();
          if (this.rows.some((row) => row.domain === normalized)) {
            await route.fulfill(
              fail("MULTI_TENANT_DOMAIN_ALREADY_EXISTS", "域名已被映射到租户"),
            );
            return;
          }
          const id = 100 + this.rows.length;
          this.rows.push({
            id,
            tenantId: Number(body.tenantId),
            domain: normalized,
            isPrimary: Boolean(body.isPrimary),
            isVerified: false,
            status: "active",
            createdAt: 1_717_200_000_000,
          });
          await route.fulfill(ok({ id }));
          return;
        }
        if (method === "PUT" && path.endsWith("/verification")) {
          const id = Number(segments[2]);
          const body = route.request().postDataJSON();
          this.lastVerifyPayload = body;
          this.lastVerifyId = id;
          const row = this.rows.find((item) => item.id === id);
          if (row) {
            row.isVerified = Boolean(body.verified);
          }
          await route.fulfill(ok({}));
          return;
        }
        if (method === "DELETE") {
          const id = Number(segments[2]);
          this.lastDeleteId = id;
          this.rows = this.rows.filter((item) => item.id !== id);
          await route.fulfill(ok({}));
          return;
        }
        await route.fulfill(ok({}));
      },
    );
  }

  async goto() {
    await this.installAuthState();
    await this.mockShell();
    await this.mockDomainApis();
    await this.page.goto("/platform/domains");
    await waitForRouteReady(this.page);
  }

  async expectDomainListWithTranslations() {
    const root = this.page.getByTestId("platform-domains-page");
    await expect(root).toBeVisible();
    await expect(
      this.page.getByText(/插件页面未找到|Plugin page not found/),
    ).toHaveCount(0);
    // Translated column header proves i18n resolved (not the raw key).
    await expect(
      this.page.locator(".vxe-header--column", { hasText: "域名" }).first(),
    ).toBeVisible();
    await expect(root.getByText("shop.alpha.com")).toBeVisible();
    await expect(root.getByText("shop.beta.com")).toBeVisible();
    // Translated add button (pages.common.add).
    await expect(this.page.getByTestId("domain-create")).toHaveText(/^新\s*增$/);
  }

  async exerciseCreate() {
    await this.page.getByTestId("domain-create").click();
    await expect(this.page.getByTestId("domain-form")).toBeVisible();
    await this.page.getByTestId("domain-tenant-input").locator("input").fill("101");
    await this.page.getByTestId("domain-input").fill("shop.gamma.com");

    const createResponse = this.page.waitForResponse(
      (response) =>
        response.request().method() === "POST" &&
        routePath(response.url()) === "/platform/domains",
    );
    await this.page
      .getByRole("dialog")
      .getByRole("button", { name: /^确\s*认$|^Confirm$/ })
      .click();
    await createResponse;
    expect(this.lastCreatePayload).toMatchObject({
      domain: "shop.gamma.com",
      tenantId: 101,
    });
    await expect(this.page.getByTestId("domain-form")).toHaveCount(0);
    await expect(
      this.page.getByTestId("platform-domains-page").getByText("shop.gamma.com"),
    ).toBeVisible();
  }

  async exerciseDuplicateRejected() {
    await this.page.getByTestId("domain-create").click();
    await expect(this.page.getByTestId("domain-form")).toBeVisible();
    await this.page.getByTestId("domain-tenant-input").locator("input").fill("101");
    await this.page.getByTestId("domain-input").fill("shop.alpha.com");

    const createResponse = this.page.waitForResponse(
      (response) =>
        response.request().method() === "POST" &&
        routePath(response.url()) === "/platform/domains",
    );
    await this.page
      .getByRole("dialog")
      .getByRole("button", { name: /^确\s*认$|^Confirm$/ })
      .click();
    const response = await createResponse;
    expect(response.status()).toBe(400);
    // Duplicate is rejected: no new row and the error toast keeps the form open.
    await expect(
      this.page
        .getByTestId("platform-domains-page")
        .getByText("shop.alpha.com"),
    ).toHaveCount(1);
  }

  async exerciseVerifyToggle() {
    const verifyResponse = this.page.waitForResponse(
      (response) =>
        response.request().method() === "PUT" &&
        /\/platform\/domains\/2\/verification$/.test(routePath(response.url())),
    );
    await this.page.getByTestId("domain-verify-2").click();
    await verifyResponse;
    expect(this.lastVerifyId).toBe(2);
    expect(this.lastVerifyPayload).toMatchObject({ verified: true });
  }

  async exerciseDelete() {
    await this.page.getByTestId("domain-delete-2").click();
    const deleteResponse = this.page.waitForResponse(
      (response) =>
        response.request().method() === "DELETE" &&
        /\/platform\/domains\/2$/.test(routePath(response.url())),
    );
    await this.page
      .getByRole("button", { name: /^确\s*定$|^OK$|^Confirm$/ })
      .click();
    await deleteResponse;
    expect(this.lastDeleteId).toBe(2);
    await expect(
      this.page.getByTestId("platform-domains-page").getByText("shop.beta.com"),
    ).toHaveCount(0);
  }
}
