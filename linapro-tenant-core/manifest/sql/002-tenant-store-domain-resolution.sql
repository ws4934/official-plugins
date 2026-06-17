-- Purpose: Stores tenant domain mappings used to resolve a tenant by request host, including custom domains, primary flag, verification state, and lifecycle status.
-- 用途：存储租户域名映射，用于按请求 Host 解析租户，包含自定义域名、主域名标记、验证状态与生命周期状态。
CREATE TABLE IF NOT EXISTS plugin_linapro_tenant_core_domain (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "tenant_id" BIGINT NOT NULL,
    "domain" VARCHAR(255) NOT NULL,
    "is_primary" BOOLEAN NOT NULL DEFAULT FALSE,
    "is_verified" BOOLEAN NOT NULL DEFAULT FALSE,
    "verification_token" VARCHAR(128) NOT NULL DEFAULT '',
    "status" VARCHAR(32) NOT NULL DEFAULT 'active',
    "created_by" BIGINT NOT NULL DEFAULT 0,
    "updated_by" BIGINT NOT NULL DEFAULT 0,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP,
    CONSTRAINT uk_plugin_linapro_tenant_core_domain_domain UNIQUE ("domain")
);

-- Tenant lookup index for listing and management filters scoped by tenant.
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_tenant_core_domain_tenant
    ON plugin_linapro_tenant_core_domain ("tenant_id");

-- Resolution index supporting host-only lookups that require verified, active domains.
CREATE INDEX IF NOT EXISTS idx_plugin_linapro_tenant_core_domain_resolve
    ON plugin_linapro_tenant_core_domain ("is_verified", "status");
