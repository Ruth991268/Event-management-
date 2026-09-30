package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"event-management/backend/config"
	"event-management/backend/handlers"
	"event-management/backend/middlewares"

	_ "github.com/lib/pq"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("DB unreachable: %v", err)
	}

	log.Println("Database connected successfully")

	// Create HTTP multiplexer
	mux := http.NewServeMux()

	// --------------------
	// Health Check
	// --------------------
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// --------------------
	// Authentication
	// --------------------
	mux.HandleFunc("/signup", handlers.Signup(db, cfg))
	mux.HandleFunc("/login", handlers.Login(db, cfg))

	// --------------------
	// Events Endpoints
	// --------------------

	// Create a new event (POST)
	mux.Handle("/events/create", middlewares.JWTMiddleware(cfg.JWTSecret, handlers.CreateEvent(db)))

	// List events (GET) with optional query params:
	// ?category_id=1&location_id=2&start_date=2025-11-23T00:00:00Z&end_date=2025-11-24T00:00:00Z&min_price=0&max_price=100&keyword=music
	mux.Handle("/events/list", middlewares.JWTMiddleware(cfg.JWTSecret, handlers.ListEvents(db)))

	// Get a single event by id: /events/get?id=3
	mux.Handle("/events/get", middlewares.JWTMiddleware(cfg.JWTSecret, handlers.GetEvent(db)))

	// Protected test route
	mux.Handle("/events/protected", middlewares.JWTMiddleware(cfg.JWTSecret,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(middlewares.UserIDKey).(int64)
			w.Write([]byte("Protected route reached by userID: " + strconv.FormatInt(userID, 10)))
		}),
	))

	// --------------------
	// Start server
	// --------------------
	log.Printf("Server running on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
