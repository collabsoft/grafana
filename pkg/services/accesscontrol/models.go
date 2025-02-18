package accesscontrol

import (
	"encoding/json"
	"strings"
	"time"
)

// RoleRegistration stores a role and its assignments to built-in roles
// (Viewer, Editor, Admin, Grafana Admin)
type RoleRegistration struct {
	Role   RoleDTO
	Grants []string
}

// Role is the model for Role in RBAC.
type Role struct {
	ID          int64  `json:"-" xorm:"pk autoincr 'id'"`
	OrgID       int64  `json:"-" xorm:"org_id"`
	Version     int64  `json:"version"`
	UID         string `xorm:"uid" json:"uid"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`

	Updated time.Time `json:"updated"`
	Created time.Time `json:"created"`
}

func (r Role) Global() bool {
	return r.OrgID == GlobalOrgID
}

func (r Role) IsFixed() bool {
	return strings.HasPrefix(r.Name, FixedRolePrefix)
}

func (r Role) GetDisplayName() string {
	if r.IsFixed() && r.DisplayName == "" {
		r.DisplayName = fallbackDisplayName(r.Name)
	}
	return r.DisplayName
}

func (r Role) MarshalJSON() ([]byte, error) {
	type Alias Role

	r.DisplayName = r.GetDisplayName()
	return json.Marshal(&struct {
		Alias
		Global bool `json:"global" xorm:"-"`
	}{
		Alias:  (Alias)(r),
		Global: r.Global(),
	})
}

type RoleDTO struct {
	Version     int64        `json:"version"`
	UID         string       `xorm:"uid" json:"uid"`
	Name        string       `json:"name"`
	DisplayName string       `json:"displayName"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions,omitempty"`

	ID    int64 `json:"-" xorm:"pk autoincr 'id'"`
	OrgID int64 `json:"-" xorm:"org_id"`

	Updated time.Time `json:"updated"`
	Created time.Time `json:"created"`
}

func (r RoleDTO) Role() Role {
	return Role{
		ID:          r.ID,
		OrgID:       r.OrgID,
		UID:         r.UID,
		Name:        r.Name,
		DisplayName: r.DisplayName,
		Description: r.Description,
		Updated:     r.Updated,
		Created:     r.Created,
	}
}

func (r RoleDTO) Global() bool {
	return r.OrgID == GlobalOrgID
}

func (r RoleDTO) IsFixed() bool {
	return strings.HasPrefix(r.Name, FixedRolePrefix)
}

func (r RoleDTO) GetDisplayName() string {
	if r.IsFixed() && r.DisplayName == "" {
		r.DisplayName = fallbackDisplayName(r.Name)
	}
	return r.DisplayName
}

func (r RoleDTO) MarshalJSON() ([]byte, error) {
	type Alias RoleDTO

	r.DisplayName = r.GetDisplayName()
	return json.Marshal(&struct {
		Alias
		Global bool `json:"global" xorm:"-"`
	}{
		Alias:  (Alias)(r),
		Global: r.Global(),
	})
}

// fallbackDisplayName provides a fallback name for role
// that can be displayed in the ui for better readability
// example: currently this would give:
// fixed:datasources:name -> datasources name
// datasources:admin      -> datasources admin
func fallbackDisplayName(rName string) string {
	// removing prefix for fixed roles
	rNameWithoutPrefix := strings.Replace(rName, FixedRolePrefix, "", 1)
	return strings.TrimSpace(strings.Replace(rNameWithoutPrefix, ":", " ", -1))
}

// Permission is the model for access control permissions.
type Permission struct {
	ID     int64  `json:"-" xorm:"pk autoincr 'id'"`
	RoleID int64  `json:"-" xorm:"role_id"`
	Action string `json:"action"`
	Scope  string `json:"scope"`

	Updated time.Time `json:"updated"`
	Created time.Time `json:"created"`
}

func (p Permission) OSSPermission() Permission {
	return Permission{
		Action: p.Action,
		Scope:  p.Scope,
	}
}

// ScopeParams holds the parameters used to fill in scope templates
type ScopeParams struct {
	OrgID     int64
	URLParams map[string]string
}

const (
	GlobalOrgID = 0
	// Permission actions

	// Users actions
	ActionUsersRead     = "users:read"
	ActionUsersWrite    = "users:write"
	ActionUsersTeamRead = "users.teams:read"
	// We can ignore gosec G101 since this does not contain any credentials.
	// nolint:gosec
	ActionUsersAuthTokenList = "users.authtoken:list"
	// We can ignore gosec G101 since this does not contain any credentials.
	// nolint:gosec
	ActionUsersAuthTokenUpdate = "users.authtoken:update"
	// We can ignore gosec G101 since this does not contain any credentials.
	// nolint:gosec
	ActionUsersPasswordUpdate    = "users.password:update"
	ActionUsersDelete            = "users:delete"
	ActionUsersCreate            = "users:create"
	ActionUsersEnable            = "users:enable"
	ActionUsersDisable           = "users:disable"
	ActionUsersPermissionsUpdate = "users.permissions:update"
	ActionUsersLogout            = "users:logout"
	ActionUsersQuotasList        = "users.quotas:list"
	ActionUsersQuotasUpdate      = "users.quotas:update"

	// Org actions
	ActionOrgUsersRead       = "org.users:read"
	ActionOrgUsersAdd        = "org.users:add"
	ActionOrgUsersRemove     = "org.users:remove"
	ActionOrgUsersRoleUpdate = "org.users.role:update"

	// LDAP actions
	ActionLDAPUsersRead    = "ldap.user:read"
	ActionLDAPUsersSync    = "ldap.user:sync"
	ActionLDAPStatusRead   = "ldap.status:read"
	ActionLDAPConfigReload = "ldap.config:reload"

	// Server actions
	ActionServerStatsRead = "server.stats:read"

	// Settings actions
	ActionSettingsRead = "settings:read"

	// Datasources actions
	ActionDatasourcesExplore = "datasources:explore"

	// Plugin actions
	ActionPluginsManage = "plugins:manage"

	// Global Scopes
	ScopeGlobalUsersAll = "global:users:*"

	// Users scope
	ScopeUsersAll = "users:*"

	// Settings scope
	ScopeSettingsAll = "settings:*"

	// Licensing related actions
	ActionLicensingRead        = "licensing:read"
	ActionLicensingUpdate      = "licensing:update"
	ActionLicensingDelete      = "licensing:delete"
	ActionLicensingReportsRead = "licensing.reports:read"
)

const RoleGrafanaAdmin = "Grafana Admin"

const FixedRolePrefix = "fixed:"

// LicensingPageReaderAccess defines permissions that grant access to the licensing and stats page
var LicensingPageReaderAccess = EvalAny(
	EvalPermission(ActionLicensingRead),
	EvalPermission(ActionServerStatsRead),
)
