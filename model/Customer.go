package model

import (
	"database/sql"
	"nyala-backend/usecase/viewmodel"
	"strings"
	"time"
)

var (
	// DefaultCustomerBy ...
	DefaultCustomerBy = "def.updated_at"
	// CustomerBy ...
	CustomerBy = []string{
		"def.created_at", "def.updated_at",
	}

	customerSelectString = `SELECT def.customer_id, def.customer_name, def.email, def.phone_number, def.dob, def.sex, def."password", def.created_at, def.updated_at, def.deleted_at FROM customers def`
)

func (model customerModel) scanRow(row *sql.Row) (d CustomerEntity, err error) {
	err = row.Scan(
		&d.CustomerID, &d.CustomerName, &d.Email, &d.PhoneNumber, &d.Dob, &d.Sex, &d.Password, &d.CreatedAt,
		&d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

// customerModel ...
type customerModel struct {
	DB *sql.DB
}

// ICustomer ...
type ICustomer interface {
	FindByEmailOrPhoneNumber(email, phoneNumber string) (CustomerEntity, error)
	FindByID(id string) (CustomerEntity, error)
	Store(body viewmodel.CustomerVM, changedAt time.Time) (string, error)
}

// CustomerEntity ....
type CustomerEntity struct {
	CustomerID   string         `db:"customer_id"`
	CustomerName string         `db:"customer_name"`
	Email        string         `db:"email"`
	PhoneNumber  string         `db:"phone_number"`
	Dob          string         `db:"dob"`
	Sex          bool           `db:"sex"`
	Password     string         `db:"password"`
	CreatedAt    string         `db:"created_at"`
	UpdatedAt    string         `db:"updated_at"`
	DeletedAt    sql.NullString `db:"deleted_at"`
}

// NewCustomerModel ...
func NewCustomerModel(db *sql.DB) ICustomer {
	return &customerModel{DB: db}
}

// FindByEmailOrMail ...
func (model customerModel) FindByEmailOrPhoneNumber(email, phoneNumber string) (res CustomerEntity, err error) {
	query := customerSelectString + ` WHERE def."deleted_at" IS NULL  AND (LOWER (def.email) = $1 OR phone_number = $2) ORDER BY def."created_at" DESC  LIMIT 1`
	row := model.DB.QueryRow(query, strings.ToLower(email), phoneNumber)
	res, err = model.scanRow(row)

	return res, err
}

// FindByID ...
func (model customerModel) FindByID(id string) (res CustomerEntity, err error) {
	query := customerSelectString + ` WHERE def."deleted_at" IS NULL  AND customer_id = $1 ORDER BY def."created_at" DESC  LIMIT 1`
	row := model.DB.QueryRow(query, id)
	res, err = model.scanRow(row)

	return res, err
}

// Store ...
func (model customerModel) Store(body viewmodel.CustomerVM, changedAt time.Time) (res string, err error) {
	sql := `INSERT INTO "customers" (
		"customer_name", "email", "phone_number", "dob", "sex", "password", "created_at", "updated_at"
		) VALUES($1, $2, $3, $4, $5, $6, $7, $7) RETURNING "customer_id"`
	err = model.DB.QueryRow(sql, body.CustomerName, body.Email, body.PhoneNumber, body.Dob, body.Sex, body.Password, changedAt).Scan(&res)

	return res, err
}
