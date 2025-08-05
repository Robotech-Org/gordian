package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Robotech-Org/gordian"
	gormadapter "github.com/Robotech-Org/gordian/adapter/gorm"
	"github.com/Robotech-Org/gordian/cmd/example-server/mailtrap"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	godotenv.Load()
}

func main() {
	// --- Step 1: Connect to the database (Carol's responsibility) ---
	dsn := os.Getenv("DATABASE_URL")
	log.Printf("Connecting to database with dsn: %s", dsn)
	if dsn == "" {
		dsn = "host=localhost user=postgres password=password dbname=gordian_test port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	mailtrapCfg := mailtrap.Config{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     os.Getenv("EMAIL_PORT"),
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		FromAddr: "noreply@diagramly.com",
	}
	emailer := mailtrap.NewEmailer(mailtrapCfg)

	db.AutoMigrate(&gordian.User{}, &gordian.Organization{}, &gordian.Membership{}, &gordian.Invite{})
	log.Println("Database migration complete.")

	// --- Step 2: Initialize the adapter and inject it into the service (The "Wiring") ---
	orgStore := gormadapter.NewOrganizationStore(db)
	userStore := gormadapter.NewUserStore(db)
	memStore := gormadapter.NewMembershipStore(db)
	invStore := gormadapter.NewInviteStore(db)

	gordianService := gordian.New(orgStore, userStore, memStore, invStore, emailer)
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

	// Get the organization by ID
	retrievedOrg, err := gordianService.GetOrganization(context.Background(), org.ID)
	if err != nil {
		log.Fatalf("ERROR: Failed to get organization: %v", err)
	}
	log.Printf("SUCCESS: Retrieved organization '%s' with ID: %s", retrievedOrg.Name, retrievedOrg.ID)

	retrievedUser, err := gordianService.GetUser(context.Background(), user.ID)
	if err != nil {
		log.Fatalf("ERROR: Failed to get user: %v", err)
	}
	log.Printf("SUCCESS: Retrieved user '%s' with ID: %s", retrievedUser.Email, retrievedUser.ID)

	log.Println("--- Testing Invitation Flow ---")
	invite, err := gordianService.CreateInvitation(context.Background(), org.ID, user.ID, "new.colleague@example.com", "editor")
	if err != nil {
		log.Fatalf("ERROR: Failed to create and send invitation: %v", err)
	}
	log.Printf("SUCCESS: Invitation created and email sent (trapped by Mailtrap). Invite Token: %s", invite.Token)
	// Create an endpoint to get the invitation token, for example the user may get a link to accept the invitation:
	// http://localhost:8080/accept-invite?token=5361fd0f-8394-4388-b95a-c273a3b567ff
	// Take the token from the URL query parameter "token" verify it and save the user as a member of the organization

	// === A working example of how to accept an invitation ===
	/* Example:
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}
	if err := gordianService.VerifyInvitation(context.Background(), token); err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	if err := gordianService.AddMemberToOrganization(context.Background(), org.ID, user.ID); err != nil {
		http.Error(w, "Failed to add user to organization", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
	*/

	//Time to simulate the user accepting the invitation
	validateToken, err := gordianService.VerifyInvitation(context.Background(), invite.Token)
	if err != nil {
		log.Fatalf("ERROR: Failed to verify invitation: %v", err)
	}
	if !validateToken {
		log.Fatalf("ERROR: Invalid token")
	}
	//We have two usecases here:
	// 1) User is already a user in the system so we can find them by email and add them to the organization
	// 2) User is not a user in the system so we can create them and add them to the organization
	newMember, err := gordianService.FindUserByEmail(context.Background(), user.Email)
	if err != nil {
		if err == fmt.Errorf("no user found") {
			user, err = gordianService.CreateUser(context.Background(), user.Email, user.Name)
			if err != nil {
				log.Fatalf("ERROR: Failed to create user: %v", err)
			}
		} else {
			log.Fatalf("ERROR: Failed to find user by email: %v", err)
		}
	}
	if err := gordianService.AddMemberToOrganization(context.Background(), org.ID, newMember.ID); err != nil {
		log.Fatalf("ERROR: Failed to add user to organization: %v", err)
		log.Printf("SUCCESS: Added user '%s' to organization '%s'", newMember.Email, org.Name)
	}

	members, err := gordianService.GetMembers(context.Background(), user.ID, org.ID)
	if err != nil {
		log.Fatalf("ERROR: Failed to get organization members: %v", err)
	}
	log.Printf("SUCCESS: Found %d members in organization '%s'", len(members), org.Name)
	for i, member := range members {
		log.Printf("  Member %d: User ID: %s, Role: %s", i+1, member.UserID, member.Role)
	}

	log.Println("Integration test successful!")
}
