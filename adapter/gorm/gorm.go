// package gorm provides a GORM-based implementation of the Gordian store interfaces.
package gorm

import (
	"context"
	"fmt"

	"github.com/Robotech-Org/gordian"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// --- OrganizationStore Implementation ---

type OrganizationStore struct {
	DB *gorm.DB
}

func NewOrganizationStore(db *gorm.DB) *OrganizationStore {
	return &OrganizationStore{DB: db}
}

// Create satisfies the gordian.OrganizationStore interface.
func (s *OrganizationStore) Create(ctx context.Context, org *gordian.Organization) error {
	return s.DB.WithContext(ctx).Create(org).Error
}


func (s *OrganizationStore) Get(ctx context.Context, id uuid.UUID) (*gordian.Organization, error) {
	var org gordian.Organization
	if err := s.DB.WithContext(ctx).First(&org, id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}


// --- UserStore Implementation ---

type UserStore struct {
	DB *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{DB: db}
}

// Create satisfies the gordian.UserStore interface.
func (s *UserStore) Create(ctx context.Context, user *gordian.User) error {
	return s.DB.WithContext(ctx).Create(user).Error
}

func (s *UserStore) Get(ctx context.Context, id uuid.UUID) (*gordian.User, error) {
	var user gordian.User
	if err := s.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (s *UserStore) GetUserRole(ctx context.Context, userID uuid.UUID) (string, error) {
	var userRole string
	err := s.DB.WithContext(ctx).Model(&gordian.User{}).Where("id = ?", userID).Pluck("role", &userRole).Error
	if err != nil {
		return "", fmt.Errorf("failed to get user role: %w", err)
	}
	return userRole, nil
}

func (s *UserStore) FindByEmail(ctx context.Context, email string) (gordian.User, error) {
	var user gordian.User
	if err := s.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return gordian.User{}, fmt.Errorf("no user found")
		}
		return gordian.User{}, fmt.Errorf("failed to find user: %w", err)
	}
	return user, nil
}


// --- MembershipStore Implementation ---

type MembershipStore struct {
	DB *gorm.DB
}

func NewMembershipStore(db *gorm.DB) *MembershipStore {
	return &MembershipStore{DB: db}
}

// Create satisfies the gordian.MembershipStore interface.
func (s *MembershipStore) Create(ctx context.Context, membership *gordian.Membership) error {
	return s.DB.WithContext(ctx).Create(membership).Error
}


func (s *MembershipStore) GetMembers(ctx context.Context, orgID uuid.UUID) ([]*gordian.Membership, error) {
	var memberships []*gordian.Membership
	err := s.DB.WithContext(ctx).Where("organization_id = ?", orgID).Find(&memberships).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}
	return memberships, nil
}


// --- InviteStore Implementation ---

type InviteStore struct {
	DB *gorm.DB
}

func NewInviteStore(db *gorm.DB) *InviteStore {
	return &InviteStore{DB: db}
}

// Create satisfies the gordian.InviteStore interface.
func (s *InviteStore) Create(ctx context.Context, invite *gordian.Invite) error {
	return s.DB.WithContext(ctx).Create(invite).Error
}


func (s *InviteStore) Verify(ctx context.Context, token string) (bool, error) {
	query := `SELECT id FROM invites WHERE token = ?`
	var inviteID string
	err := s.DB.WithContext(ctx).Raw(query, token).Scan(&inviteID).Error
	if err != nil {
		return false, fmt.Errorf("failed to verify invitation: %w", err)
	}
	return true, nil
}

