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

type AdminController struct {
    AdminGateway *gateways.AdminGateway
}

func NewAdminController(db *mongo.Database) *AdminController {
    return &AdminController{
        AdminGateway: gateways.NewAdminGateway(db),
    }
}

// Handle GET requests for system administrators with pagination
func (ac *AdminController) GetAdmins(w http.ResponseWriter, r *http.Request) {
    limit, _ := r.Context().Value("limit").(int) // Retrieve limit from context
    page, _ := r.Context().Value("page").(int)   // Retrieve page from context
    identifier := r.URL.Query().Get("identifier")

    offset := (page - 1) * limit // Calculate the actual offset

    query := r.URL.Query().Get("name")
    var admins []models.SystemAdmin
    var err error

    if identifier == "search" {
        admins, err = ac.AdminGateway.SearchAdmins(query, limit, offset)
    } else {
        admins, err = ac.AdminGateway.GetAdmins(limit, offset)
    }

    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch system administrators")
        return
    }

    totalEntries, err := ac.AdminGateway.GetTotalAdminCount(query)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch total system administrators count")
        return
    }

    totalPages := int(math.Ceil(float64(totalEntries) / float64(limit)))

    response := map[string]interface{}{
        "data":         admins,
        "total_pages":  totalPages,
        "page": page,
        "total_entries": totalEntries,
    }

    json.NewEncoder(w).Encode(response)
}

// Handle POST requests for system administrators
func (ac *AdminController) CreateAdmin(w http.ResponseWriter, r *http.Request) {
    var admin models.SystemAdmin
    err := json.NewDecoder(r.Body).Decode(&admin)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    err = ac.AdminGateway.CreateAdmin(&admin)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to create system administrator")
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(admin)
}

// Handle PUT requests for system administrators
func (ac *AdminController) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
    var admin models.SystemAdmin
    err := json.NewDecoder(r.Body).Decode(&admin)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    err = ac.AdminGateway.UpdateAdmin(&admin)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to update system administrator")
        return
    }
    json.NewEncoder(w).Encode(admin)
}

// Handle DELETE requests for system administrators
func (ac *AdminController) DeleteAdmin(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    adminID := vars["id"]
    if adminID == "" {
        utils.ErrorHandler(w, http.StatusBadRequest, nil, "Missing admin ID")
        return
    }

    err := ac.AdminGateway.DeleteAdmin(adminID)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to delete system administrator")
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"message": "System administrator deleted successfully"})
}