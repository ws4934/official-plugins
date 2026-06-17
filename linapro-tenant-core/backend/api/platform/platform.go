// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package platform

import (
	"context"

	"lina-plugin-linapro-tenant-core/backend/api/platform/v1"
)

type IPlatformV1 interface {
	DomainList(ctx context.Context, req *v1.DomainListReq) (res *v1.DomainListRes, err error)
	DomainCreate(ctx context.Context, req *v1.DomainCreateReq) (res *v1.DomainCreateRes, err error)
	DomainDelete(ctx context.Context, req *v1.DomainDeleteReq) (res *v1.DomainDeleteRes, err error)
	DomainVerify(ctx context.Context, req *v1.DomainVerifyReq) (res *v1.DomainVerifyRes, err error)
	TenantList(ctx context.Context, req *v1.TenantListReq) (res *v1.TenantListRes, err error)
	TenantCreate(ctx context.Context, req *v1.TenantCreateReq) (res *v1.TenantCreateRes, err error)
	TenantDelete(ctx context.Context, req *v1.TenantDeleteReq) (res *v1.TenantDeleteRes, err error)
	TenantEndImpersonate(ctx context.Context, req *v1.TenantEndImpersonateReq) (res *v1.TenantEndImpersonateRes, err error)
	TenantGet(ctx context.Context, req *v1.TenantGetReq) (res *v1.TenantGetRes, err error)
	TenantImpersonate(ctx context.Context, req *v1.TenantImpersonateReq) (res *v1.TenantImpersonateRes, err error)
	TenantStatus(ctx context.Context, req *v1.TenantStatusReq) (res *v1.TenantStatusRes, err error)
	TenantUpdate(ctx context.Context, req *v1.TenantUpdateReq) (res *v1.TenantUpdateRes, err error)
}
