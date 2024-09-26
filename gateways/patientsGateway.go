package gateways

import (
	"database/sql"
)

type Patient struct {
	ID   int
	Name string
	Age  int
}

type PatientGateway struct {
	db *sql.DB
}

func NewPatientGateway(db *sql.DB) *PatientGateway {
	return &PatientGateway{db: db}
}

func (g *PatientGateway) GetPatients() ([]Patient, error) {
	rows, err := g.db.Query("SELECT id, name, age FROM patients")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []Patient
	for rows.Next() {
		var p Patient
		if err := rows.Scan(&p.ID, &p.Name, &p.Age); err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}
	return patients, nil
}

func (g *PatientGateway) CreatePatient(patient Patient) error {
	_, err := g.db.Exec("INSERT INTO patients (name, age) VALUES (?, ?)", patient.Name, patient.Age)
	return err
}
