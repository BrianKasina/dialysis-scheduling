package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"errors"
	"github.com/joho/godotenv"
	"github.com/BrianKasina/dialysis-scheduling/utils" 
	"github.com/BrianKasina/dialysis-scheduling/controllers"
	"github.com/BrianKasina/dialysis-scheduling/gateways"// Import the new utilities

)

func main() {
	// Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	// Load environment variables
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	// secretKey := os.Getenv("JWT_SECRET_KEY")

	// Initialize database connection
	database := utils.NewDatabase(dbHost+":"+dbPort, dbName, dbUser, dbPass)
	db, err := database.GetConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// // Initialize JWT utility
	// jwtUtil := utils.NewJWTUtil(secretKey)

	// Initialize the error handler utility 
	ErrorHandler := func(w http.ResponseWriter, statusCode int, err error, message string) {
		utils.ErrorHandler(w, statusCode, err, message)
	}

	// Initialize router
	router := mux.NewRouter()

	// // Add CORS middleware
	// router.Use(CorsMiddleware)

	allowedEndpoints := map[string]bool{
		"patients":               true,
		"hospital_staff":         true,
		"dialysis_appointments":  true,
		"nephrologist_appointments": true,
		"posts":                  true,
		"system_admins":          true,
		"notifications":          true,
		"patient_history":        true,
		"payment_details":        true,
	}

	// Define routes
	router.HandleFunc("/{endpoint}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		endpoint := vars["endpoint"]

		// Check if the endpoint is valid and allowed
		if _, ok := allowedEndpoints[endpoint]; !ok {
			ErrorHandler(w, http.StatusNotFound, errors.New("endpoint not found"), "Endpoint not found")
			return
		}

		// Handle GET, POST, PUT, DELETE requests
		switch r.Method {
		case http.MethodGet:
			handleGetRequest(w, r, endpoint, db)
		case http.MethodPost:
			handlePostRequest(w, r, endpoint, db)
		case http.MethodPut:
			handlePutRequest(w, r, endpoint, db)
		case http.MethodDelete:
			handleDeleteRequest(w, r, endpoint, db)
		default:
			ErrorHandler(w, http.StatusMethodNotAllowed, errors.New("invalid method"),"Method not allowed")
		}

	}).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Handle GET requests
func handleGetRequest(w http.ResponseWriter, r *http.Request, endpoint string, db *sql.DB) {
	switch endpoint {
	case "patients":
		// Logic for GET patients
		json.NewEncoder(w).Encode(map[string]string{"message": "GET request for " + endpoint})
	case "hospital_staff":
		// Logic for GET hospital_staff
		json.NewEncoder(w).Encode(map[string]string{"message": "GET request for " + endpoint})
	case "dialysis_appointments":
		// Logic for GET dialysis_appointments
		json.NewEncoder(w).Encode(map[string]string{"message": "GET request for " + endpoint})
	case "nephrologist_appointments":
		// Logic for GET nephrologist_appointments
		json.NewEncoder(w).Encode(map[string]string{"message": "GET request for " + endpoint})
	case "payment_details":
		// Logic for GET payment_details
		json.NewEncoder(w).Encode(map[string]string{"message": "GET request for " + endpoint})
	default:
		// Default response
		//indicate wrong endpoint
		http.Error(w, "Invalid endpoint", http.StatusNotFound)
	}
}

// Handle POST requests
func handlePostRequest(w http.ResponseWriter, r *http.Request, endpoint string, db *sql.DB) {
	// Logic for POST with JWT encoding
	// payload := map[string]interface{}{"endpoint": endpoint, "timestamp": time.Now().Unix()}
	// token, err := jwtUtil.Encode(payload, time.Hour)
	// if err != nil {
	// 	ErrorHandler(w, r, http.StatusInternalServerError, "Failed to create JWT")
	// 	return
	// }
	switch endpoint {
	case "patients":
		// Logic for POST patients
		json.NewEncoder(w).Encode(map[string]string{"message": "POST request for " + endpoint})
	case "hospital_staff":
		// Logic for POST hospital_staff
		json.NewEncoder(w).Encode(map[string]string{"message": "POST request for " + endpoint})
	case "dialysis_appointments":
		// Logic for POST dialysis_appointments
		json.NewEncoder(w).Encode(map[string]string{"message": "POST request for " + endpoint})
	case "nephrologist_appointments":
		// Logic for POST nephrologist_appointments
		json.NewEncoder(w).Encode(map[string]string{"message": "POST request for " + endpoint})
	case "payment_details":
		// Logic for POST payment_details
		json.NewEncoder(w).Encode(map[string]string{"message": "POST request for " + endpoint})
	default:
		// Default response
		//indicate wrong endpoint
		http.Error(w, "Invalid endpoint", http.StatusNotFound)
	}
}

// Handle PUT requests
func handlePutRequest(w http.ResponseWriter, r *http.Request, endpoint string, db *sql.DB) {
  switch endpoint {
	case "patients":
		// Logic for PUT patients
		json.NewEncoder(w).Encode(map[string]string{"message": "PUT request for " + endpoint})
	case "hospital_staff":
		// Logic for PUT hospital_staff
		json.NewEncoder(w).Encode(map[string]string{"message": "PUT request for " + endpoint})
	case "dialysis_appointments":
		// Logic for PUT dialysis_appointments
		json.NewEncoder(w).Encode(map[string]string{"message": "PUT request for " + endpoint})
	case "nephrologist_appointments":
		// Logic for PUT nephrologist_appointments
		json.NewEncoder(w).Encode(map[string]string{"message": "PUT request for " + endpoint})
	case "payment_details":
		// Logic for PUT payment_details
		json.NewEncoder(w).Encode(map[string]string{"message": "PUT request for " + endpoint})
	default:
		// Default response
		//indicate wrong endpoint
		http.Error(w, "Invalid endpoint", http.StatusNotFound)
  }
}

// Handle DELETE requests
func handleDeleteRequest(w http.ResponseWriter, r *http.Request, endpoint string, db *sql.DB) {

	switch endpoint {
	case "patients":
		// Logic for DELETE patients
		json.NewEncoder(w).Encode(map[string]string{"message": "DELETE request for " + endpoint})
	case "hospital_staff":
		// Logic for DELETE hospital_staff
		json.NewEncoder(w).Encode(map[string]string{"message": "DELETE request for " + endpoint})
	case "dialysis_appointments":
		// Logic for DELETE dialysis_appointments
		json.NewEncoder(w).Encode(map[string]string{"message": "DELETE request for " + endpoint})
	case "nephrologist_appointments":
		// Logic for DELETE nephrologist_appointments
		json.NewEncoder(w).Encode(map[string]string{"message": "DELETE request for " + endpoint})
	case "payment_details":
		// Logic for DELETE payment_details
		json.NewEncoder(w).Encode(map[string]string{"message": "DELETE request for " + endpoint})
	default:
		// Default response
		//indicate wrong endpoint
		http.Error(w, "Invalid endpoint", http.StatusNotFound)
	}
}
