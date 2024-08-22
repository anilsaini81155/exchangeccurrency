package handlers

import (
	"encoding/json"

	"github.com/anilsaini81155/exchangeccurrency/config"
	"github.com/anilsaini81155/exchangeccurrency/models"

	"net/http"

	"github.com/anilsaini81155/exchangeccurrency/utils"
)

// Register registers a new user
func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	var users = config.Users
	json.NewDecoder(r.Body).Decode(&user)

	if _, exists := users[user.Username]; exists {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	_, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// users[user.Username] = hashedPassword

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User created")
}
