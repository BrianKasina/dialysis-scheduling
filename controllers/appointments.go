package controllers

import (
    "encoding/json"
    "net/http"
    "github.com/BrianKasina/dialysis-scheduling/gateways"
    "github.com/BrianKasina/dialysis-scheduling/utils"
    "database/sql"
)

// AppointmentController struct to manage both dialysis and nephrologist appointments
type AppointmentController struct {
    DialysisGateway       *gateways.DialysisGateway
    NephrologistGateway   *gateways.NephrologistAppointmentGateway
}

// NewAppointmentController creates a new AppointmentController instance
func NewAppointmentController(db *sql.DB) *AppointmentController {
    return &AppointmentController{
        DialysisGateway:     gateways.NewDialysisGateway(db),
        NephrologistGateway: gateways.NewNephrologistAppointmentGateway(db),
    }
}

// Handle GET requests for appointments, type distinguishes between dialysis and nephrologist
func (ac *AppointmentController) GetAppointments(w http.ResponseWriter, r *http.Request) {
    appointmentType := r.URL.Query().Get("type")
    identifier := r.URL.Query().Get("identifier")
    var appointments interface{}
    var err error

    switch identifier {
    case "dialysis":
        appointments, err = ac.DialysisGateway.GetAppointments()
    case "nephrologist":
        appointments, err = ac.NephrologistGateway.GetAppointments()
    case "search":
        query := r.URL.Query().Get("name")
        if query == "" {
            utils.ErrorHandler(w, http.StatusBadRequest, nil, "Missing search query")
            return
        }
        if appointmentType == "dialysis" {
            appointments, err = ac.DialysisGateway.SearchAppointments(query)
        } else if appointmentType == "nephrologist" {
            appointments, err = ac.NephrologistGateway.SearchAppointments(query)
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
    json.NewEncoder(w).Encode(appointments)
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
