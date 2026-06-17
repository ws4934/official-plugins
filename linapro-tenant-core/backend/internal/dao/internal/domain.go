// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DomainDao is the data access object for the table plugin_linapro_tenant_core_domain.
type DomainDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DomainColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DomainColumns defines and stores column names for the table plugin_linapro_tenant_core_domain.
type DomainColumns struct {
	Id                string //
	TenantId          string //
	Domain            string //
	IsPrimary         string //
	IsVerified        string //
	VerificationToken string //
	Status            string //
	CreatedBy         string //
	UpdatedBy         string //
	CreatedAt         string //
	UpdatedAt         string //
	DeletedAt         string //
}

// domainColumns holds the columns for the table plugin_linapro_tenant_core_domain.
var domainColumns = DomainColumns{
	Id:                "id",
	TenantId:          "tenant_id",
	Domain:            "domain",
	IsPrimary:         "is_primary",
	IsVerified:        "is_verified",
	VerificationToken: "verification_token",
	Status:            "status",
	CreatedBy:         "created_by",
	UpdatedBy:         "updated_by",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
	DeletedAt:         "deleted_at",
}

// NewDomainDao creates and returns a new DAO object for table data access.
func NewDomainDao(handlers ...gdb.ModelHandler) *DomainDao {
	return &DomainDao{
		group:    "default",
		table:    "plugin_linapro_tenant_core_domain",
		columns:  domainColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DomainDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DomainDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DomainDao) Columns() DomainColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DomainDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DomainDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *DomainDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
