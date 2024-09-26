package controllers

import (
    "encoding/json"
    "net/http"
    "github.com/BrianKasina/dialysis-scheduling/gateways"
    "github.com/BrianKasina/dialysis-scheduling/utils"
    "database/sql"
    "github.com/gorilla/mux"
)

type PatientController struct {
    PatientGateway *gateways.PatientGateway
}

func NewPatientController(db *sql.DB) *PatientController {
    return &PatientController{
        PatientGateway: gateways.NewPatientGateway(db),
    }
}

// Handle GET requests for patients
func (pc *PatientController) GetPatients(w http.ResponseWriter, r *http.Request) {
    identifier := r.URL.Query().Get("identifier")
    var patients interface{}
    var err error

    switch identifier {
    case "withPayment":
        patients, err = pc.PatientGateway.GetPatientsWithPayment()
    case "withDialysisAppointments":
        patients, err = pc.PatientGateway.GetPatientsWithDialysisAppointments()
    case "withNephrologistAppointments":
        patients, err = pc.PatientGateway.GetPatientsWithNephrologistAppointments()
    case "withNotifications":
        patients, err = pc.PatientGateway.GetPatientsWithNotifications()
    default:
        patients, err = pc.PatientGateway.GetPatients()
    }

    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch patients")
        return
    }
    json.NewEncoder(w).Encode(patients)
}

// Handle POST requests for patients
func (pc *PatientController) CreatePatient(w http.ResponseWriter, r *http.Request) {
    // Extract patient data from request body
    var patient gateways.Patient
    err := json.NewDecoder(r.Body).Decode(&patient)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    // Create patient in DB
    err = pc.PatientGateway.CreatePatient(&patient)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to create patient")
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(patient)
}

// Handle PUT requests for patients
func (pc *PatientController) UpdatePatient(w http.ResponseWriter, r *http.Request) {
    // Extract patient data from request body
    var patient gateways.Patient
    err := json.NewDecoder(r.Body).Decode(&patient)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    // Update patient in DB
    err = pc.PatientGateway.UpdatePatient(&patient)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to update patient")
        return
    }
    json.NewEncoder(w).Encode(patient)
}

// Handle DELETE requests for patients
func (pc *PatientController) DeletePatient(w http.ResponseWriter, r *http.Request) {
    // Extract patient ID from request URL
    vars := mux.Vars(r)
    patientID := vars["id"]
    if patientID == "" {
        utils.ErrorHandler(w, http.StatusBadRequest, nil, "Missing patient ID")
        return
    }

    // Delete patient in DB
    err := pc.PatientGateway.DeletePatient(patientID)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to delete patient")
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"message": "Patient deleted successfully"})
}