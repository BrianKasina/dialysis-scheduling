package main

import (
	"context"
	"errors"
	"log"
	"net/http"
    "strconv"
	"os"
	"github.com/BrianKasina/dialysis-scheduling/controllers"
	"github.com/BrianKasina/dialysis-scheduling/utils"
	"github.com/gorilla/mux")

// Middleware to extract pagination parameters
func paginationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		pageStr := r.URL.Query().Get("page")

		limit := 10 // Default limit
		page := 1   // Default page

		if limitStr != "" {
			if limitVal, err := strconv.Atoi(limitStr); err == nil {
				limit = limitVal
			}
		}
		if pageStr != "" {
			if pageVal, err := strconv.Atoi(pageStr); err == nil {
				page = pageVal
			}
		}

		// Store limit and page in request context
		ctx := context.WithValue(r.Context(), "limit", limit)
		ctx = context.WithValue(ctx, "page", page)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware to set Content-Type header
func setJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// // Load environment variables from .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file: %v", err)
	// }

	// Load environment variables
	dbHost := os.Getenv("MONGO_HOST")
	dbPort := os.Getenv("MONGO_PORT")
	dbName := os.Getenv("MONGO_DATABASE")
	dbUser := os.Getenv("MONGO_USER")
	dbPass := os.Getenv("MONGO_PASSWORD")

	// Initialize database connection
	database, err := utils.NewDatabase(dbHost, dbPort, dbName, dbUser, dbPass)
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.GetConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Initialize controllers
	controllersMap := map[string]interface{}{
		"patients":        controllers.NewPatientController(db),
		"appointments":    controllers.NewAppointmentController(db),
		"hospital_staff":  controllers.NewHospitalStaffController(db),
		"system_admins":   controllers.NewAdminController(db),
		"notifications":   controllers.NewNotificationController(db),
		"posts":           controllers.NewPostController(db),
		"payment_details": controllers.NewPaymentDetailsController(db),
		"patient_history": controllers.NewPatientHistoryController(db),
	}

	// Initialize router
	router := mux.NewRouter()

	// Apply middleware to set Content-Type header
	router.Use(utils.CorsMiddleware)
	router.Use(setJSONContentType)
	router.Use(paginationMiddleware)

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

	}).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Handle GET requests
func handleGetRequest(w http.ResponseWriter, r *http.Request, endpoint string, controllersMap map[string]interface{}) {
	switch endpoint {
	case "patients":
		controllersMap["patients"].(*controllers.PatientController).GetPatients(w, r)
	case "appointments":
		controllersMap["appointments"].(*controllers.AppointmentController).GetAppointments(w, r)
	case "hospital_staff":
		controllersMap["hospital_staff"].(*controllers.HospitalStaffController).GetHospitalStaff(w, r)
	case "system_admins":
		controllersMap["system_admins"].(*controllers.AdminController).GetAdmins(w, r)
	case "notifications":
		controllersMap["notifications"].(*controllers.NotificationController).GetNotifications(w, r)
	case "posts":
		controllersMap["posts"].(*controllers.PostController).GetPosts(w, r)
	case "payment_details":
		controllersMap["payment_details"].(*controllers.PaymentDetailsController).GetPaymentDetails(w, r)
	case "patient_history":
		controllersMap["patient_history"].(*controllers.PatientHistoryController).HandlePatientHistory(w, r)
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
	case "hospital_staff":
		controllersMap["hospital_staff"].(*controllers.HospitalStaffController).CreateHospitalStaff(w, r)
	case "system_admins":
		controllersMap["system_admins"].(*controllers.AdminController).CreateAdmin(w, r)
	case "notifications":
		controllersMap["notifications"].(*controllers.NotificationController).CreateNotification(w, r)
	case "posts":
		controllersMap["posts"].(*controllers.PostController).CreatePost(w, r)
	case "payment_details":
		controllersMap["payment_details"].(*controllers.PaymentDetailsController).CreatePaymentDetail(w, r)
	case "patient_history":
		controllersMap["patient_history"].(*controllers.PatientHistoryController).UploadPatientHistory(w, r)
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
	case "hospital_staff":
		controllersMap["hospital_staff"].(*controllers.HospitalStaffController).UpdateHospitalStaff(w, r)
	case "system_admins":
		controllersMap["system_admins"].(*controllers.AdminController).UpdateAdmin(w, r)
	case "notifications":
		controllersMap["notifications"].(*controllers.NotificationController).UpdateNotification(w, r)
	case "posts":
		controllersMap["posts"].(*controllers.PostController).UpdatePost(w, r)
	case "payment_details":
		controllersMap["payment_details"].(*controllers.PaymentDetailsController).UpdatePaymentDetail(w, r)
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
	case "hospital_staff":
		controllersMap["hospital_staff"].(*controllers.HospitalStaffController).DeleteHospitalStaff(w, r)
	case "system_admins":
		controllersMap["system_admins"].(*controllers.AdminController).DeleteAdmin(w, r)
	case "notifications":
		controllersMap["notifications"].(*controllers.NotificationController).DeleteNotification(w, r)
	case "posts":
		controllersMap["posts"].(*controllers.PostController).DeletePost(w, r)
	case "payment_details":
		controllersMap["payment_details"].(*controllers.PaymentDetailsController).DeletePaymentDetail(w, r)
	default:
		http.Error(w, "Invalid endpoint", http.StatusNotFound)
	}
}
