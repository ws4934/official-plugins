// This file declares platform tenant domain verification DTOs for the linapro-tenant-core source plugin.

package v1

import "github.com/gogf/gf/v2/frame/g"

// DomainVerifyReq defines the request for setting tenant domain verification.
type DomainVerifyReq struct {
	g.Meta   `path:"/platform/domains/{id}/verification" method:"put" tags:"Platform Tenant Domains" summary:"Set tenant domain verification" dc:"Set the verification flag of a tenant domain mapping. Only verified, active domains resolve tenants." permission:"system:tenant:domain:verify"`
	Id       int64 `json:"id" v:"required|min:1#gf.gvalid.rule.required|Domain ID must be a positive integer" dc:"Domain mapping ID" eg:"1"`
	Verified bool  `json:"verified" dc:"Target verification state" eg:"true"`
}

// DomainVerifyRes defines the tenant domain verification response.
type DomainVerifyRes struct{}
