// resolver_domain.go implements the custom-domain tenant resolver. It maps a
// full request host to a tenant through verified, active domain mappings and is
// the basis for storefront-by-domain resolution. It never resolves to the
// platform tenant: an unmatched, unverified, disabled, or inactive-tenant host
// reports no match so the resolution chain can fail closed.

package resolver

import (
	"context"
	"database/sql"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	"lina-plugin-linapro-tenant-core/backend/internal/service/shared"
)

// domainResolver resolves a tenant by matching the full request host against
// verified, active domain mappings owned by an active tenant.
type domainResolver struct{}

// Name returns the configured resolver name.
func (r domainResolver) Name() string {
	return shared.ResolverDomain
}

// Resolve matches the normalized request host against verified, active domain
// mappings. It returns the mapped tenant only when the domain is verified and
// active and its tenant is active. An empty host, unknown host, unverified or
// disabled domain, or inactive tenant reports no match without returning the
// platform tenant. It returns an error only on a database access failure.
func (r domainResolver) Resolve(ctx context.Context, request *ghttp.Request, identity Identity, config Config) (*Result, bool, error) {
	host := normalizeHost(request.Host)
	if host == "" {
		return nil, false, nil
	}
	domainRow := struct {
		TenantId int64 `orm:"tenant_id"`
	}{}
	err := shared.Model(ctx, shared.TableDomain).
		Fields("tenant_id").
		Where("domain", host).
		Where("is_verified", true).
		Where("status", string(shared.DomainStatusActive)).
		Scan(&domainRow)
	if err != nil {
		if gerror.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, true, err
	}
	if domainRow.TenantId == 0 {
		return nil, false, nil
	}
	active, err := tenantIsActive(ctx, domainRow.TenantId)
	if err != nil {
		return nil, true, err
	}
	if !active {
		return nil, false, nil
	}
	return &Result{TenantID: domainRow.TenantId, Source: shared.ResolverDomain}, true, nil
}

// tenantIsActive reports whether the tenant exists and is in the active status.
// Soft-deleted tenants are excluded by GoFrame automatic soft-delete filtering.
func tenantIsActive(ctx context.Context, tenantID int64) (bool, error) {
	count, err := shared.Model(ctx, shared.TableTenant).
		Where("id", tenantID).
		Where("status", string(shared.TenantStatusActive)).
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// normalizeHost lowercases the request host and strips any port suffix.
func normalizeHost(host string) string {
	hostname := strings.ToLower(strings.TrimSpace(host))
	if colon := strings.LastIndex(hostname, ":"); colon >= 0 {
		hostname = hostname[:colon]
	}
	return hostname
}
