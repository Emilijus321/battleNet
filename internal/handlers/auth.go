package handlers

import (
	"battleNet/models"
	"log"
	"net/http"

	"battleNet/templates"

	"golang.org/x/crypto/bcrypt"
)

// HandleLoginPage displays login page
func (h *Handler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	component := templates.LoginPage()
	component.Render(r.Context(), w)
}

// HandleLogin processes login form
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Get user from database
	user, err := h.userRepo.GetUserByEmail(r.Context(), email)
	if err != nil {
		log.Printf("Login failed for email %s: %v", email, err)
		component := templates.LoginPageWithError("Invalid email or password")
		component.Render(r.Context(), w)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Printf("Invalid password for user %s", email)
		component := templates.LoginPageWithError("Invalid email or password")
		component.Render(r.Context(), w)
		return
	}

	// Update last login
	if err := h.userRepo.UpdateLastLogin(r.Context(), user.UserID); err != nil {
		log.Printf("Failed to update last login for user %s: %v", user.UserID, err)
	}

	// Set session data
	h.sessionManager.Put(r.Context(), "userID", user.UserID.String())
	h.sessionManager.Put(r.Context(), "email", user.Email)
	h.sessionManager.Put(r.Context(), "role", user.Role)
	h.sessionManager.Put(r.Context(), "name", user.FirstName+" "+user.LastName)
	h.sessionManager.Put(r.Context(), "authenticated", true)

	log.Printf("User logged in: %s (role: %s)", user.Email, user.Role)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// HandleSignupPage displays signup page
func (h *Handler) HandleSignupPage(w http.ResponseWriter, r *http.Request) {
	component := templates.SignupPage()
	component.Render(r.Context(), w)
}

// HandleSignup processes signup form
func (h *Handler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	signupReq := models.SignupRequest{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
		FirstName:       r.FormValue("first_name"),
		LastName:        r.FormValue("last_name"),
		Username:        r.FormValue("username"),
	}

	// Validate passwords match
	if signupReq.Password != signupReq.ConfirmPassword {
		component := templates.SignupPageWithError("Passwords do not match")
		component.Render(r.Context(), w)
		return
	}

	// Check if user already exists
	_, err := h.userRepo.GetUserByEmail(r.Context(), signupReq.Email)
	if err == nil {
		component := templates.SignupPageWithError("Email already registered")
		component.Render(r.Context(), w)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create user
	user := &models.User{
		Email:        signupReq.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    signupReq.FirstName,
		LastName:     signupReq.LastName,
		Username:     signupReq.Username,
		Role:         "user", // Default role
		IsActive:     true,
	}

	if err := h.userRepo.CreateUser(r.Context(), user); err != nil {
		log.Printf("Failed to create user: %v", err)
		component := templates.SignupPageWithError("Failed to create account")
		component.Render(r.Context(), w)
		return
	}

	// Auto-login after signup
	h.sessionManager.Put(r.Context(), "userID", user.UserID.String())
	h.sessionManager.Put(r.Context(), "email", user.Email)
	h.sessionManager.Put(r.Context(), "role", user.Role)
	h.sessionManager.Put(r.Context(), "name", user.FirstName+" "+user.LastName)
	h.sessionManager.Put(r.Context(), "authenticated", true)

	log.Printf("New user registered: %s (id: %s)", user.Email, user.UserID)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// HandleLogout logs out user
func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	email := h.sessionManager.GetString(r.Context(), "email")
	h.sessionManager.Destroy(r.Context())
	log.Printf("User logged out: %s", email)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
