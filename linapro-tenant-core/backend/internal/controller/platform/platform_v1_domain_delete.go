// This file implements the platform tenant-domain delete endpoint.

package platform

import (
	"context"

	v1 "lina-plugin-linapro-tenant-core/backend/api/platform/v1"
)

// DomainDelete soft-deletes a tenant domain mapping by ID.
func (c *ControllerV1) DomainDelete(ctx context.Context, req *v1.DomainDeleteReq) (res *v1.DomainDeleteRes, err error) {
	if err := c.domainSvc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DomainDeleteRes{}, nil
}
