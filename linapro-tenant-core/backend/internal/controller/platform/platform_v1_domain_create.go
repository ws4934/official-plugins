// This file implements the platform tenant-domain create endpoint.

package platform

import (
	"context"

	v1 "lina-plugin-linapro-tenant-core/backend/api/platform/v1"
	domainsvc "lina-plugin-linapro-tenant-core/backend/internal/service/domain"
)

// DomainCreate maps a domain host to a tenant.
func (c *ControllerV1) DomainCreate(ctx context.Context, req *v1.DomainCreateReq) (res *v1.DomainCreateRes, err error) {
	id, err := c.domainSvc.Create(ctx, domainsvc.CreateInput{
		TenantId:  req.TenantId,
		Domain:    req.Domain,
		IsPrimary: req.IsPrimary,
	})
	if err != nil {
		return nil, err
	}
	return &v1.DomainCreateRes{Id: id}, nil
}
