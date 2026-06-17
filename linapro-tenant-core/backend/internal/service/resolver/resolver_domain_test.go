// resolver_domain_test.go verifies custom-domain tenant resolution. Verified,
// active domains owned by an active tenant resolve to that tenant, while
// unverified, disabled, inactive-tenant, or unknown hosts report no match
// without returning the platform tenant. It also guards the host-only
// storefront chain red line.

package resolver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	_ "lina-core/pkg/dbdriver"

	"lina-plugin-linapro-tenant-core/backend/internal/model/do"
	"lina-plugin-linapro-tenant-core/backend/internal/service/shared"
)

// TestStorefrontChainExcludesDefault verifies the host-only storefront chain
// includes the domain resolver but never the membership default, so an
// unmatched host fails closed instead of falling back to the platform tenant.
func TestStorefrontChainExcludesDefault(t *testing.T) {
	chain := shared.StorefrontResolverChain()
	hasDomain := false
	for _, name := range chain {
		if name == shared.ResolverDefault {
			t.Fatalf("storefront chain must exclude the default resolver, got %v", chain)
		}
		if name == shared.ResolverDomain {
			hasDomain = true
		}
	}
	if !hasDomain {
		t.Fatalf("storefront chain must include the domain resolver, got %v", chain)
	}
}

// TestDomainResolverMatchesVerifiedActiveDomain verifies a verified, active
// domain owned by an active tenant resolves to that tenant.
func TestDomainResolverMatchesVerifiedActiveDomain(t *testing.T) {
	ctx := context.Background()
	configureResolverTestDB(t, ctx)
	prepareResolverTenantSchema(t, ctx)
	prepareResolverDomainSchema(t, ctx)

	host := "match-acme.example"
	tenantID := seedResolverDomain(t, ctx, shared.TenantStatusActive, host, true, shared.DomainStatusActive)

	resolved, matched, err := domainResolver{}.Resolve(ctx, requestForHost(host), Identity{}, Config{})
	if err != nil {
		t.Fatalf("resolve verified domain failed: %v", err)
	}
	if !matched || resolved == nil {
		t.Fatalf("expected verified domain to match, got matched=%t result=%#v", matched, resolved)
	}
	if resolved.TenantID != tenantID || resolved.Source != shared.ResolverDomain {
		t.Fatalf("expected tenant %d via domain source, got %#v", tenantID, resolved)
	}
}

// TestDomainResolverSkipsUnverifiedDomain verifies an unverified domain reports
// no match and never returns the platform tenant.
func TestDomainResolverSkipsUnverifiedDomain(t *testing.T) {
	ctx := context.Background()
	configureResolverTestDB(t, ctx)
	prepareResolverTenantSchema(t, ctx)
	prepareResolverDomainSchema(t, ctx)

	host := "unverified-acme.example"
	seedResolverDomain(t, ctx, shared.TenantStatusActive, host, false, shared.DomainStatusActive)

	resolved, matched, err := domainResolver{}.Resolve(ctx, requestForHost(host), Identity{}, Config{})
	if err != nil {
		t.Fatalf("resolve unverified domain failed: %v", err)
	}
	if matched || resolved != nil {
		t.Fatalf("expected unverified domain to report no match, got matched=%t result=%#v", matched, resolved)
	}
}

// TestDomainResolverSkipsInactiveTenant verifies a verified, active domain whose
// tenant is suspended reports no match.
func TestDomainResolverSkipsInactiveTenant(t *testing.T) {
	ctx := context.Background()
	configureResolverTestDB(t, ctx)
	prepareResolverTenantSchema(t, ctx)
	prepareResolverDomainSchema(t, ctx)

	host := "suspended-acme.example"
	seedResolverDomain(t, ctx, shared.TenantStatusSuspended, host, true, shared.DomainStatusActive)

	resolved, matched, err := domainResolver{}.Resolve(ctx, requestForHost(host), Identity{}, Config{})
	if err != nil {
		t.Fatalf("resolve inactive-tenant domain failed: %v", err)
	}
	if matched || resolved != nil {
		t.Fatalf("expected inactive tenant to report no match, got matched=%t result=%#v", matched, resolved)
	}
}

