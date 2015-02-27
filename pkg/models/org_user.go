package models

import (
	"errors"
	"time"
)

// Typed errors
var (
	ErrInvalidRoleType = errors.New("Invalid role type")
	ErrLastOrgAdmin    = errors.New("Cannot remove last organization admin")
)

type RoleType string

const (
	ROLE_VIEWER         RoleType = "Viewer"
	ROLE_EDITOR         RoleType = "Editor"
	ROLE_ADMIN          RoleType = "Admin"
	ROLE_RAINTANK_ADMIN RoleType = "RaintankAdmin"
)

func (r RoleType) IsValid() bool {
	return r == ROLE_VIEWER || r == ROLE_ADMIN || r == ROLE_EDITOR
}

type OrgUser struct {
	OrgId   int64
	UserId  int64
	Role    RoleType
	Created time.Time
	Updated time.Time
}

// ---------------------
// COMMANDS

type RemoveOrgUserCommand struct {
	UserId int64
	OrgId  int64
}

type AddOrgUserCommand struct {
	LoginOrEmail string   `json:"loginOrEmail" binding:"Required"`
	Role         RoleType `json:"role" binding:"Required"`

	OrgId  int64 `json:"-"`
	UserId int64 `json:"-"`
}

// ----------------------
// QUERIES

type GetOrgUsersQuery struct {
	OrgId  int64
	Result []*OrgUserDTO
}

// ----------------------
// Projections and DTOs

type OrgUserDTO struct {
	OrgId  int64  `json:"orgId"`
	UserId int64  `json:"userId"`
	Email  string `json:"email"`
	Login  string `json:"login"`
	Role   string `json:"role"`
}
