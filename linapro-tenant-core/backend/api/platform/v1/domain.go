// This file declares platform tenant domain list DTOs for the linapro-tenant-core source plugin.

package v1

import "github.com/gogf/gf/v2/frame/g"

// DomainListReq defines the request for listing tenant domain mappings.
type DomainListReq struct {
	g.Meta   `path:"/platform/domains" method:"get" tags:"Platform Tenant Domains" summary:"Get tenant domain list" dc:"Query tenant domain mappings by page with optional tenant, domain, and status filters. Platform governance data gated by platform permission." permission:"system:tenant:domain:list"`
	PageNum  int    `json:"pageNum" d:"1" v:"min:1" dc:"Page number" eg:"1"`
	PageSize int    `json:"pageSize" d:"10" v:"min:1|max:100" dc:"Number of items per page" eg:"10"`
	TenantId int64  `json:"tenantId" dc:"Filter by owning tenant ID. Zero means no tenant filter." eg:"1"`
	Domain   string `json:"domain" dc:"Filter by domain host substring" eg:"acme"`
	Status   string `json:"status" dc:"Filter by domain status. One of active, disabled. Empty means no status filter." eg:"active"`
}

// DomainListRes defines the tenant domain list response.
type DomainListRes struct {
	List  []*DomainItem `json:"list" dc:"Tenant domain list" eg:"[]"`
	Total int           `json:"total" dc:"Total domain count" eg:"3"`
}

// DomainItem is the tenant domain list item projection.
type DomainItem struct {
	Id         int64  `json:"id" dc:"Domain mapping ID" eg:"1"`
	TenantId   int64  `json:"tenantId" dc:"Owning tenant ID" eg:"1"`
	Domain     string `json:"domain" dc:"Mapped domain host, stored lowercase" eg:"shop.acme.com"`
	IsPrimary  bool   `json:"isPrimary" dc:"Whether this is the tenant primary domain" eg:"true"`
	IsVerified bool   `json:"isVerified" dc:"Whether the domain is verified and usable for resolution" eg:"true"`
	Status     string `json:"status" dc:"Domain status. One of active, disabled." eg:"active"`
	CreatedAt  *int64 `json:"createdAt" dc:"Creation time. Unix timestamp in milliseconds." eg:"1717000000000"`
}
