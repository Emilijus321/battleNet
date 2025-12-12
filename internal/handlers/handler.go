package handlers

import (
	"battleNet/repository"
	"net/http"

	//"battleNet/templates"

	"github.com/alexedwards/scs/v2"
)

type Handler struct {
	userRepo       *repository.UserRepository
	movieRepo      *repository.MovieRepository
	reviewRepo     *repository.ReviewRepository
	watchlistRepo  *repository.WatchlistRepository
	jwtSecret      string
	sessionManager *scs.SessionManager
}

func NewHandler(
	userRepo *repository.UserRepository,
	movieRepo *repository.MovieRepository,
	reviewRepo *repository.ReviewRepository,
	watchlistRepo *repository.WatchlistRepository,
	jwtSecret string,
	sessionManager *scs.SessionManager,
) *Handler {
	return &Handler{
		userRepo:       userRepo,
		movieRepo:      movieRepo,
		reviewRepo:     reviewRepo,
		watchlistRepo:  watchlistRepo,
		jwtSecret:      jwtSecret,
		sessionManager: sessionManager,
	}
}
func (h *Handler) HandleAdminMovies(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Admin Movies Page"))
}

func (h *Handler) HandleCreateMoviePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create Movie Page"))
}

func (h *Handler) HandleCreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Movie Created"))
}

func (h *Handler) HandleAdminUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Admin Users Page"))
}

func (h *Handler) HandleAPICreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status":"success"}`))
}
func (h *Handler) HandleAPIUpdateMovie(w http.ResponseWriter, r *http.Request) {
	// TODO: čia įrašyk logiką filmų atnaujinimui per API
	w.Write([]byte(`{"status":"success", "message":"Movie updated"}`))
}

func (h *Handler) HandleAPIDeleteMovie(w http.ResponseWriter, r *http.Request) {
	// TODO: čia įrašyk logiką filmų ištrynimui per API
	w.Write([]byte(`{"status":"success", "message":"Movie deleted"}`))
}
