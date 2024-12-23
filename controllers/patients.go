package controllers

import (
    "encoding/json"
    "net/http"
    "math"
    "github.com/BrianKasina/dialysis-scheduling/gateways"
    "github.com/BrianKasina/dialysis-scheduling/utils"
	"github.com/BrianKasina/dialysis-scheduling/models"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/gorilla/mux"
)

type PatientController struct {
    PatientGateway *gateways.PatientGateway
}

func NewPatientController(db *mongo.Database) *PatientController {
    return &PatientController{
        PatientGateway: gateways.NewPatientGateway(db),
    }
}

// Handle GET requests for patients
func (pc *PatientController) GetPatients(w http.ResponseWriter, r *http.Request) {
    limit := r.Context().Value("limit").(int)
    page := r.Context().Value("page").(int)
    identifier := r.URL.Query().Get("identifier")
    query := r.URL.Query().Get("query")

    offset := (page - 1) * limit

    var patients interface{}
    var err error

    switch identifier {
    case "history":
        patients, err = pc.PatientGateway.GetPatientsWithHistory( limit, offset )
	case "search":
        name := r.URL.Query().Get("name")
        if name == "" {
            patients, err = pc.PatientGateway.GetPatients( limit, offset )
			if err != nil {
				utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch patients")
				return
			}
            return
        }
        patients, err = pc.PatientGateway.SearchPatients(name, limit, offset)
    default:
        patients, err = pc.PatientGateway.GetPatients( limit, offset )
    }

    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch patients")
        return
    }

    totalEntries, err := pc.PatientGateway.GetTotalPatientCount(query)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch total patients count")
        return
    }

    totalPages := int(math.Ceil(float64(totalEntries) / float64(limit)))

    response := map[string]interface{}{
        "data":         patients,
        "total_pages":  totalPages,
        "page": page,
        "total_entries": totalEntries,
    }
    json.NewEncoder(w).Encode(response)
}

// Handle POST requests for patients
func (pc *PatientController) CreatePatient(w http.ResponseWriter, r *http.Request) {
    // Extract patient data from request body
    var patient models.Patient
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
    var patient models.Patient
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