package main

import (
	"log"
	"net/http"

	"github.com/anilsaini81155/exchangeccurrency/router"
)

func main() {
	r := router.SetupRouter()
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
