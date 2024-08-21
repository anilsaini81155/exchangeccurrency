package router

import (
	"github.com/anilsaini81155/exchangeccurrency/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/register", handlers.Register).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/user", handlers.GetUser).Methods("GET")
	// router.HandleFunc("/rate", handlers.StoreExchangeRate).Methods("POST").Handler(middleware.JWTAuth(http.HandlerFunc(handlers.StoreExchangeRate)))
	router.HandleFunc("/rate", handlers.StoreExchangeRate).Methods("POST")
	router.HandleFunc("/rate", handlers.GetExchangeRate).Methods("GET")

	return router
}
