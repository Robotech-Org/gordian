package services

import (
	"context"
	"github.com/Robotech-Org/gordian"
	"github.com/Robotech-Org/gordian/cmd/example-server/models"
	"gorm.io/gorm"
)

// UserService is Carol's high-level service.
type UserService struct {
	gordian *gordian.Service
	db      *gorm.DB
}

func NewUserService(gordianService *gordian.Service, db *gorm.DB) *UserService {
	return &UserService{
		gordian: gordianService,
		db:      db,
	}
}

// CreateUserWithStripe is Carol's custom method. It wraps Gordian's method.
func (s *UserService) CreateUserWithStripe(ctx context.Context, email, name string) (*models.User, error) {
	// 1. Use Gordian to create the base user.
	// This handles the core logic of creating the user record.
	baseUser, err := s.gordian.CreateUser(ctx, email, name)
	if err != nil {
		return nil, err
	}

	// 2. Retrieve the full user model (including the Stripe ID created by the hook).
	// We need to fetch the record again to get the value set by the AfterCreate hook.
	var fullUser models.User
	if err := s.db.First(&fullUser, baseUser.ID).Error; err != nil {
		return nil, err
	}

	return &fullUser, nil
}
