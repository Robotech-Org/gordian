# Gordian Documentation

## 1. Overview

Gordian is a multi-tenancy toolkit for Go designed to simplify the management of organizations, user roles, and invitations. It provides a clean, decoupled architecture that allows you to easily integrate tenancy into your Go applications. The library is built with a "batteries-included but replaceable" philosophy, offering a ready-to-use GORM adapter while allowing you to implement your own storage or service adapters.

The primary goal of Gordian is to handle the common backend logic for SaaS applications, such as:
- Creating and managing tenant accounts (Organizations).
- Associating users with multiple organizations.
- Assigning roles to users within a specific organization.
- A complete invitation flow for adding new members.

## 2. Core Concepts

Gordian's data model is built around four primary entities.

```go gordian/tenancy.go#L8-89
// User represents a global user in the system.
// A user is NOT a tenant. They are an identity who can join tenants.
type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	CreatedAt time.Time
}

// Organization is the Tenant. It is the top-level container for users and resources.
type Organization struct {
	ID        uuid.UUID
	Name      string
	OwnerID   uuid.UUID // The user who created and owns the organization
	CreatedAt time.Time
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
```

-   **Organization**: This represents the tenant. It is the top-level container for resources and members.
-   **User**: This is a global entity. A single `User` account can have access to multiple `Organization`s.
-   **Membership**: This is the crucial link between a `User` and an `Organization`. It defines that a user is part of an organization and specifies their `Role` (e.g., "owner", "admin", "editor").
-   **Invite**: This entity manages the process of inviting a user (by email) to join a specific organization with a designated role.

## 3. Architecture

Gordian uses a layered architecture to separate concerns, making it flexible and testable.

-   **Service Layer (`gordian.go`)**: The `gordian.Service` is the core of the library. It contains all the business logic for managing users, organizations, and memberships. It orchestrates the data stores and other services like the emailer.

-   **Interfaces (`stores.go`)**: The service layer depends on a set of interfaces for its data persistence and external communication needs. This decouples the business logic from the implementation details.
    -   `OrganizationStore`: Handles `Organization` persistence.
    -   `UserStore`: Handles `User` persistence.
    -   `MembershipStore`: Handles `Membership` persistence.
    -   `InvitationStore`: Handles `Invite` persistence.
    -   `Emailer`: Defines a contract for sending emails, such as invitations.

-   **Adapters (`adapter/`)**: Adapters are concrete implementations of the store interfaces. Gordian provides a `gorm` adapter out of the box.
    -   `gordian/adapter/gorm/gorm.go`: This package provides GORM-based implementations for all the store interfaces, designed to work with a PostgreSQL database. You can easily create your own adapters for different databases (e.g., MongoDB, SQLC) by implementing the interfaces defined in `stores.go`.

## 4. Getting Started & Example Usage

The `cmd/example-server` directory provides a working example of how to initialize and use the Gordian service.

### Step 1: Configuration

Create a `.env` file in the root of your project with your database and email credentials.

```gordian/.env
DATABASE_URL=postgres://user:password@localhost:5432/gordian_db?sslmode=disable
EMAIL_HOST=sandbox.smtp.mailtrap.io
EMAIL_PORT=2525
EMAIL_USERNAME=your-username
EMAIL_PASSWORD=your-password
EMAIL_AUTH=PLAIN
```

### Step 2: Initialization

In your application's entry point, you need to connect to your database, initialize the stores, create an emailer, and inject them into the `gordian.Service`.

```go gordian/cmd/example-server/main.go#L30-L54
// --- Step 1: Connect to the database ---
dsn := os.Getenv("DATABASE_URL")
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//...

// --- Configure Mailtrap Emailer ---
mailtrapCfg := mailtrap.Config{
    Host:     os.Getenv("EMAIL_HOST"),
    Port:     os.Getenv("EMAIL_PORT"),
    Username: os.Getenv("EMAIL_USERNAME"),
    Password: os.Getenv("EMAIL_PASSWORD"),
    FromAddr: "noreply@yourapp.com",
}
emailer := mailtrap.NewEmailer(mailtrapCfg)

db.AutoMigrate(&gordian.User{}, &gordian.Organization{}, &gordian.Membership{}, &gordian.Invite{})


// --- Step 2: Initialize the adapter and inject it into the service ---
orgStore := gormadapter.NewOrganizationStore(db)
userStore := gormadapter.NewUserStore(db)
memStore := gormadapter.NewMembershipStore(db)
invStore := gormadapter.NewInviteStore(db)

gordianService := gordian.New(orgStore, userStore, memStore, invStore, emailer)
```

### Step 3: Core Operations

Once the service is initialized, you can use it to perform multi-tenancy operations.

#### Create a User and Organization
When an organization is created, the user who creates it is automatically assigned the "owner" role.

```go gordian/cmd/example-server/main.go#L60-L70
// Create a user first
user, err := gordianService.CreateUser(context.Background(), "test.user@example.com", "Test User")
// ...

// Now create an organization with that user as the owner
org, err := gordianService.CreateOrganization(context.Background(), "My First Test Org", user.ID)
// ...
```

#### The Invitation Flow
The invitation process involves creating an invite, sending it, verifying it, and finally adding the user as a member.

1.  **Create and Send Invitation**: An existing member of an organization creates an invitation for a new user's email address. Gordian generates a unique token and uses the `Emailer` to send the invite.

    ```go gordian/cmd/example-server/main.go#L81-L85
    invite, err := gordianService.CreateInvitation(context.Background(), org.ID, user.ID, "new.colleague@example.com", "editor")
    if err != nil {
        log.Fatalf("ERROR: Failed to create and send invitation: %v", err)
    }
    ```

2.  **Accept Invitation**: When the invited user clicks the link in their email, your application will receive the token. You must then verify it. If valid, you check if a user with that email already exists. If not, you create one. Finally, you create a membership to add them to the organization.

    ```go gordian/cmd/example-server/main.go#L95-L121
    //Simulate the user accepting the invitation
    validateToken, err := gordianService.VerifyInvitation(context.Background(), invite.Token)
    // ...

    //We have two usecases here:
    // 1) User is already a user in the system so we can find them by email and add them to the organization
    // 2) User is not a user in the system so we can create them and add them to the organization
    newMember, err := gordianService.FindUserByEmail(context.Background(), user.Email)
    if err != nil {
        if err == fmt.Errorf("no user found") {
            user, err = gordianService.CreateUser(context.Background(), user.Email, user.Name)
            //...
        }
        //...
    }
    if err := gordianService.AddMemberToOrganization(context.Background(), org.ID, newMember.ID); err != nil {
        //...
    }
    ```

## 5. Tenancy Middleware

Gordian provides an HTTP middleware to enforce tenancy at the request level.

```go gordian/middleware.go#L14-L57
func (s *Service) TenancyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ... (implementation)
	})
}
```

The `TenancyMiddleware` performs the following actions:
1.  Extracts the `user_id` from the request context (this assumes a preceding authentication middleware has already identified the user).
2.  Reads the `X-Tenant-ID` header from the request to determine which organization the user is trying to act on.
3.  Verifies that the user is a member of the specified organization.
4.  If the check passes, it injects the active `OrganizationID`, `Role`, and `MembershipID` into the context for downstream handlers to use.

This ensures that all subsequent logic in your request handler is correctly scoped to a single tenant and that the user has the appropriate permissions.