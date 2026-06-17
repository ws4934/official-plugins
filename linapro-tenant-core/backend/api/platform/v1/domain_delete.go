// This file declares platform tenant domain delete DTOs for the linapro-tenant-core source plugin.

package v1

import "github.com/gogf/gf/v2/frame/g"

// DomainDeleteReq defines the request for deleting a tenant domain mapping.
type DomainDeleteReq struct {
	g.Meta `path:"/platform/domains/{id}" method:"delete" tags:"Platform Tenant Domains" summary:"Delete tenant domain" dc:"Soft-delete a tenant domain mapping by ID." permission:"system:tenant:domain:remove"`
	Id     int64 `json:"id" v:"required|min:1#gf.gvalid.rule.required|Domain ID must be a positive integer" dc:"Domain mapping ID" eg:"1"`
}

// DomainDeleteRes defines the tenant domain delete response.
type DomainDeleteRes struct{}
