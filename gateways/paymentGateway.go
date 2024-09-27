package gateways

import (
    "database/sql"
    "github.com/BrianKasina/dialysis-scheduling/models"
)

type PaymentDetailsGateway struct {
    db *sql.DB
}

func NewPaymentDetailsGateway(db *sql.DB) *PaymentDetailsGateway {
    return &PaymentDetailsGateway{db: db}
}

func (pg *PaymentDetailsGateway) GetPaymentDetails(limit, offset int) ([]models.PaymentDetails, error) {
    rows, err := pg.db.Query(`
        SELECT pd.payment_details_id, pd.payment_name
        FROM payment_details pd
        LIMIT ? OFFSET ?
    `, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var paymentDetails []models.PaymentDetails
    for rows.Next() {
        var paymentDetail models.PaymentDetails
        if err := rows.Scan(&paymentDetail.ID, &paymentDetail.PaymentName); err != nil {
            return nil, err
        }
        paymentDetails = append(paymentDetails, paymentDetail)
    }
    return paymentDetails, nil
}

func (pg *PaymentDetailsGateway) SearchPaymentDetails(query string, limit, offset int) ([]models.PaymentDetails, error) {
    searchQuery := "%" + query + "%"
    rows, err := pg.db.Query(`
        SELECT pd.payment_details_id, pd.payment_name
        FROM payment_details pd
        WHERE CONCAT(pd.payment_name) LIKE ?
        LIMIT ? OFFSET ?
    `, searchQuery, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var paymentDetails []models.PaymentDetails
    for rows.Next() {
        var paymentDetail models.PaymentDetails
        if err := rows.Scan(&paymentDetail.ID, &paymentDetail.PaymentName); err != nil {
            return nil, err
        }
        paymentDetails = append(paymentDetails, paymentDetail)
    }
    return paymentDetails, nil
}

func (pg *PaymentDetailsGateway) GetTotalPaymentDetailsCount(query string) (int, error) {
    var row *sql.Row
    if query != "" {
        searchQuery := "%" + query + "%"
        row = pg.db.QueryRow(`
            SELECT COUNT(*)
            FROM payment_details pd
            WHERE CONCAT(pd.payment_name) LIKE ?
        `, searchQuery)
    } else {
        row = pg.db.QueryRow("SELECT COUNT(*) FROM payment_details")
    }

    var count int
    err := row.Scan(&count)
    if err != nil {
        return 0, err
    }

    return count, nil
}

func (pg *PaymentDetailsGateway) CreatePaymentDetail(paymentDetail *models.PaymentDetails) error {
    _, err := pg.db.Exec("INSERT INTO payment_details (payment_name) VALUES (?)",
        paymentDetail.PaymentName)
    return err
}

func (pg *PaymentDetailsGateway) UpdatePaymentDetail(paymentDetail *models.PaymentDetails) error {
    _, err := pg.db.Exec(`UPDATE payment_details SET payment_name = ? WHERE payment_details_id = ?`,
        paymentDetail.PaymentName, paymentDetail.ID)
    return err
}

func (pg *PaymentDetailsGateway) DeletePaymentDetail(paymentDetailID string) error {
    _, err := pg.db.Exec("DELETE FROM payment_details WHERE payment_details_id = ?", paymentDetailID)
    return err
}