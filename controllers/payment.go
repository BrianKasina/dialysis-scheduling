package controllers

import (
    "encoding/json"
    "net/http"
    "math"
    "github.com/BrianKasina/dialysis-scheduling/gateways"
    "github.com/BrianKasina/dialysis-scheduling/utils"
    "github.com/BrianKasina/dialysis-scheduling/models"
    "database/sql"
    "github.com/gorilla/mux"
)

type PaymentDetailsController struct {
    PaymentDetailsGateway *gateways.PaymentDetailsGateway
}

func NewPaymentDetailsController(db *sql.DB) *PaymentDetailsController {
    return &PaymentDetailsController{
        PaymentDetailsGateway: gateways.NewPaymentDetailsGateway(db),
    }
}

// Handle GET requests for payment details with pagination
func (pc *PaymentDetailsController) GetPaymentDetails(w http.ResponseWriter, r *http.Request) {
    limit, _ := r.Context().Value("limit").(int) // Retrieve limit from context
    page, _ := r.Context().Value("page").(int)   // Retrieve page from context
    identifier := r.URL.Query().Get("identifier")

    offset := (page - 1) * limit // Calculate the actual offset

    query := r.URL.Query().Get("query")
    var paymentDetails []models.PaymentDetails
    var err error

    if identifier == "search" {
        paymentDetails, err = pc.PaymentDetailsGateway.SearchPaymentDetails(query, limit, offset)
    } else {
        paymentDetails, err = pc.PaymentDetailsGateway.GetPaymentDetails(limit, offset)
    }

    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch payment details")
        return
    }

    totalEntries, err := pc.PaymentDetailsGateway.GetTotalPaymentDetailsCount(query)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch total payment details count")
        return
    }

    totalPages := int(math.Ceil(float64(totalEntries) / float64(limit)))

    response := map[string]interface{}{
        "data":         paymentDetails,
        "total_pages":  totalPages,
        "page": page,
        "total_entries": totalEntries,
    }

    json.NewEncoder(w).Encode(response)
}

// Handle POST requests for payment details
func (pc *PaymentDetailsController) CreatePaymentDetail(w http.ResponseWriter, r *http.Request) {
    var paymentDetail models.PaymentDetails
    err := json.NewDecoder(r.Body).Decode(&paymentDetail)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    err = pc.PaymentDetailsGateway.CreatePaymentDetail(&paymentDetail)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to create payment detail")
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(paymentDetail)
}

// Handle PUT requests for payment details
func (pc *PaymentDetailsController) UpdatePaymentDetail(w http.ResponseWriter, r *http.Request) {
    var paymentDetail models.PaymentDetails
    err := json.NewDecoder(r.Body).Decode(&paymentDetail)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    err = pc.PaymentDetailsGateway.UpdatePaymentDetail(&paymentDetail)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to update payment detail")
        return
    }
    json.NewEncoder(w).Encode(paymentDetail)
}

// Handle DELETE requests for payment details
func (pc *PaymentDetailsController) DeletePaymentDetail(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    paymentDetailID := vars["id"]
    if paymentDetailID == "" {
        utils.ErrorHandler(w, http.StatusBadRequest, nil, "Missing payment detail ID")
        return
    }

    err := pc.PaymentDetailsGateway.DeletePaymentDetail(paymentDetailID)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to delete payment detail")
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"message": "Payment detail deleted successfully"})
}