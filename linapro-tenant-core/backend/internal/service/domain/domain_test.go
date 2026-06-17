// domain_test.go verifies tenant domain mapping commands and queries: domain
// normalization, global uniqueness, invalid input, not-found handling,
// verification flips, and tenant-scoped list filtering, using the local
// PostgreSQL test database.

package domain

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"lina-core/pkg/bizerr"
	_ "lina-core/pkg/dbdriver"

	"lina-plugin-linapro-tenant-core/backend/internal/service/shared"
)

// TestCreateNormalizesAndRejectsDuplicate verifies domains are stored lowercase
// without port and that a duplicate domain in any case is rejected.
func TestCreateNormalizesAndRejectsDuplicate(t *testing.T) {
	ctx := context.Background()
	configureDomainTestDB(t, ctx)
	prepareDomainSchema(t, ctx)
	cleanupDomains(t, ctx, "shop.acme.com")
	svc := New(nil)

	id, err := svc.Create(ctx, CreateInput{TenantId: 7, Domain: "SHOP.Acme.COM:8443", IsPrimary: true})
	if err != nil {
		t.Fatalf("create domain failed: %v", err)
	}
	if id == 0 {
		t.Fatal("expected a created domain id")
	}
	stored, err := shared.Model(ctx, shared.TableDomain).Fields("domain").Where("id", id).Value()
	if err != nil {
		t.Fatalf("read stored domain failed: %v", err)
	}
	if stored.String() != "shop.acme.com" {
		t.Fatalf("expected normalized domain shop.acme.com, got %q", stored.String())
	}
	if _, err := svc.Create(ctx, CreateInput{TenantId: 9, Domain: "shop.acme.com"}); !bizerr.Is(err, CodeDomainAlreadyExists) {
		t.Fatalf("expected duplicate domain rejection, got %v", err)
	}
}

// TestCreateRejectsInvalid verifies missing tenant or empty domain is rejected.
func TestCreateRejectsInvalid(t *testing.T) {
	ctx := context.Background()
	configureDomainTestDB(t, ctx)
	prepareDomainSchema(t, ctx)
	svc := New(nil)

	if _, err := svc.Create(ctx, CreateInput{TenantId: 0, Domain: "x.example"}); !bizerr.Is(err, CodeDomainInvalid) {
		t.Fatalf("expected invalid for missing tenant, got %v", err)
	}
	if _, err := svc.Create(ctx, CreateInput{TenantId: 1, Domain: "   "}); !bizerr.Is(err, CodeDomainInvalid) {
		t.Fatalf("expected invalid for empty domain, got %v", err)
	}
}

// TestDeleteMissingReturnsNotFound verifies deleting an unknown mapping fails closed.
func TestDeleteMissingReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	configureDomainTestDB(t, ctx)
	prepareDomainSchema(t, ctx)
	svc := New(nil)

	if err := svc.Delete(ctx, 999999); !bizerr.Is(err, CodeDomainNotFound) {
		t.Fatalf("expected not found for missing domain, got %v", err)
	}
}

// TestSetVerifiedFlipsFlag verifies verification can be set on a mapping.
func TestSetVerifiedFlipsFlag(t *testing.T) {
	ctx := context.Background()
	configureDomainTestDB(t, ctx)
	prepareDomainSchema(t, ctx)
	cleanupDomains(t, ctx, "verify.example")
	svc := New(nil)

	id, err := svc.Create(ctx, CreateInput{TenantId: 3, Domain: "verify.example"})
	if err != nil {
		t.Fatalf("create domain failed: %v", err)
	}
	if err := svc.SetVerified(ctx, id, true); err != nil {
		t.Fatalf("set verified failed: %v", err)
	}
	verified, err := shared.Model(ctx, shared.TableDomain).Fields("is_verified").Where("id", id).Value()
	if err != nil {
		t.Fatalf("read verified flag failed: %v", err)
	}
	if !verified.Bool() {
		t.Fatal("expected domain to be verified")
	}
}

// TestListFiltersByTenant verifies list filtering scopes results to one tenant.
func TestListFiltersByTenant(t *testing.T) {
	ctx := context.Background()
	configureDomainTestDB(t, ctx)
	prepareDomainSchema(t, ctx)
	cleanupDomains(t, ctx, "t1.example", "t2.example")
	svc := New(nil)

	if _, err := svc.Create(ctx, CreateInput{TenantId: 101, Domain: "t1.example"}); err != nil {
		t.Fatalf("seed tenant 101 domain failed: %v", err)
	}
	if _, err := svc.Create(ctx, CreateInput{TenantId: 202, Domain: "t2.example"}); err != nil {
		t.Fatalf("seed tenant 202 domain failed: %v", err)
	}
	out, err := svc.List(ctx, ListInput{TenantId: 101})
	if err != nil {
		t.Fatalf("list domains failed: %v", err)
	}
	if out.Total != 1 || len(out.List) != 1 || out.List[0].Domain != "t1.example" {
		t.Fatalf("expected only tenant 101 domain, got total=%d list=%#v", out.Total, out.List)
	}
}

// configureDomainTestDB points the package test at the local PostgreSQL database
// initialized by the repository test workflow and restores config on cleanup.
func configureDomainTestDB(t *testing.T, ctx context.Context) {
	t.Helper()

	originalConfig := gdb.GetAllConfig()
	if err := gdb.SetConfig(gdb.Config{
		gdb.DefaultGroupName: gdb.ConfigGroup{{Link: "pgsql:postgres:postgres@tcp(127.0.0.1:5432)/linapro?sslmode=disable"}},
	}); err != nil {
		t.Fatalf("configure domain test database failed: %v", err)
	}
	db := g.DB()
	t.Cleanup(func() {
		if err := db.Close(ctx); err != nil {
			t.Errorf("close domain test database failed: %v", err)
		}
		if err := gdb.SetConfig(originalConfig); err != nil {
			t.Errorf("restore domain test database config failed: %v", err)
		}
	})
}

// prepareDomainSchema creates the domain mapping table needed by domain tests.
func prepareDomainSchema(t *testing.T, ctx context.Context) {
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
		t.Fatalf("prepare domain schema failed: %v", err)
	}
}

// cleanupDomains registers hard-delete cleanup for the named domains so unique
// constraints do not accumulate across runs.
func cleanupDomains(t *testing.T, ctx context.Context, domains ...string) {
	t.Helper()

	t.Cleanup(func() {
		for _, domain := range domains {
			if _, err := shared.Model(ctx, shared.TableDomain).Where("domain", domain).Unscoped().Delete(); err != nil {
				t.Errorf("cleanup domain %q failed: %v", domain, err)
			}
		}
	})
}
