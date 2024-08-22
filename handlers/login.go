package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anilsaini81155/exchangeccurrency/config"
	"github.com/anilsaini81155/exchangeccurrency/models"

	"github.com/anilsaini81155/exchangeccurrency/utils"
)

// Mock user database
// var users = map[string]string{
// 	"username": "abc@g.com",
// 	"password": "12ds34",
// }

// Login logs in a user and returns a JWT token
func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	var users = config.Users
	json.NewDecoder(r.Body).Decode(&user)

	storedPassword, exists := users[user.Username]
	if !exists || !utils.CheckPasswordHash(user.Password, storedPassword["password"]) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
