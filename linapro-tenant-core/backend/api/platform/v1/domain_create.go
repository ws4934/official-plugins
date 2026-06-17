// This file declares platform tenant domain create DTOs for the linapro-tenant-core source plugin.

package v1

import "github.com/gogf/gf/v2/frame/g"

// DomainCreateReq defines the request for creating a tenant domain mapping.
type DomainCreateReq struct {
	g.Meta    `path:"/platform/domains" method:"post" tags:"Platform Tenant Domains" summary:"Create tenant domain" dc:"Map a domain host to a tenant. The domain is normalized to lowercase and must be globally unique. New mappings start unverified." permission:"system:tenant:domain:add"`
	TenantId  int64  `json:"tenantId" v:"required|min:1#gf.gvalid.rule.required|Tenant ID must be a positive integer" dc:"Owning tenant ID" eg:"1"`
	Domain    string `json:"domain" v:"required|max-length:255#gf.gvalid.rule.required|Domain must have at most 255 characters" dc:"Domain host to map. Stored lowercase without port." eg:"shop.acme.com"`
	IsPrimary bool   `json:"isPrimary" dc:"Whether to mark this mapping as the tenant primary domain" eg:"false"`
}

// DomainCreateRes defines the tenant domain create response.
type DomainCreateRes struct {
	Id int64 `json:"id" dc:"Domain mapping ID" eg:"1"`
}
