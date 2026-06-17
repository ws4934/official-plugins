// This file implements the platform tenant-domain list endpoint.

package platform

import (
	"context"

	v1 "lina-plugin-linapro-tenant-core/backend/api/platform/v1"
	domainsvc "lina-plugin-linapro-tenant-core/backend/internal/service/domain"
)

// DomainList lists tenant domain mappings with paging and filters.
func (c *ControllerV1) DomainList(ctx context.Context, req *v1.DomainListReq) (res *v1.DomainListRes, err error) {
	out, err := c.domainSvc.List(ctx, domainsvc.ListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		TenantId: req.TenantId,
		Domain:   req.Domain,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*v1.DomainItem, 0, len(out.List))
	for _, item := range out.List {
		items = append(items, toAPIDomain(item))
	}
	return &v1.DomainListRes{List: items, Total: out.Total}, nil
}
