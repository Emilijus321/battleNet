package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"battleNet/config"
	"battleNet/internal/handlers"
	"battleNet/middlewaree"
	"battleNet/repository"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	sessionManager *scs.SessionManager
	cfg            *config.Config
	db             *repository.Database
)

func main() {
	// Load configuration
	cfg = config.Load()

	// Connect to database
	var err error
	db, err = repository.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize session manager
	initSessionManager()

	// Create repository instances
	userRepo := repository.NewUserRepository(db.Pool)
	movieRepo := repository.NewMovieRepository(db.Pool)
	reviewRepo := repository.NewReviewRepository(db.Pool)
	watchlistRepo := repository.NewWatchlistRepository(db.Pool)

	// Initialize handlers
	handler := handlers.NewHandler(userRepo, movieRepo, reviewRepo, watchlistRepo, cfg.JWTSecret, sessionManager)

	// Setup router
	router := setupRouter(handler)

	// Start server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("🚀 Server starting on port %s", cfg.Port)
		log.Printf("🌐 Frontend: http://localhost:%s", cfg.Port)
		log.Printf("🔗 API: http://localhost:%s/api/v1", cfg.Port)
		log.Printf("📊 Environment: %s", cfg.Environment)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server shutdown failed: %v", err)
	}

	log.Println("✅ Server stopped")
}

func initSessionManager() {
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = false
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
}

func setupRouter(handler *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(sessionManager.LoadAndSave)

	// CORS middlewaree
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	// Public routes
	r.Get("/", handler.HandleHome)
	r.Get("/login", handler.HandleLoginPage)
	r.Post("/login", handler.HandleLogin)
	r.Get("/signup", handler.HandleSignupPage)
	r.Post("/signup", handler.HandleSignup)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middlewaree.RequireAuth(sessionManager))

		r.Get("/dashboard", handler.HandleDashboard)
		r.Get("/profile", handler.HandleProfile)
		r.Get("/logout", handler.HandleLogout)
		r.Get("/movies", handler.HandleMovies)
		r.Get("/movies/{id}", handler.HandleMovieDetail)
		r.Get("/watchlist", handler.HandleWatchlist)
		r.Post("/watchlist/add", handler.HandleAddToWatchlist)
		r.Post("/reviews", handler.HandleCreateReview)

		// Admin routes
		r.Group(func(r chi.Router) {
			r.Use(middlewaree.RequireRole(sessionManager, "admin"))

			r.Get("/admin/movies", handler.HandleAdminMovies)
			r.Get("/admin/movies/create", handler.HandleCreateMoviePage)
			r.Post("/admin/movies/create", handler.HandleCreateMovie)
			r.Get("/admin/users", handler.HandleAdminUsers)
		})
	})

	// API routes (REST API)
	r.Route("/api/v1", func(r chi.Router) {
		// Public API endpoints
		r.Get("/movies", handler.HandleAPIMovies)
		r.Get("/movies/{id}", handler.HandleAPIMovieDetail)
		r.Get("/reviews", handler.HandleAPIReviews)

		// Protected API endpoints
		r.Group(func(r chi.Router) {
			r.Use(middlewaree.RequireAuthAPI(sessionManager))

			r.Post("/reviews", handler.HandleAPICreateReview)
			r.Get("/watchlist", handler.HandleAPIWatchlist)
			r.Post("/watchlist", handler.HandleAPIAddToWatchlist)
			r.Delete("/watchlist/{movieId}", handler.HandleAPIRemoveFromWatchlist)
		})

		// Admin API endpoints
		r.Group(func(r chi.Router) {
			r.Use(middlewaree.RequireAuthAPI(sessionManager))
			r.Use(middlewaree.RequireRoleAPI(sessionManager, "admin"))

			r.Post("/movies", handler.HandleAPICreateMovie)
			r.Put("/movies/{id}", handler.HandleAPIUpdateMovie)
			r.Delete("/movies/{id}", handler.HandleAPIDeleteMovie)
		})
	})

	return r
}
