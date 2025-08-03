package main

import (
	"context"
	"log"
	"os"

	"github.com/Robotech-Org/gordian"
	gormadapter "github.com/Robotech-Org/gordian/adapter/gorm" // Import your new adapter
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type NoOpEmailer struct{}

func (e *NoOpEmailer) SendInvitation(ctx context.Context, invite *gordian.Invite) error {
	log.Printf("SIMULATING SENDING EMAIL: Would send invite with token %s to %s", invite.Token, invite.InviteeEmail)
	return nil
}

func main() {
	// --- Step 1: Connect to the database (Carol's responsibility) ---
	// In a real app, this DSN comes from env vars or a config file.
	dsn := os.Getenv("DATABASE_URL")
	log.Printf("Connecting to database with dsn: %s", dsn)
	if dsn == "" {
		dsn = "host=localhost user=postgres password=password dbname=gordian_test port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// This is a great place to auto-migrate the schema for testing purposes.
	// In a real app, Carol would use a separate migration tool.
	db.AutoMigrate(&gordian.User{}, &gordian.Organization{}, &gordian.Membership{}, &gordian.Invite{})
	log.Println("Database migration complete.")

	// --- Step 2: Initialize the adapter and inject it into the service (The "Wiring") ---
	orgStore := gormadapter.NewOrganizationStore(db)
	userStore := gormadapter.NewUserStore(db)
	memStore := gormadapter.NewMembershipStore(db)
	invStore := gormadapter.NewInviteStore(db)
	dummyEmailer := &NoOpEmailer{}

	// For now, we don't have an emailer, so we can pass nil or a mock.
	// We'll need to adjust the New() function to allow a nil emailer for now.
	// Or create a dummy emailer.

	gordianService := gordian.New(orgStore, userStore, memStore, invStore, dummyEmailer)
	log.Println("Gordian service initialized.")

	// --- Step 3: Use the service to perform a real operation (The Test) ---
	log.Println("Attempting to create a new user and organization...")

	// Create a user first
	user, err := gordianService.CreateUser(context.Background(), "test.user@example.com", "Test User")
	if err != nil {
		log.Fatalf("ERROR: Failed to create user: %v", err)
	}
	log.Printf("SUCCESS: Created user with ID: %s", user.ID)

	// Now create an organization with that user as the owner
	org, err := gordianService.CreateOrganization(context.Background(), "My First Test Org", user.ID)
	if err != nil {
		log.Fatalf("ERROR: Failed to create organization: %v", err)
	}
	log.Printf("SUCCESS: Created organization '%s' with ID: %s", org.Name, org.ID)
	log.Println("Integration test successful!")
}
