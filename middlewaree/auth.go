package middlewaree

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	EmailKey  contextKey = "email"
	RoleKey   contextKey = "role"
)

// - Middleware for web routes
func RequireAuth(sm *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authenticated := sm.GetBool(r.Context(), "authenticated")
			if !authenticated {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, sm.GetString(r.Context(), "userID"))
			ctx = context.WithValue(ctx, EmailKey, sm.GetString(r.Context(), "email"))
			ctx = context.WithValue(ctx, RoleKey, sm.GetString(r.Context(), "role"))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuthAPI - Middleware for API routes
func RequireAuthAPI(sm *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authenticated := sm.GetBool(r.Context(), "authenticated")
			if !authenticated {
				http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, sm.GetString(r.Context(), "userID"))
			ctx = context.WithValue(ctx, EmailKey, sm.GetString(r.Context(), "email"))
			ctx = context.WithValue(ctx, RoleKey, sm.GetString(r.Context(), "role"))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Middleware for specific roles (web)
func RequireRole(sm *scs.SessionManager, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := sm.GetString(r.Context(), "role")
			hasAccess := false

			for _, role := range allowedRoles {
				if userRole == role {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				http.Error(w, "Access Denied - Insufficient privileges", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Middleware for specific roles (API)
func RequireRoleAPI(sm *scs.SessionManager, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := sm.GetString(r.Context(), "role")
			hasAccess := false

			for _, role := range allowedRoles {
				if userRole == role {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				http.Error(w, `{"error": "Access Denied - Insufficient privileges"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
