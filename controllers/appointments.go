package controllers

import (
    "encoding/json"
    "net/http"
    "math"
    "github.com/BrianKasina/dialysis-scheduling/gateways"
    "github.com/BrianKasina/dialysis-scheduling/utils"
    "go.mongodb.org/mongo-driver/mongo"
)

// AppointmentController struct to manage both dialysis and nephrologist appointments
type AppointmentController struct {
    DialysisGateway       *gateways.DialysisGateway
    NephrologistGateway   *gateways.NephrologistAppointmentGateway
}

// NewAppointmentController creates a new AppointmentController instance
func NewAppointmentController(db *mongo.Database) *AppointmentController {
    return &AppointmentController{
        DialysisGateway:     gateways.NewDialysisGateway(db),
        NephrologistGateway: gateways.NewNephrologistAppointmentGateway(db),
    }
}

// Handle GET requests for appointments, type distinguishes between dialysis and nephrologist
func (ac *AppointmentController) GetAppointments(w http.ResponseWriter, r *http.Request) {
    limit := r.Context().Value("limit").(int)
    page := r.Context().Value("page").(int)
    appointmentType := r.URL.Query().Get("type")
    identifier := r.URL.Query().Get("identifier")

    offset := (page - 1) * limit

    query := r.URL.Query().Get("query")
    var appointments interface{}
    var err error

    switch identifier {
    case "dialysis":
        appointments, err = ac.DialysisGateway.GetAppointments( limit, offset )
    case "nephrologist":
        appointments, err = ac.NephrologistGateway.GetAppointments( limit, offset )
    case "search":
        query := r.URL.Query().Get("name")
        if query == "" {
            utils.ErrorHandler(w, http.StatusBadRequest, nil, "Missing search query")
            return
        }
        if appointmentType == "dialysis" {
            appointments, err = ac.DialysisGateway.SearchAppointments(query, limit, offset)
        } else if appointmentType == "nephrologist" {
            appointments, err = ac.NephrologistGateway.SearchAppointments(query, limit, offset)
        } else {
            utils.ErrorHandler(w, http.StatusBadRequest, nil, "Invalid appointment type for search")
            return
        }
    default:
        utils.ErrorHandler(w, http.StatusBadRequest, nil, "Invalid appointment type")
        return
    }

    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch appointments")
        return
    }

    var totalEntries int
    if appointmentType == "dialysis" {
        totalEntries, err = ac.DialysisGateway.GetTotalDialysisAppointmentCount(query)
    } else if appointmentType == "nephrologist" {
        totalEntries, err = ac.NephrologistGateway.GetTotalNephrologistAppointmentCount(query)
    }

    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch total appointments count")
        return
    }

    totalPages := int(math.Ceil(float64(totalEntries) / float64(limit)))

    response := map[string]interface{}{
        "data":         appointments,
        "total_pages":  totalPages,
        "page": page,
        "total_entries": totalEntries,
    }

    json.NewEncoder(w).Encode(response)
}

// Handle POST requests for creating appointments
func (ac *AppointmentController) CreateAppointment(w http.ResponseWriter, r *http.Request) {
    appointmentType := r.URL.Query().Get("type")

    switch appointmentType {
    case "dialysis":
        ac.DialysisGateway.CreateAppointment(w, r)
    case "nephrologist":
        ac.NephrologistGateway.CreateAppointment(w, r)
    default:
        utils.ErrorHandler(w, http.StatusBadRequest, nil, "Invalid appointment type")
    }
}

// Handle PUT requests for updating appointments
func (ac *AppointmentController) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
    appointmentType := r.URL.Query().Get("type")

    switch appointmentType {
    case "dialysis":
        ac.DialysisGateway.UpdateAppointment(w, r)
    case "nephrologist":
        ac.NephrologistGateway.UpdateAppointment(w, r)
    default:
        utils.ErrorHandler(w, http.StatusBadRequest, nil, "Invalid appointment type")
    }
}

// Handle DELETE requests for deleting appointments
func (ac *AppointmentController) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
    appointmentType := r.URL.Query().Get("type")

    switch appointmentType {
    case "dialysis":
        ac.DialysisGateway.DeleteAppointment(w, r)
    case "nephrologist":
        ac.NephrologistGateway.DeleteAppointment(w, r)
    default:
        utils.ErrorHandler(w, http.StatusBadRequest, nil, "Invalid appointment type")
    }
}
