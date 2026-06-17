// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// Domain is the golang structure for table domain.
type Domain struct {
	Id                int64      `json:"id"                orm:"id"                 description:""`
	TenantId          int64      `json:"tenantId"          orm:"tenant_id"          description:""`
	Domain            string     `json:"domain"            orm:"domain"             description:""`
	IsPrimary         bool       `json:"isPrimary"         orm:"is_primary"         description:""`
	IsVerified        bool       `json:"isVerified"        orm:"is_verified"        description:""`
	VerificationToken string     `json:"verificationToken" orm:"verification_token" description:""`
	Status            string     `json:"status"            orm:"status"             description:""`
	CreatedBy         int64      `json:"createdBy"         orm:"created_by"         description:""`
	UpdatedBy         int64      `json:"updatedBy"         orm:"updated_by"         description:""`
	CreatedAt         *time.Time `json:"createdAt"         orm:"created_at"         description:""`
	UpdatedAt         *time.Time `json:"updatedAt"         orm:"updated_at"         description:""`
	DeletedAt         *time.Time `json:"deletedAt"         orm:"deleted_at"         description:""`
}
