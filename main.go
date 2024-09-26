package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	"dialysis-scheduling/utils" // Import the new utilities
)

func main() {
	// Load environment variables
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	secretKey := os.Getenv("JWT_SECRET_KEY")

	// Initialize database connection
	database := utils.NewDatabase(dbHost+":"+dbPort, dbName, dbUser, dbPass)
	db, err := database.GetConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// // Initialize JWT utility
	// jwtUtil := utils.NewJWTUtil(secretKey)

	// Initialize router
	router := mux.NewRouter()

	// // Add CORS middleware
	// router.Use(CorsMiddleware)

	// Define routes
	router.HandleFunc("/{endpoint}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		endpoint := vars["endpoint"]

		// Check if the endpoint is valid and allowed
		if _, ok := allowedEndpoints[endpoint]; !ok {
			ErrorHandler(w, r, http.StatusNotFound, "Endpoint not found")
			return
		}

		// Handle GET, POST, PUT, DELETE requests
		switch r.Method {
		case http.MethodGet:
			handleGetRequest(w, r, endpoint, db)
		case http.MethodPost:
			handlePostRequest(w, r, endpoint, db, jwtUtil)
		case http.MethodPut:
			handlePutRequest(w, r, endpoint, db)
		case http.MethodDelete:
			handleDeleteRequest(w, r, endpoint, db)
		default:
			ErrorHandler(w, r, http.StatusMethodNotAllowed, "Method not allowed")
		}

	}).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Handle GET requests
func handleGetRequest(w http.ResponseWriter, r *http.Request, endpoint string, db *sql.DB) {
	// Logic for GET
	json.NewEncoder(w).Encode(map[string]string{"message": "GET request for " + endpoint})
}

// Handle POST requests
func handlePostRequest(w http.ResponseWriter, r *http.Request, endpoint string, db *sql.DB, jwtUtil *utils.JWTUtil) {
	// Logic for POST with JWT encoding
	// payload := map[string]interface{}{"endpoint": endpoint, "timestamp": time.Now().Unix()}
	// token, err := jwtUtil.Encode(payload, time.Hour)
	// if err != nil {
	// 	ErrorHandler(w, r, http.StatusInternalServerError, "Failed to create JWT")
	// 	return
	// }

	json.NewEncoder(w).Encode(map[string]string{"message": "POST request for " + endpoint, "token": token})
}

// Handle PUT requests
func handlePutRequest(w http.ResponseWriter, r *http.Request, endpoint string, db *sql.DB) {
	// Logic for PUT
	json.NewEncoder(w).Encode(map[string]string{"message": "PUT request for " + endpoint})
}

// Handle DELETE requests
func handleDeleteRequest(w http.ResponseWriter, r *http.Request, endpoint string, db *sql.DB) {
	// Logic for DELETE
	json.NewEncoder(w).Encode(map[string]string{"message": "DELETE request for " + endpoint})
}
