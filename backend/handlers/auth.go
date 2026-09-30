package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"event-management/backend/config"
	"event-management/backend/middlewares"
	"event-management/backend/models"

	"golang.org/x/crypto/bcrypt"
)

func Signup(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		if req.Email == "" || req.Password == "" || req.Name == "" {
			http.Error(w, "Missing fields", http.StatusBadRequest)
			return
		}

		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		var userID int64
		err := db.QueryRow(`
			INSERT INTO users (name, email, password)
			VALUES ($1, $2, $3)
			RETURNING id
		`, req.Name, req.Email, string(hashed)).Scan(&userID)
		if err != nil {
			http.Error(w, "Email already exists", http.StatusBadRequest)
			return
		}

		token, _ := middlewares.GenerateJWT(userID, cfg.JWTSecret)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Signup successful",
			"user": map[string]interface{}{
				"id":    userID,
				"name":  req.Name,
				"email": req.Email,
			},
			"token": token,
		})
	}
}
func Login(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
		if req.Email == "" || req.Password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		var user models.User
		err := db.QueryRow(`SELECT id, name, email, password FROM users WHERE email=$1`, req.Email).
			Scan(&user.ID, &user.Name, &user.Email, &user.Password)

		if err != nil {
			if err == sql.ErrNoRows {
				// Log for debugging
				log.Println("Login failed, user not found:", req.Email)
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			} else {
				log.Println("Login DB error:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Compare password hash
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			log.Println("Login failed, invalid password for:", req.Email)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Generate JWT
		token, err := middlewares.GenerateJWT(user.ID, cfg.JWTSecret)
		if err != nil {
			log.Println("JWT generation error:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Login successful",
			"user": map[string]interface{}{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
			"token": token,
		})
	}
}
