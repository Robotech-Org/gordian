// package gorm provides a GORM-based implementation of the Gordian store interfaces.
package gorm

import (
	"context"

	"github.com/Robotech-Org/gordian"
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
