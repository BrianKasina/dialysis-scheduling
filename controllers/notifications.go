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

type NotificationController struct {
    NotificationGateway *gateways.NotificationGateway
}

func NewNotificationController(db *sql.DB) *NotificationController {
    return &NotificationController{
        NotificationGateway: gateways.NewNotificationGateway(db),
    }
}

// Handle GET requests for notifications with pagination
func (nc *NotificationController) GetNotifications(w http.ResponseWriter, r *http.Request) {
    limit, _ := r.Context().Value("limit").(int) // Retrieve limit from context
    page, _ := r.Context().Value("page").(int)   // Retrieve page from context
    identifier := r.URL.Query().Get("identifier")

    offset := (page - 1) * limit // Calculate the actual offset

    query := r.URL.Query().Get("query")
    var notifications []models.Notification
    var err error

    if identifier == "search" {
        notifications, err = nc.NotificationGateway.SearchNotifications(query, limit, offset)
    } else {
        notifications, err = nc.NotificationGateway.GetNotifications(limit, offset)
    }

    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch notifications")
        return
    }

    totalEntries, err := nc.NotificationGateway.GetTotalNotificationCount(query)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to fetch total notifications count")
        return
    }

    totalPages := int(math.Ceil(float64(totalEntries) / float64(limit)))

    response := map[string]interface{}{
        "data":         notifications,
        "total_pages":  totalPages,
        "page": page,
        "total_entries": totalEntries,
    }

    json.NewEncoder(w).Encode(response)
}

// Handle POST requests for notifications
func (nc *NotificationController) CreateNotification(w http.ResponseWriter, r *http.Request) {
    var notification models.Notification
    err := json.NewDecoder(r.Body).Decode(&notification)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    err = nc.NotificationGateway.CreateNotification(&notification)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to create notification")
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(notification)
}

// Handle PUT requests for notifications
func (nc *NotificationController) UpdateNotification(w http.ResponseWriter, r *http.Request) {
    var notification models.Notification
    err := json.NewDecoder(r.Body).Decode(&notification)
    if err != nil {
        utils.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    err = nc.NotificationGateway.UpdateNotification(&notification)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to update notification")
        return
    }
    json.NewEncoder(w).Encode(notification)
}

// Handle DELETE requests for notifications
func (nc *NotificationController) DeleteNotification(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    notificationID := vars["id"]
    if notificationID == "" {
        utils.ErrorHandler(w, http.StatusBadRequest, nil, "Missing notification ID")
        return
    }

    err := nc.NotificationGateway.DeleteNotification(notificationID)
    if err != nil {
        utils.ErrorHandler(w, http.StatusInternalServerError, err, "Failed to delete notification")
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"message": "Notification deleted successfully"})
}