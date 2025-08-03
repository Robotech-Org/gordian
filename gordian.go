package gordian

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	orgStore  OrganizationStore
	userStore UserStore
	memStore  MembershipStore
	invStore  InvitationStore
	emailer   Emailer
}

func New(
	orgStore OrganizationStore,
	userStore UserStore,
	membershipStore MembershipStore,
	invitationStore InvitationStore,
	emailer Emailer,

) *Service {
	return &Service{
		orgStore:  orgStore,
		userStore: userStore,
		memStore:  membershipStore,
		invStore:  invitationStore,
		emailer:   emailer,
	}
}

func (s *Service) CreateOrganization(ctx context.Context, name string, ownerID uuid.UUID) (*Organization, error) {
	// 1. Validation Step
	if len(name) < 3 {
		return nil, fmt.Errorf("invalid organization name")
	}

	// 2. Create the organization
	org := NewOrganization(ownerID, name)
	if err := s.orgStore.Create(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// 3. Create the owner membership[the user that created the organization is the owner / has the role of owner]
	ownerMembership := NewMembership(ownerID, org.ID, "owner")
	if err := s.memStore.Create(ctx, ownerMembership); err != nil {
		return nil, fmt.Errorf("failed to create membership: %w", err)
	}

	return org, nil
}

func (s *Service) CreateUser(ctx context.Context, email, name string) (*User, error) {
	user := NewUser(email, name)
	if err := s.userStore.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *Service) CreateMembership(ctx context.Context, userID, orgID uuid.UUID, role string) (*Membership, error) {
	membership := NewMembership(userID, orgID, role)
	if err := s.memStore.Create(ctx, membership); err != nil {
		return nil, fmt.Errorf("failed to create membership: %w", err)
	}
	return membership, nil
}

func (s *Service) CreateInvitation(ctx context.Context, organizationID, inviterID uuid.UUID, inviteeEmail, role string) (*Invite, error) {
	// 1. Validate input
	if inviteeEmail == "" {
		return nil, errors.New("invitee email cannot be empty")
	}

	//2. Create Token
	token := uuid.New().String()

	// 3. Create the invitation
	invitation := NewInvite(organizationID, inviterID, inviteeEmail, role, token)
	if err := s.invStore.Create(ctx, invitation); err != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	// I will have to research on how to securily send email to invitee with token
	err := s.emailer.SendInvitation(ctx, invitation)
	if err != nil {
		return nil, fmt.Errorf("failed to send invitation email: %w", err)
	}

	return invitation, nil
}
