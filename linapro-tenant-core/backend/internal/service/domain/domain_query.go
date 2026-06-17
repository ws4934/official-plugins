// This file implements tenant domain list queries with database-side filtering,
// pagination, and a minimal projection.

package domain

import (
	"context"
	"strings"

	"lina-plugin-linapro-tenant-core/backend/internal/service/shared"
)

// defaultPageNum and defaultPageSize bound list queries when the caller omits them.
const (
	defaultPageNum  = 1
	defaultPageSize = 10
)

// List queries domain mappings with database-side tenant, domain, and status
// filtering and pagination, returning a minimal projection and total count.
func (s *serviceImpl) List(ctx context.Context, in ListInput) (*ListOutput, error) {
	model := shared.Model(ctx, shared.TableDomain)
	if in.TenantId > 0 {
		model = model.Where("tenant_id", in.TenantId)
	}
	if domain := strings.ToLower(strings.TrimSpace(in.Domain)); domain != "" {
		model = model.WhereLike("domain", "%"+domain+"%")
	}
	if status := strings.TrimSpace(in.Status); status != "" {
		model = model.Where("status", status)
	}
	total, err := model.Count()
	if err != nil {
		return nil, err
	}
	pageNum := in.PageNum
	if pageNum < 1 {
		pageNum = defaultPageNum
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	rows := make([]*Entity, 0, pageSize)
	if err := model.
		Fields("id, tenant_id, domain, is_primary, is_verified, status, created_at").
		Page(pageNum, pageSize).
		OrderDesc("id").
		Scan(&rows); err != nil {
		return nil, err
	}
	return &ListOutput{List: rows, Total: total}, nil
}
