package gordian

import (
	"context"

	"github.com/google/uuid"
)

// Defines contract for storing organizations.
type OrganizationStore interface {
	Create(ctx context.Context, org *Organization) error
	Get(ctx context.Context, id uuid.UUID) (*Organization, error)
}

// Defines contract for storing users.
type UserStore interface {
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserRole(ctx context.Context, userID uuid.UUID) (string, error)
	FindByEmail(ctx context.Context, email string) (User, error)
}

// Defines contract for storing memberships.
type MembershipStore interface {
	Create(ctx context.Context, membership *Membership) error
	GetMembers(ctx context.Context, orgID uuid.UUID) ([]*Membership, error)
	GetMembership(ctx context.Context, userID uuid.UUID, orgID uuid.UUID) (Membership, error)
}

// Defines contract for storing invitations.
type InvitationStore interface {
	Create(ctx context.Context, invite *Invite) error
	Verify(ctx context.Context, token string) (bool, error)
}

// Defines the contract for sending emails.
type Emailer interface {
	SendInvitation(ctx context.Context, invite *Invite) error
}
