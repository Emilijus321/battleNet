package handlers

import (
	//"battleNet/models"
	"battleNet/templates"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HandleEditProfilePage - rodo profilio redagavimo formą
func (h *Handler) HandleEditProfilePage(w http.ResponseWriter, r *http.Request) {
	userIDStr := h.sessionManager.GetString(r.Context(), "userID")
	if userIDStr == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusBadRequest)
		return
	}

	// Gauti vartotojo duomenis
	user, err := h.userRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		http.Error(w, "Failed to load profile", http.StatusInternalServerError)
		return
	}

	email := h.sessionManager.GetString(r.Context(), "email")
	role := h.sessionManager.GetString(r.Context(), "role")

	component := templates.EditProfilePage(email, role, user, "")
	component.Render(r.Context(), w)
}

// HandleUpdateProfile - apdoroja profilio atnaujinimą
func (h *Handler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userIDStr := h.sessionManager.GetString(r.Context(), "userID")
	if userIDStr == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusBadRequest)
		return
	}

	// Atnaujinti profilio duomenis
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	username := r.FormValue("username")

	// Gauti vartotoją, kad parodytume formą su error
	user, err := h.userRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		http.Error(w, "Failed to load profile", http.StatusInternalServerError)
		return
	}

	email := h.sessionManager.GetString(r.Context(), "email")
	role := h.sessionManager.GetString(r.Context(), "role")

	// Atnaujinti profilio duomenis
	err = h.userRepo.UpdateUserProfile(r.Context(), userID, firstName, lastName, username)
	if err != nil {
		log.Printf("Error updating profile: %v", err)
		component := templates.EditProfilePage(email, role, user, "Failed to update profile")
		component.Render(r.Context(), w)
		return
	}

	// Atnaujinti vardą sesijoje
	h.sessionManager.Put(r.Context(), "name", firstName+" "+lastName)
	h.sessionManager.Put(r.Context(), "username", username)

	// Nukreipti į profilio puslapį
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// HandleChangePasswordPage - rodo slaptažodžio keitimo formą
func (h *Handler) HandleChangePasswordPage(w http.ResponseWriter, r *http.Request) {
	email := h.sessionManager.GetString(r.Context(), "email")
	role := h.sessionManager.GetString(r.Context(), "role")

	component := templates.ChangePasswordPage(email, role, "")
	component.Render(r.Context(), w)
}

// HandleChangePassword - apdoroja slaptažodžio keitimą
func (h *Handler) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userIDStr := h.sessionManager.GetString(r.Context(), "userID")
	if userIDStr == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusBadRequest)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	email := h.sessionManager.GetString(r.Context(), "email")
	role := h.sessionManager.GetString(r.Context(), "role")

	// Patikrinti ar nauji slaptažodžiai sutampa
	if newPassword != confirmPassword {
		component := templates.ChangePasswordPage(email, role, "New passwords do not match")
		component.Render(r.Context(), w)
		return
	}

	// Gauti vartotojo dabartinį slaptažodį
	user, err := h.userRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Patikrinti dabartinį slaptažodį
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		component := templates.ChangePasswordPage(email, role, "Current password is incorrect")
		component.Render(r.Context(), w)
		return
	}

	// Sukurti naują slaptažodžio hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		component := templates.ChangePasswordPage(email, role, "Failed to update password")
		component.Render(r.Context(), w)
		return
	}

	// Atnaujinti duomenų bazėje
	err = h.userRepo.UpdatePassword(r.Context(), userID, string(hashedPassword))
	if err != nil {
		log.Printf("Error updating password: %v", err)
		component := templates.ChangePasswordPage(email, role, "Failed to update password")
		component.Render(r.Context(), w)
		return
	}

	// Sėkmės puslapis
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Password Changed</title>
		<link rel="stylesheet" href="/static/style.css">
		<meta http-equiv="refresh" content="3;url=/profile">
	</head>
	<body>
		<div class="container">
			<nav>
				<div class="nav-links">
					<a href="/">Home</a>
					<a href="/dashboard">Dashboard</a>
					<a href="/profile">Profile</a>
				</div>
				<a href="/logout">Logout</a>
			</nav>
			
			<div class="content">
				<div class="alert alert-success">
					✓ Password changed successfully!
				</div>
				<p>Redirecting to profile page in 3 seconds...</p>
				<a href="/profile" class="btn">← Go to Profile now</a>
			</div>
		</div>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
