// This file implements the platform tenant-domain verification endpoint.

package platform

import (
	"context"

	v1 "lina-plugin-linapro-tenant-core/backend/api/platform/v1"
)

// DomainVerify sets the verification flag of a tenant domain mapping.
func (c *ControllerV1) DomainVerify(ctx context.Context, req *v1.DomainVerifyReq) (res *v1.DomainVerifyRes, err error) {
	if err := c.domainSvc.SetVerified(ctx, req.Id, req.Verified); err != nil {
		return nil, err
	}
	return &v1.DomainVerifyRes{}, nil
}
