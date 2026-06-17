// This file defines tenant domain business error codes and runtime metadata.

package domain

import (
	"github.com/gogf/gf/v2/errors/gcode"

	"lina-core/pkg/bizerr"
)

var (
	// CodeDomainInvalid reports that the supplied domain host or tenant is invalid.
	CodeDomainInvalid = bizerr.MustDefine(
		"MULTI_TENANT_DOMAIN_INVALID",
		"Domain host or tenant is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeDomainAlreadyExists reports that the domain is already mapped to a tenant.
	CodeDomainAlreadyExists = bizerr.MustDefine(
		"MULTI_TENANT_DOMAIN_ALREADY_EXISTS",
		"Domain is already mapped to a tenant",
		gcode.CodeInvalidParameter,
	)
	// CodeDomainNotFound reports that the tenant domain mapping does not exist.
	CodeDomainNotFound = bizerr.MustDefine(
		"MULTI_TENANT_DOMAIN_NOT_FOUND",
		"Tenant domain mapping does not exist",
		gcode.CodeNotFound,
	)
)
