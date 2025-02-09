// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"Bite-Speed/controllers"
	"Bite-Speed/database"

	"github.com/gorilla/mux"
)

func main() {
	// Initializing DB
	database.ConnectDB()
	defer database.DB.Close()

	// Creating Route
	r := mux.NewRouter()
	r.HandleFunc("/identify", controllers.IdentifyHandler).Methods("POST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
