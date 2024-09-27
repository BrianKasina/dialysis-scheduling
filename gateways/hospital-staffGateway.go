package gateways

import (
    "database/sql"
    "github.com/BrianKasina/dialysis-scheduling/models"
)

type HospitalStaffGateway struct {
    db *sql.DB
}

func NewHospitalStaffGateway(db *sql.DB) *HospitalStaffGateway {
    return &HospitalStaffGateway{db: db}
}

func (hsg *HospitalStaffGateway) GetHospitalStaff(limit, offset int) ([]models.HospitalStaff, error) {
    rows, err := hsg.db.Query("SELECT * FROM hospital_staff LIMIT ? OFFSET ?", limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var staff []models.HospitalStaff
    for rows.Next() {
        var member models.HospitalStaff
        if err := rows.Scan(&member.ID, &member.Name, &member.Gender, &member.Specialization, &member.PhoneNumber, &member.Status); err != nil {
            return nil, err
        }
        staff = append(staff, member)
    }
    return staff, nil
}

func (hsg *HospitalStaffGateway) SearchHospitalStaff(query string, limit, offset int) ([]models.HospitalStaff, error) {
    searchQuery := "%" + query + "%"
    rows, err := hsg.db.Query(`
        SELECT * FROM hospital_staff
        WHERE CONCAT(name, ' ', specialization, ' ', phonenumber) LIKE ? LIMIT ? OFFSET ?
    `, searchQuery, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var staff []models.HospitalStaff
    for rows.Next() {
        var member models.HospitalStaff
        if err := rows.Scan(&member.ID, &member.Name, &member.Gender, &member.Specialization, &member.PhoneNumber, &member.Status); err != nil {
            return nil, err
        }
        staff = append(staff, member)
    }
    return staff, nil
}

func (hsg *HospitalStaffGateway) GetTotalStaffCount(query string) (int, error) {
    var row *sql.Row
    if query != "" {
        searchQuery := "%" + query + "%"
        row = hsg.db.QueryRow(`
            SELECT COUNT(*) FROM hospital_staff
            WHERE CONCAT(name, ' ', specialization, ' ', phonenumber) LIKE ?
        `, searchQuery)
    } else {
        row = hsg.db.QueryRow("SELECT COUNT(*) FROM hospital_staff")
    }

    var count int
    err := row.Scan(&count)
    if err != nil {
        return 0, err
    }

    return count, nil
}

func (hsg *HospitalStaffGateway) CreateHospitalStaff(member *models.HospitalStaff) error {
    _, err := hsg.db.Exec("INSERT INTO hospital_staff (name, gender, specialization, phonenumber, status) VALUES (?, ?, ?, ?, ?)",
        member.Name, member.Gender, member.Specialization, member.PhoneNumber, member.Status)
    return err
}

func (hsg *HospitalStaffGateway) UpdateHospitalStaff(member *models.HospitalStaff) error {
    _, err := hsg.db.Exec(
        `UPDATE hospital_staff SET name = ?, gender = ?, specialization = ?, phonenumber = ?, status = ? 
        WHERE staff_id = (SELECT staff_id FROM hospital_staff WHERE name = ?)`,
        member.Name, member.Gender, member.Specialization, member.PhoneNumber, member.Status, member.Name)
    return err
}

func (hsg *HospitalStaffGateway) DeleteHospitalStaff(staffID string) error {
    _, err := hsg.db.Exec("DELETE FROM hospital_staff WHERE staff_id = ?", staffID)
    return err
}