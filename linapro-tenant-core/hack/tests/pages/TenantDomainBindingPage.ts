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

const tenant = {
  id: 101,
  code: "alpha",
  name: "Alpha BU",
  remark: "Alpha tenant remark",
  status: "active",
  createdAt: "2026-05-09 10:00:00",
};

const initialDomains: DomainRow[] = [
  {
    id: 1,
    tenantId: 101,
    domain: "shop.alpha.com",
    isPrimary: true,
    isVerified: false,
    status: "active",
    createdAt: 1_717_000_000_000,
  },
];

function ok(data: unknown) {
  return {
    contentType: "application/json",
    body: JSON.stringify({ code: 0, data }),
  };
}

function routePath(url: string) {
  return new URL(url).pathname
    .replace(/^\/x\/linapro-tenant-core\/api\/v1/, "")
    .replace(/^\/api\/v1/, "")
    .replace(/^\/api/, "");
}

// TenantDomainBindingPage drives the tenant management page and verifies domain
// binding inside the tenant edit modal (list, add, verification toggle, delete)
// against mocked workbench, tenant, and domain APIs, with translated i18n text.
export class TenantDomainBindingPage {
  private rows: DomainRow[] = initialDomains.map((row) => ({ ...row }));
  private lastCreatePayload: Record<string, any> | null = null;
  private lastVerifyId = 0;
  private lastVerifyPayload: Record<string, any> | null = null;
  private lastDeleteId = 0;

  constructor(private page: Page) {}

  private menuTree() {
    return [
      {
        children: [
          {
            children: [],
            meta: { order: 10 },
            name: "PlatformTenantManagement",
            path: "/platform/tenants",
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
            component: "#/views/system/plugin/dynamic-page",
            meta: {
              icon: "lucide:building",
              order: 10,
              title: "plugin.linapro-tenant-core.menu.platform.tenants.name",
            },
            name: "PluginMultiTenantPlatformTenants",
            path: "/platform/tenants",
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
      homePath: "/platform/tenants",
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

  private async mockTenantApis() {
    await this.page.route(
      /(?:\/api(?:\/v1)?|\/x\/linapro-tenant-core\/api\/v1)\/platform\/tenants(?:\/\d+)?(?:\?.*)?$/,
      async (route) => {
        const method = route.request().method();
        if (method === "GET") {
          await route.fulfill(ok({ list: [tenant], total: 1 }));
          return;
        }
        if (method === "PUT") {
          await route.fulfill(ok({ ...tenant, ...route.request().postDataJSON() }));
          return;
        }
        await route.fulfill(ok(tenant));
      },
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
          const tenantId = new URL(route.request().url()).searchParams.get(
            "tenantId",
          );
          const list = tenantId
            ? this.rows.filter((row) => row.tenantId === Number(tenantId))
            : this.rows;
          await route.fulfill(ok({ list, total: list.length }));
          return;
        }
        if (method === "POST") {
          const body = route.request().postDataJSON();
          this.lastCreatePayload = body;
          const id = 100 + this.rows.length;
          this.rows.push({
            id,
            tenantId: Number(body.tenantId),
            domain: String(body.domain ?? "").trim().toLowerCase(),
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
          this.lastVerifyId = id;
          this.lastVerifyPayload = route.request().postDataJSON();
          const row = this.rows.find((item) => item.id === id);
          if (row) {
            row.isVerified = Boolean(this.lastVerifyPayload?.verified);
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

  async gotoTenants() {
    await this.installAuthState();
    await this.mockShell();
    await this.mockTenantApis();
    await this.mockDomainApis();
    await this.page.goto("/platform/tenants");
    await waitForRouteReady(this.page);
  }

  async openTenantEditDomains() {
    await expect(this.page.getByTestId("platform-tenants-page")).toBeVisible();
    await this.page.getByTestId("tenant-edit-101").click();
    await expect(this.page.getByTestId("tenant-form")).toBeVisible();
    const domains = this.page.getByTestId("tenant-domains");
    await expect(domains).toBeVisible();
    // Translated section title proves i18n resolved (not the raw key).
    await expect(domains).toContainText("域名");
    await expect(
      this.page.getByTestId("tenant-domains-section").getByText("shop.alpha.com"),
    ).toBeVisible();
  }

  async exerciseAddDomain() {
    await this.page.getByTestId("tenant-domain-input").fill("acme.example.com");
    const createResponse = this.page.waitForResponse(
      (response) =>
        response.request().method() === "POST" &&
        routePath(response.url()) === "/platform/domains",
    );
    await this.page.getByTestId("tenant-domain-add").click();
    await createResponse;
    expect(this.lastCreatePayload).toMatchObject({
      domain: "acme.example.com",
      tenantId: 101,
    });
    await expect(
      this.page
        .getByTestId("tenant-domains-section")
        .getByText("acme.example.com"),
    ).toBeVisible();
  }

  async exerciseVerifyToggle() {
    const verifyResponse = this.page.waitForResponse(
      (response) =>
        response.request().method() === "PUT" &&
        /\/platform\/domains\/1\/verification$/.test(routePath(response.url())),
    );
    await this.page.getByTestId("tenant-domain-verify-1").click();
    await verifyResponse;
    expect(this.lastVerifyId).toBe(1);
    expect(this.lastVerifyPayload).toMatchObject({ verified: true });
  }

  async exerciseDeleteDomain() {
    await this.page.getByTestId("tenant-domain-delete-1").click();
    const deleteResponse = this.page.waitForResponse(
      (response) =>
        response.request().method() === "DELETE" &&
        /\/platform\/domains\/1$/.test(routePath(response.url())),
    );
    await this.page
      .getByRole("button", { name: /^确\s*定$|^OK$|^Confirm$/ })
      .click();
    await deleteResponse;
    expect(this.lastDeleteId).toBe(1);
    await expect(
      this.page
        .getByTestId("tenant-domains-section")
        .getByText("shop.alpha.com"),
    ).toHaveCount(0);
  }
}
