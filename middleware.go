package gordian

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

const (
	ActiveOrgIDKey        = "active_org_id"
	ActiveRoleKey         = "active_role"
	ActiveMembershipIDKey = "active_membership_id"
)

func (s *Service) TenancyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID_from_ctx := r.Context().Value("user_id") // Using a generic "user_id" key for example
		if userID_from_ctx == nil {
			// This is a server error. The auth middleware should have run.
			log.Println("ERROR: user_id not found in context. Is the auth middleware missing?")
			http.Error(w, "Server Configuration Error", http.StatusInternalServerError)
			return
		}
		userID, ok := userID_from_ctx.(uuid.UUID)
		if !ok {
			http.Error(w, "Invalid user_id in context", http.StatusInternalServerError)
			return
		}

		tenantID_from_header := r.Header.Get("X-Tenant-ID")
		if tenantID_from_header == "" {
			http.Error(w, "Missing X-Tenant-ID header", http.StatusBadRequest)
			return
		}
		orgID, err := uuid.Parse(tenantID_from_header)
		if err != nil {
			http.Error(w, "Invalid X-Tenant-ID header", http.StatusBadRequest)
			return
		}

		membershipID, role, err := s.GetMemberships(r.Context(), userID, orgID)
		if err != nil {
			http.Error(w, "Failed to get membership", http.StatusInternalServerError)
			return
		}

		if role == "" || role != "admin" {
			http.Error(w, "User not a member of organization", http.StatusForbidden)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ActiveOrgIDKey, orgID)
		ctx = context.WithValue(ctx, ActiveRoleKey, role)
		ctx = context.WithValue(ctx, ActiveMembershipIDKey, membershipID)

		next.ServeHTTP(w, r)
	})
}
