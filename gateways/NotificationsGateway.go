package gateways

import (
    "database/sql"
    "github.com/BrianKasina/dialysis-scheduling/models"
)

type NotificationGateway struct {
    db *sql.DB
}

func NewNotificationGateway(db *sql.DB) *NotificationGateway {
    return &NotificationGateway{db: db}
}

func (ng *NotificationGateway) GetNotifications(limit, offset int) ([]models.Notification, error) {
    rows, err := ng.db.Query(`
        SELECT n.notification_id, n.message, n.sent_date, n.sent_time, sa.name AS admin_name, p.name AS patient_name
        FROM notifications n
        JOIN system_admin sa ON n.admin_id = sa.admin_id
        JOIN patient p ON n.patient_id = p.patient_id
        LIMIT ? OFFSET ?
    `, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notifications []models.Notification
    for rows.Next() {
        var notification models.Notification
        if err := rows.Scan(&notification.ID, &notification.Message, &notification.SentDate, &notification.SentTime, &notification.AdminName, &notification.PatientName); err != nil {
            return nil, err
        }
        notifications = append(notifications, notification)
    }
    return notifications, nil
}

func (ng *NotificationGateway) SearchNotifications(query string, limit, offset int) ([]models.Notification, error) {
    searchQuery := "%" + query + "%"
    rows, err := ng.db.Query(`
        SELECT n.notification_id, n.message, n.sent_date, n.sent_time, sa.name AS admin_name, p.name AS patient_name
        FROM notifications n
        JOIN system_admin sa ON n.admin_id = sa.admin_id
        JOIN patient p ON n.patient_id = p.patient_id
        WHERE CONCAT(n.message, ' ', sa.name, ' ', p.name) LIKE ?
        LIMIT ? OFFSET ?
    `, searchQuery, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var notifications []models.Notification
    for rows.Next() {
        var notification models.Notification
        if err := rows.Scan(&notification.ID, &notification.Message, &notification.SentDate, &notification.SentTime, &notification.AdminName, &notification.PatientName); err != nil {
            return nil, err
        }
        notifications = append(notifications, notification)
    }
    return notifications, nil
}

func (ng *NotificationGateway) GetTotalNotificationCount(query string) (int, error) {
    var row *sql.Row
    if query != "" {
        searchQuery := "%" + query + "%"
        row = ng.db.QueryRow(`
            SELECT COUNT(*)
            FROM notifications n
            JOIN system_admin sa ON n.admin_id = sa.admin_id
            JOIN patient p ON n.patient_id = p.patient_id
            WHERE CONCAT(n.message, ' ', sa.name, ' ', p.name) LIKE ?
        `, searchQuery)
    } else {
        row = ng.db.QueryRow("SELECT COUNT(*) FROM notifications")
    }

    var count int
    err := row.Scan(&count)
    if err != nil {
        return 0, err
    }

    return count, nil
}

func (ng *NotificationGateway) CreateNotification(notification *models.Notification) error {
    _, err := ng.db.Exec("INSERT INTO notifications (message, sent_date, sent_time, admin_id, patient_id) VALUES (?, ?, ?, ?, ?)",
        notification.Message, notification.SentDate, notification.SentTime, notification.AdminID, notification.PatientID)
    return err
}

func (ng *NotificationGateway) UpdateNotification(notification *models.Notification) error {
    _, err := ng.db.Exec("UPDATE notifications SET message = ?, sent_date = ?, sent_time = ?, admin_id = ?, patient_id = ? WHERE notification_id = ?",
        notification.Message, notification.SentDate, notification.SentTime, notification.AdminID, notification.PatientID, notification.ID)
    return err
}

func (ng *NotificationGateway) DeleteNotification(notificationID string) error {
    _, err := ng.db.Exec("DELETE FROM notifications WHERE notification_id = ?", notificationID)
    return err
}