package gordian

import "context"

// Defines contract for storing organizations.
type OrganizationStore interface {
	Create(ctx context.Context, org *Organization) error
}

// Defines contract for storing users.
type UserStore interface {
	Create(ctx context.Context, user *User) error
}

// Defines contract for storing memberships.
type MembershipStore interface {
	Create(ctx context.Context, membership *Membership) error
}

// Defines contract for storing invitations.
type InvitationStore interface {
	Create(ctx context.Context, invite *Invite) error
}

// Defines the contract for sending emails.
type Emailer interface {
	SendInvitation(ctx context.Context, invite *Invite) error
}
