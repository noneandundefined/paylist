package constants

import "time"

const (
	FreeMaxSharedMembers = 2
	ShareInviteTTL       = 7 * 24 * time.Hour

	MemberRoleOwner    = "owner"
	MemberRoleMember   = "member"
	MemberRoleObserver = "observer"
)
