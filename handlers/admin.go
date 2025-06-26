package handlers

import (
	"encoding/json"
	"gofamtree/config"
	"gofamtree/models"
	"gofamtree/utils"
	"net/http"
	"time"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string       `json:"message"`
	Admin   models.Admin `json:"admin"`
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var admin models.Admin
	if err := config.DB.Where("username = ?", input.Username).First(&admin).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(input.Password, admin.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		Message: "Login successful",
		Admin:   admin,
	})
}

func AdminRegister(w http.ResponseWriter, r *http.Request) {
	var input RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if username already exists
	var existingAdmin models.Admin
	if err := config.DB.Where("username = ?", input.Username).First(&existingAdmin).Error; err == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	admin := models.Admin{
		Username:  input.Username,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		http.Error(w, "Failed to create admin", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Admin created successfully",
		"admin":   admin,
	})
}
