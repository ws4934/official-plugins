// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Domain is the golang structure of table plugin_linapro_tenant_core_domain for DAO operations like Where/Data.
type Domain struct {
	g.Meta            `orm:"table:plugin_linapro_tenant_core_domain, do:true"`
	Id                any        //
	TenantId          any        //
	Domain            any        //
	IsPrimary         any        //
	IsVerified        any        //
	VerificationToken any        //
	Status            any        //
	CreatedBy         any        //
	UpdatedBy         any        //
	CreatedAt         *time.Time //
	UpdatedAt         *time.Time //
	DeletedAt         *time.Time //
}
