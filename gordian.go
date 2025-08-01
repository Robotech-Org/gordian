package gordian

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type OrganizationStore interface {
	Create(ctx context.Context, org *Organization) error
}

type UserStore interface {
	Create(ctx context.Context, user *User) error
}

type MembershipStore interface {
	Create(ctx context.Context, membership *Membership) error
}

type InvitationStore interface {
	Create(ctx context.Context, invitation *Invite) error
}

type GordianService struct {
	OrganizationStore OrganizationStore
	UserStore         UserStore
	MembershipStore   MembershipStore
	InvitationStore   InvitationStore
}

func NewGordianService(
	orgStore OrganizationStore,
	userStore UserStore,
	membershipStore MembershipStore,
	invitationStore InvitationStore,
) *GordianService {
	return &GordianService{
		OrganizationStore: orgStore,
		UserStore:         userStore,
		MembershipStore:   membershipStore,
		InvitationStore:   invitationStore,
	}
}

func (s *GordianService) CreateOrganization(ctx context.Context, name string, ownerID uuid.UUID) (*Organization, error) {
	// 1. Validation Step
	if len(name) < 3 {
		return nil, fmt.Errorf("invalid organization name")
	}

	// 2. Create the organization
	org := NewOrganization(ownerID, name)
	if err := s.OrganizationStore.Create(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// 3. Create the owner membership[the user that created the organization is the owner / has the role of owner]
	ownerMembership := NewMembership(ownerID, org.ID, "owner")
	if err := s.MembershipStore.Create(ctx, ownerMembership); err != nil {
		return nil, fmt.Errorf("failed to create membership: %w", err)
	}

	return org, nil
}

func (s *GordianService) CreateUser(ctx context.Context, user *User) error {
	return s.UserStore.Create(ctx, user)
}

func (s *GordianService) CreateMembership(ctx context.Context, membership *Membership) error {
	return s.MembershipStore.Create(ctx, membership)
}

func (s *GordianService) CreateInvitation(ctx context.Context, organizationID, inviterID uuid.UUID, inviteeEmail, role string) error {
	// 1. Validate input
	if inviteeEmail == "" {
		return errors.New("invitee email cannot be empty")
	}

	//2. Create Token
	token := uuid.New().String()

	// 3. Create the invitation
	invitation := NewInvite(organizationID, inviterID, inviteeEmail, role, token)
	if err := s.InvitationStore.Create(ctx, invitation); err != nil {
		return fmt.Errorf("failed to create invitation: %w", err)
	}

	// I will have to research on how to securily send email to invitee with token
	// err := s.EmailService.SendInvitationEmail(ctx, inviteeEmail, token)
	// if err != nil {
	// 	return fmt.Errorf("failed to send invitation email: %w", err)
	// }

	return nil
}
