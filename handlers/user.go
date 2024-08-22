package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anilsaini81155/exchangeccurrency/utils"
)

// GetUser gets the authenticated user
func GetUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenString := r.Header.Get("Authorization")
	claims, err := utils.ValidateJWT(tokenString)

	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"username": claims.Username})
}
