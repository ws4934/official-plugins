// Package domain implements tenant domain mapping CRUD and verification for the
// linapro-tenant-core source plugin. Domains map request hosts to tenants and
// are platform-governed data: management is gated by platform permission rather
// than row-level tenant data scope, matching tenant master-data management.
package domain

import (
	"context"

	"lina-core/pkg/plugin/capability/bizctxcap"
)

// Service defines tenant domain mapping CRUD and verification for plugin-owned
// domain rows.
type Service interface {
	// List queries domain mappings with pagination and tenant/domain/status
	// filters applied at the database layer. It is read-only and returns database
	// errors.
	List(ctx context.Context, in ListInput) (*ListOutput, error)
	// Create maps a normalized, globally unique domain to a tenant and records
	// creator metadata from ctx. New mappings start unverified. It returns
	// CodeDomainInvalid, CodeDomainAlreadyExists, or database errors.
	Create(ctx context.Context, in CreateInput) (int64, error)
	// Delete soft-deletes a visible domain mapping or returns CodeDomainNotFound.
	Delete(ctx context.Context, id int64) error
	// SetVerified sets the verification flag of a visible domain mapping or
	// returns CodeDomainNotFound. Only verified, active domains resolve tenants.
	SetVerified(ctx context.Context, id int64, verified bool) error
}

// Ensure serviceImpl implements Service.
var _ Service = (*serviceImpl)(nil)

// serviceImpl implements Service.
type serviceImpl struct {
	bizCtxSvc bizctxcap.Service
}

// New creates and returns a new domain Service instance.
func New(bizCtxSvc bizctxcap.Service) Service {
	return &serviceImpl{bizCtxSvc: bizCtxSvc}
}

// Entity is the service-layer domain projection.
type Entity struct {
	Id         int64  `json:"id" orm:"id"`
	TenantId   int64  `json:"tenantId" orm:"tenant_id"`
	Domain     string `json:"domain" orm:"domain"`
	IsPrimary  bool   `json:"isPrimary" orm:"is_primary"`
	IsVerified bool   `json:"isVerified" orm:"is_verified"`
	Status     string `json:"status" orm:"status"`
	CreatedAt  string `json:"createdAt" orm:"created_at"`
}

// ListInput defines domain list filters.
type ListInput struct {
	PageNum  int
	PageSize int
	TenantId int64
	Domain   string
	Status   string
}

// ListOutput defines domain list output.
type ListOutput struct {
	List  []*Entity
	Total int
}

// CreateInput defines domain creation fields.
type CreateInput struct {
	TenantId  int64
	Domain    string
	IsPrimary bool
}

// domainInsertData is a typed insert payload for plugin_linapro_tenant_core_domain.
type domainInsertData struct {
	TenantId          int64  `orm:"tenant_id"`
	Domain            string `orm:"domain"`
	IsPrimary         bool   `orm:"is_primary"`
	IsVerified        bool   `orm:"is_verified"`
	VerificationToken string `orm:"verification_token"`
	Status            string `orm:"status"`
	CreatedBy         int64  `orm:"created_by"`
	UpdatedBy         int64  `orm:"updated_by"`
}

// domainVerifyData is a typed update payload for the domain verification flag.
type domainVerifyData struct {
	IsVerified bool  `orm:"is_verified"`
	UpdatedBy  int64 `orm:"updated_by"`
}
