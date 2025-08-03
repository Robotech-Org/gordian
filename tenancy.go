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

func NewUser(email string, name string) *User {
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
	}
}

// Organization is the Tenant. It is the top-level container for users and resources.
type Organization struct {
	ID        uuid.UUID
	Name      string
	OwnerID   uuid.UUID // The user who created and owns the organization
	CreatedAt time.Time
}

func NewOrganization(ownerID uuid.UUID, name string) *Organization {
	return &Organization{
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
	}
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

func NewMembership(userID, organizationID uuid.UUID, role string) *Membership {
	return &Membership{
		UserID:         userID,
		OrganizationID: organizationID,
		Role:           role,
		JoinedAt:       time.Now(),
	}
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

func NewInvite(organizationID, inviterID uuid.UUID, inviteeEmail, role, token string) *Invite {
	createdAt := time.Now()
	expiresAt := createdAt.Add(time.Hour * 24)
	return &Invite{
		OrganizationID: organizationID,
		InviterID:      inviterID,
		InviteeEmail:   inviteeEmail,
		Role:           role,
		Token:          token,
		CreatedAt:      createdAt,
		ExpiresAt:      expiresAt,
	}
}
