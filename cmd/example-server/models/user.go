package models

import (
	"github.com/Robotech-Org/gordian"
	"gorm.io/gorm"
)

// User is Carol's application-specific user model.
// It contains all of Gordian's fields, plus custom ones.
type User struct {
	gordian.User
	StripeCustomerID string `gorm:"column:stripe_customer_id;unique"`
	PhoneNumber      string `gorm:"column:phone_number;unique"`
}

// This is an advanced GORM feature: Hooks.
// This hook ensures that every time a new User is created,
// we also create a corresponding Stripe customer.
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	// In a real app, this would be a call to the Stripe API.
	// stripeCustomer, err := stripe.Customers.New(...)
	// We simulate this by creating a fake ID.
	fakeStripeID := "cus_" + u.ID.String()

	// Update the user record with the new Stripe ID within the same transaction.
	return tx.Model(u).Update("stripe_customer_id", fakeStripeID).Error
}
