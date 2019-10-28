package rbac

import (
	"errors"
	"sync"
)

// Role is an interface.
// You should implement this interface for your own role structures.
type Role interface {
	ID() string
	Permit(Permission) bool
	Permissions() []Permission
	Revoke(p Permission) error
	Assign(p Permission) error
}

var (
	// ErrRoleNotExist occurred if a role cann't be found
	ErrPermNotExist = errors.New("Permission does not exist")
)

// Roles is a map
type Roles map[string]Role

// NewStdRole is the default role factory function.
// It matches the declaration to RoleFactoryFunc.
func NewStdRole(id string) *StdRole {
	role := &StdRole{
		IDStr:       id,
		permissions: make(Permissions),
	}
	return role
}

// StdRole is the default role implement.
// You can combine this struct into your own Role implement.
type StdRole struct {
	sync.RWMutex
	// IDStr is the identity of role
	IDStr       string `json:"id"`
	permissions Permissions
}

// ID returns the role's identity name.
func (role *StdRole) ID() string {
	return role.IDStr
}

// Assign a permission to the role.
func (role *StdRole) Assign(p Permission) error {
	role.Lock()
	role.permissions[p.ID()] = p
	role.Unlock()
	return nil
}

// Permit returns true if the role has specific permission.
func (role *StdRole) Permit(p Permission) (rslt bool) {
	if p == nil {
		return false
	}

	role.RLock()
	for _, rp := range role.permissions {
		if rp.Match(p) {
			rslt = true
			break
		}
	}
	role.RUnlock()
	return
}

// Revoke the specific permission.
func (role *StdRole) Revoke(p Permission) error {
	if role.Permit(p) {
		role.Lock()
		delete(role.permissions, p.ID())
		role.Unlock()
		return nil
	} else {
		role.Unlock()
		return ErrPermNotExist
	}
}

// Permissions returns all permissions into a slice.
func (role *StdRole) Permissions() []Permission {
	role.RLock()
	result := make([]Permission, 0, len(role.permissions))
	for _, p := range role.permissions {
		result = append(result, p)
	}
	role.RUnlock()
	return result
}
