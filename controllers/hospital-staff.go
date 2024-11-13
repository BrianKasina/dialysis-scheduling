package controllers

import (
    "encoding/json"
    "net/http"
    "github.com/BrianKasina/dialysis-scheduling/gateways"
    "github.com/BrianKasina/dialysis-scheduling/utils"
    "github.com/BrianKasina/dialysis-scheduling/models"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/gorilla/mux"
	"math"
)

type HospitalStaffController struct {
    HospitalStaffGateway *gateways.HospitalStaffGateway
}

func NewHospitalStaffController(db *mongo.Database) *HospitalStaffController {
    return &HospitalStaffController{
        HospitalStaffGateway: gateways.NewHospitalStaffGateway(db),
    }
}

// Handle GET requests for hospital staff
func (hsc *HospitalStaffController) GetHospitalStaff(w http.ResponseWriter, r *http.Request) {
    limit, _ := r.Context().Value("limit").(int) // Retrieve limit from context
    page, _ := r.Context().Value("page").(int)   // Retrieve page from context
	identifier := r.URL.Query().Get("identifier")

    offset := (page - 1) * limit // Calculate the actual offset

    query := r.URL.Query().Get("name")
    var staff []models.HospitalStaff
    var err error

    if identifier == "search" {
        staff, err = hsc.HospitalStaffGateway.SearchHospitalStaff(query, limit, offset)
    } else {
        staff, err = hsc.HospitalStaffGateway.GetHospitalStaff(limit, offset)
    }

    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch hospital staff")
        return
    }

    totalEntries, err := hsc.HospitalStaffGateway.GetTotalStaffCount(query)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch total hospital staff count")
        return
    }

    totalPages := int(math.Ceil(float64(totalEntries) / float64(limit)))

    response := map[string]interface{}{
        "data":         staff,
        "total_pages":  totalPages,
        "page": page,
        "total_entries": totalEntries,
    }

    json.NewEncoder(w).Encode(response)
}


// Handle POST requests for hospital staff
func (hsc *HospitalStaffController) CreateHospitalStaff(w http.ResponseWriter, r *http.Request) {
    var member models.HospitalStaff
    err := json.NewDecoder(r.Body).Decode(&member)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    err = hsc.HospitalStaffGateway.CreateHospitalStaff(&member)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to create hospital staff")
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(member)
}

// Handle PUT requests for hospital staff
func (hsc *HospitalStaffController) UpdateHospitalStaff(w http.ResponseWriter, r *http.Request) {
    var member models.HospitalStaff
    err := json.NewDecoder(r.Body).Decode(&member)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    err = hsc.HospitalStaffGateway.UpdateHospitalStaff(&member)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to update hospital staff")
        return
    }
    json.NewEncoder(w).Encode(member)
}

// Handle DELETE requests for hospital staff
func (hsc *HospitalStaffController) DeleteHospitalStaff(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    staffID := vars["id"]
    if staffID == "" {
        utils.ErrorHandler(w, http.StatusBadRequest, nil, "Missing staff ID")
        return
    }

    err := hsc.HospitalStaffGateway.DeleteHospitalStaff(staffID)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to delete hospital staff")
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"message": "Hospital staff deleted successfully"})
}