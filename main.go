package main

import (
    "errors"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/BrianKasina/dialysis-scheduling/controllers"
    "github.com/BrianKasina/dialysis-scheduling/utils"
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

    // Initialize database connection
    database := utils.NewDatabase(dbHost+":"+dbPort, dbName, dbUser, dbPass)
    db, err := database.GetConnection()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Initialize controllers
    patientController := controllers.NewPatientController(db)

    // Initialize router
    router := mux.NewRouter()

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
            utils.ErrorHandler(w, http.StatusNotFound, errors.New("endpoint not found"), "Endpoint not found")
            return
        }

        // Handle GET, POST, PUT, DELETE requests
        switch r.Method {
        case http.MethodGet:
            handleGetRequest(w, r, endpoint, patientController)
        case http.MethodPost:
            handlePostRequest(w, r, endpoint, patientController)
        case http.MethodPut:
            handlePutRequest(w, r, endpoint, patientController)
        case http.MethodDelete:
            handleDeleteRequest(w, r, endpoint, patientController)
        default:
            utils.ErrorHandler(w, http.StatusMethodNotAllowed, errors.New("invalid method"), "Method not allowed")
        }

    }).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)

    log.Fatal(http.ListenAndServe(":8080", router))
}

// Handle GET requests
func handleGetRequest(w http.ResponseWriter, r *http.Request, endpoint string, patientController *controllers.PatientController) {
    switch endpoint {
    case "patients":
        patientController.GetPatients(w, r)
    default:
        http.Error(w, "Invalid endpoint", http.StatusNotFound)
    }
}

// Handle POST requests
func handlePostRequest(w http.ResponseWriter, r *http.Request, endpoint string, patientController *controllers.PatientController) {
    switch endpoint {
    case "patients":
        patientController.CreatePatient(w, r)
    default:
        http.Error(w, "Invalid endpoint", http.StatusNotFound)
    }
}

// Handle PUT requests
func handlePutRequest(w http.ResponseWriter, r *http.Request, endpoint string, patientController *controllers.PatientController) {
    switch endpoint {
    case "patients":
        patientController.UpdatePatient(w, r)
    default:
        http.Error(w, "Invalid endpoint", http.StatusNotFound)
    }
}

// Handle DELETE requests
func handleDeleteRequest(w http.ResponseWriter, r *http.Request, endpoint string, patientController *controllers.PatientController) {
    switch endpoint {
    case "patients":
        patientController.DeletePatient(w, r)
    default:
        http.Error(w, "Invalid endpoint", http.StatusNotFound)
    }
}