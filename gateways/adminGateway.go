package gateways

import (
    "database/sql"
    "github.com/BrianKasina/dialysis-scheduling/models"
)

type AdminGateway struct {
    db *sql.DB
}

func NewAdminGateway(db *sql.DB) *AdminGateway {
    return &AdminGateway{db: db}
}

func (ag *AdminGateway) GetAdmins(limit, offset int) ([]models.SystemAdmin, error) {
    rows, err := ag.db.Query("SELECT * FROM system_admin LIMIT ? OFFSET ?", limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var admins []models.SystemAdmin
    for rows.Next() {
        var admin models.SystemAdmin
        if err := rows.Scan(&admin.ID, &admin.Name, &admin.Email, &admin.PhoneNumber); err != nil {
            return nil, err
        }
        admins = append(admins, admin)
    }
    return admins, nil
}

func (ag *AdminGateway) SearchAdmins(query string, limit, offset int) ([]models.SystemAdmin, error) {
    searchQuery := "%" + query + "%"
    rows, err := ag.db.Query(`
        SELECT * FROM system_admin
        WHERE CONCAT(name, ' ', email, ' ', phonenumber) LIKE ? LIMIT ? OFFSET ?
    `, searchQuery, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var admins []models.SystemAdmin
    for rows.Next() {
        var admin models.SystemAdmin
        if err := rows.Scan(&admin.ID, &admin.Name, &admin.Email, &admin.PhoneNumber); err != nil {
            return nil, err
        }
        admins = append(admins, admin)
    }
    return admins, nil
}

func (ag *AdminGateway) GetTotalAdminCount(query string) (int, error) {
    var row *sql.Row
    if query != "" {
        searchQuery := "%" + query + "%"
        row = ag.db.QueryRow(`
            SELECT COUNT(*) FROM system_admin
            WHERE CONCAT(name, ' ', email, ' ', phonenumber) LIKE ?
        `, searchQuery)
    } else {
        row = ag.db.QueryRow("SELECT COUNT(*) FROM system_admin")
    }

    var count int
    err := row.Scan(&count)
    if err != nil {
        return 0, err
    }

    return count, nil
}

func (ag *AdminGateway) CreateAdmin(admin *models.SystemAdmin) error {
    _, err := ag.db.Exec("INSERT INTO system_admin (name, email, phonenumber) VALUES (?, ?, ?)",
        admin.Name, admin.Email, admin.PhoneNumber)
    return err
}

func (ag *AdminGateway) UpdateAdmin(admin *models.SystemAdmin) error {
    _, err := ag.db.Exec("UPDATE system_admin SET name = ?, email = ?, phonenumber = ? WHERE admin_id = ?",
        admin.Name, admin.Email, admin.PhoneNumber, admin.ID)
    return err
}

func (ag *AdminGateway) DeleteAdmin(adminID string) error {
    _, err := ag.db.Exec("DELETE FROM system_admin WHERE admin_id = ?", adminID)
    return err
}