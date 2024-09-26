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

// Middleware to set Content-Type header
func setJSONContentType(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}

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
    controllersMap := map[string]interface{}{
        "patients":     controllers.NewPatientController(db),
        "appointments": controllers.NewAppointmentController(db),
    }

    // Initialize router
    router := mux.NewRouter()

    // Apply middleware to set Content-Type header
    router.Use(setJSONContentType)

    allowedEndpoints := map[string]bool{
        "patients":        true,
        "hospital_staff":  true,
        "appointments":    true,
        "posts":           true,
        "system_admins":   true,
        "notifications":   true,
        "patient_history": true,
        "payment_details": true,
    }

    // Define routes
    router.HandleFunc("/{endpoint}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        endpoint := vars["endpoint"]

        // Check if the endpoint is valid and allowed
        if _, ok := allowedEndpoints[endpoint]; !ok {
            utils.ErrorHandler(w, http.StatusNotFound, err, "Endpoint not found")
            return
        }

        // Dispatch the request to the appropriate controller method
        switch r.Method {
        case http.MethodGet:
            handleGetRequest(w, r, endpoint, controllersMap)
        case http.MethodPost:
            handlePostRequest(w, r, endpoint, controllersMap)
        case http.MethodPut:
            handlePutRequest(w, r, endpoint, controllersMap)
        case http.MethodDelete:
            handleDeleteRequest(w, r, endpoint, controllersMap)
        default:
            utils.ErrorHandler(w, http.StatusMethodNotAllowed, errors.New("invalid method"), "Method not allowed")
        }

    }).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete)

    log.Fatal(http.ListenAndServe(":8080", router))
}

// Handle GET requests
func handleGetRequest(w http.ResponseWriter, r *http.Request, endpoint string, controllersMap map[string]interface{}) {
    switch endpoint {
    case "patients":
        controllersMap["patients"].(*controllers.PatientController).GetPatients(w, r)
    case "appointments":
        controllersMap["appointments"].(*controllers.AppointmentController).GetAppointments(w, r)
    default:
        http.Error(w, "Invalid endpoint", http.StatusNotFound)
    }
}

// Handle POST requests
func handlePostRequest(w http.ResponseWriter, r *http.Request, endpoint string, controllersMap map[string]interface{}) {
    switch endpoint {
    case "patients":
        controllersMap["patients"].(*controllers.PatientController).CreatePatient(w, r)
    case "appointments":
        controllersMap["appointments"].(*controllers.AppointmentController).CreateAppointment(w, r)
    default:
        http.Error(w, "Invalid endpoint", http.StatusNotFound)
    }
}

// Handle PUT requests
func handlePutRequest(w http.ResponseWriter, r *http.Request, endpoint string, controllersMap map[string]interface{}) {
    switch endpoint {
    case "patients":
        controllersMap["patients"].(*controllers.PatientController).UpdatePatient(w, r)
    case "appointments":
        controllersMap["appointments"].(*controllers.AppointmentController).UpdateAppointment(w, r)
    default:
        http.Error(w, "Invalid endpoint", http.StatusNotFound)
    }
}

// Handle DELETE requests
func handleDeleteRequest(w http.ResponseWriter, r *http.Request, endpoint string, controllersMap map[string]interface{}) {
    switch endpoint {
    case "patients":
        controllersMap["patients"].(*controllers.PatientController).DeletePatient(w, r)
    case "appointments":
        controllersMap["appointments"].(*controllers.AppointmentController).DeleteAppointment(w, r)
    default:
        http.Error(w, "Invalid endpoint", http.StatusNotFound)
    }
}