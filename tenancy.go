package gordian

import (
	"time"

	"github.com/google/uuid" // A good choice for IDs
)

// User represents a global user in the system.
// A user is NOT a tenant. They are an identity who can join tenants.
type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	CreatedAt time.Time
}

// Organization is the Tenant. It is the top-level container for users and resources.
type Organization struct {
	ID        uuid.UUID
	Name      string
	OwnerID   uuid.UUID // The user who created and owns the organization
	CreatedAt time.Time
}

// Membership is the junction entity that links a User to an Organization.
// This is the "User-under-Organization" block. It's the most important struct.
type Membership struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID // Foreign Key to Organization
	UserID         uuid.UUID // Foreign Key to User
	Role           string    // e.g., "owner", "admin", "member"
	JoinedAt       time.Time
}

// Invite represents a pending invitation for a user to join an organization.
type Invite struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID // The organization the user is being invited to
	InviterID      uuid.UUID // The user who sent the invite
	InviteeEmail   string    // The email of the person being invited
	Role           string    // The role they will have when they accept
	Token          string    // A unique, secret token for the invite link
	ExpiresAt      time.Time
	CreatedAt      time.Time
}
