// This file implements tenant domain create, delete, and verification commands,
// including domain normalization, uniqueness checks, and actor attribution.

package domain

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/util/guid"

	"lina-core/pkg/bizerr"
	"lina-plugin-linapro-tenant-core/backend/internal/service/shared"
)

// Create maps a normalized, unique domain to a tenant. The mapping starts
// unverified and carries a generated verification token for later proof flows.
func (s *serviceImpl) Create(ctx context.Context, in CreateInput) (int64, error) {
	domain := normalizeDomain(in.Domain)
	if domain == "" || in.TenantId <= 0 {
		return 0, bizerr.NewCode(CodeDomainInvalid)
	}
	existing, err := shared.Model(ctx, shared.TableDomain).Where("domain", domain).Count()
	if err != nil {
		return 0, err
	}
	if existing > 0 {
		return 0, bizerr.NewCode(CodeDomainAlreadyExists)
	}
	actor := s.actorID(ctx)
	id, err := shared.Model(ctx, shared.TableDomain).Data(domainInsertData{
		TenantId:          in.TenantId,
		Domain:            domain,
		IsPrimary:         in.IsPrimary,
		IsVerified:        false,
		VerificationToken: guid.S(),
		Status:            string(shared.DomainStatusActive),
		CreatedBy:         actor,
		UpdatedBy:         actor,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Delete soft-deletes a visible domain mapping after confirming it exists.
func (s *serviceImpl) Delete(ctx context.Context, id int64) error {
	visible, err := s.exists(ctx, id)
	if err != nil {
		return err
	}
	if !visible {
		return bizerr.NewCode(CodeDomainNotFound)
	}
	if _, err := shared.Model(ctx, shared.TableDomain).Where("id", id).Delete(); err != nil {
		return err
	}
	return nil
}

// SetVerified sets the verification flag of a visible domain mapping.
func (s *serviceImpl) SetVerified(ctx context.Context, id int64, verified bool) error {
	visible, err := s.exists(ctx, id)
	if err != nil {
		return err
	}
	if !visible {
		return bizerr.NewCode(CodeDomainNotFound)
	}
	if _, err := shared.Model(ctx, shared.TableDomain).Where("id", id).Data(domainVerifyData{
		IsVerified: verified,
		UpdatedBy:  s.actorID(ctx),
	}).Update(); err != nil {
		return err
	}
	return nil
}

// exists reports whether a non-deleted domain mapping with the id is present.
func (s *serviceImpl) exists(ctx context.Context, id int64) (bool, error) {
	count, err := shared.Model(ctx, shared.TableDomain).Where("id", id).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// actorID returns the current user id from the business context for audit fields.
func (s *serviceImpl) actorID(ctx context.Context) int64 {
	if s.bizCtxSvc == nil {
		return 0
	}
	return int64(s.bizCtxSvc.Current(ctx).UserID)
}

// normalizeDomain lowercases the domain host and strips any port suffix so
// stored mappings match resolver host normalization.
func normalizeDomain(domain string) string {
	hostname := strings.ToLower(strings.TrimSpace(domain))
	if colon := strings.LastIndex(hostname, ":"); colon >= 0 {
		hostname = hostname[:colon]
	}
	return hostname
}
