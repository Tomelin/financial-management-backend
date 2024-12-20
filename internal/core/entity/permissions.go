package entity

import (
	"errors"
)

type PermissionLevel string

const (
	PermissionOwner       PermissionLevel = "owner"
	PermissionAdmin       PermissionLevel = "admin"
	PermissionEdit        PermissionLevel = "edit"
	PermissionView        PermissionLevel = "view"
	PermissionSystemAdmin PermissionLevel = "system-admin"
	PermissionSystemView  PermissionLevel = "system-view"
)

func (p *PermissionLevel) Validate() error {
	switch *p {
	case PermissionOwner, PermissionAdmin, PermissionEdit, PermissionView:
		return nil
	default:
		return errors.New("invalid permission")
	}
}

type Module string

const (
	ModuleUser    Module = "user"
	ModuleTenant  Module = "tenant"
	ModuleProduct Module = "product"
	ModuleOrder   Module = "order"
	ModuleInvoice Module = "invoice"
)

func (m *Module) Validate() error {
	switch *m {
	case ModuleUser, ModuleTenant, ModuleProduct, ModuleOrder, ModuleInvoice:
		return nil
	default:
		return errors.New("invalid module")
	}
}

type TenantPermission struct {
	Module    Module          `json:"module" firestore:"module"`
	Level     PermissionLevel `json:"level" firestore:"level"`
	TennantID string          `json:"tenant_id" firestore:"tenant_id"`
	Owner     bool            `json:"owner_id" firestore:"owner_id"`
}

type WalletPermission struct {
	Level     PermissionLevel `json:"level" firestore:"level"`
	TennantID string          `json:"tenant_id" firestore:"tenant_id"`
	WalletID  string          `json:"wallet_id" firestore:"wallet_id"`
	Owner     bool            `json:"owner_id" firestore:"owner_id"`
}
