package gordian

import "context"

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

func (s *GordianService) CreateOrganization(ctx context.Context, org *Organization) error {
	return s.OrganizationStore.Create(ctx, org)
}

func (s *GordianService) CreateUser(ctx context.Context, user *User) error {
	return s.UserStore.Create(ctx, user)
}

func (s *GordianService) CreateMembership(ctx context.Context, membership *Membership) error {
	return s.MembershipStore.Create(ctx, membership)
}

func (s *GordianService) CreateInvitation(ctx context.Context, invitation *Invite) error {
	return s.InvitationStore.Create(ctx, invitation)
}