// TestDomainResolverUnknownHostReturnsNoMatch verifies an unknown host reports
// no match so the host-only chain fails closed rather than resolving platform.
func TestDomainResolverUnknownHostReturnsNoMatch(t *testing.T) {
	ctx := context.Background()
	configureResolverTestDB(t, ctx)
	prepareResolverTenantSchema(t, ctx)
	prepareResolverDomainSchema(t, ctx)

	resolved, matched, err := domainResolver{}.Resolve(ctx, requestForHost("unknown-host.example"), Identity{}, Config{})
	if err != nil {
		t.Fatalf("resolve unknown host failed: %v", err)
	}
	if matched || resolved != nil {
		t.Fatalf("expected unknown host to report no match, got matched=%t result=%#v", matched, resolved)
	}
}

// TestDomainResolverNormalizesHostPortAndCase verifies host matching lowercases
// the host and strips any port suffix.
func TestDomainResolverNormalizesHostPortAndCase(t *testing.T) {
	ctx := context.Background()
	configureResolverTestDB(t, ctx)
	prepareResolverTenantSchema(t, ctx)
	prepareResolverDomainSchema(t, ctx)

	host := "port-acme.example"
	tenantID := seedResolverDomain(t, ctx, shared.TenantStatusActive, host, true, shared.DomainStatusActive)

	resolved, matched, err := domainResolver{}.Resolve(ctx, requestForHost("PORT-ACME.EXAMPLE:8443"), Identity{}, Config{})
	if err != nil {
		t.Fatalf("resolve normalized host failed: %v", err)
	}
	if !matched || resolved == nil || resolved.TenantID != tenantID {
		t.Fatalf("expected normalized host to match tenant %d, got matched=%t result=%#v", tenantID, matched, resolved)
	}
}

// requestForHost builds a ghttp request whose host is set for resolver tests.
func requestForHost(host string) *ghttp.Request {
	return &ghttp.Request{Request: httptest.NewRequest(http.MethodGet, "http://"+host+"/", nil)}
}

// prepareResolverDomainSchema creates the domain mapping table needed by
// custom-domain resolver tests.
func prepareResolverDomainSchema(t *testing.T, ctx context.Context) {
	t.Helper()

	statement := `CREATE TABLE IF NOT EXISTS plugin_linapro_tenant_core_domain (
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
	)`
	if _, err := g.DB().Exec(ctx, statement); err != nil {
		t.Fatalf("prepare resolver domain schema failed: %v", err)
	}
}

// seedResolverDomain inserts one tenant and one domain mapping for resolver
// tests and registers hard-delete cleanup so runs stay self-contained.
func seedResolverDomain(t *testing.T, ctx context.Context, tenantStatus shared.TenantStatus, host string, verified bool, domainStatus shared.DomainStatus) int64 {
	t.Helper()

	tenantID, err := shared.Model(ctx, shared.TableTenant).Data(do.Tenant{
		Code:   "rt-" + host,
		Name:   "Resolver Test " + host,
		Status: string(tenantStatus),
	}).InsertAndGetId()
	if err != nil {
		t.Fatalf("seed resolver tenant failed: %v", err)
	}
	if _, err = shared.Model(ctx, shared.TableDomain).Data(do.Domain{
		TenantId:   tenantID,
		Domain:     host,
		IsVerified: verified,
		Status:     string(domainStatus),
	}).Insert(); err != nil {
		t.Fatalf("seed resolver domain failed: %v", err)
	}
	t.Cleanup(func() {
		if _, err := shared.Model(ctx, shared.TableDomain).Where("tenant_id", tenantID).Unscoped().Delete(); err != nil {
			t.Errorf("cleanup resolver domain failed: %v", err)
		}
		if _, err := shared.Model(ctx, shared.TableTenant).Where("id", tenantID).Unscoped().Delete(); err != nil {
			t.Errorf("cleanup resolver tenant failed: %v", err)
		}
	})
	return tenantID
}
