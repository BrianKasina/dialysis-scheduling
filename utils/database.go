package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Host     string
	Name     string
	User     string
	Password string
}

func NewDatabase(host, name, user, password string) *Database {
	return &Database{Host: host, Name: name, User: user, Password: password}
}

// GetConnection establishes and returns a new DB connection
func (db *Database) GetConnection() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", db.User, db.Password, db.Host, db.Name)
	return sql.Open("mysql", dsn)
}
