package model

import (
	"database/sql"
	"kriyapeople/usecase/viewmodel"
	"strings"
	"time"
)

var (
	// DefaultAdminBy ...
	DefaultAdminBy = "def.updated_at"
	// AdminBy ...
	AdminBy = []string{
		"def.created_at", "def.updated_at",
	}

	adminSelectString = `SELECT def.id, def."data" ->> 'email' as email, def."data" ->> 'password' as password, def."data" ->> 'username' as username, def."data" -> 'status' ->> 'is_active' as status, r."id" as role_id, r."data" ->> 'role_name' as role_name, def.created_at, def.updated_at, def.deleted_at FROM "users" def LEFT JOIN "roles" r ON r."id" = def."role_id"`
)

func (model adminModel) scanRows(rows *sql.Rows) (d UserEntity, err error) {
	err = rows.Scan(
		&d.ID, &d.Email, &d.Password, &d.UserName, &d.Status, &d.RoleID, &d.Role.Name, &d.CreatedAt,
		&d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

func (model adminModel) scanRow(row *sql.Row) (d UserEntity, err error) {
	err = row.Scan(
		&d.ID, &d.Email, &d.Password, &d.UserName, &d.Status, &d.RoleID, &d.Role.Name, &d.CreatedAt,
		&d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

// adminModel ...
type adminModel struct {
	DB *sql.DB
}

// IAdmin ...
type IAdmin interface {
	FindAll(search string, offset, limit int, by, sort string) ([]UserEntity, int, error)
	FindByID(id string) (UserEntity, error)
	FindByEmail(email string) (UserEntity, error)
	Store(body viewmodel.UserVM, changedAt time.Time) (string, error)
	Update(id string, body viewmodel.UserVM, changedAt time.Time) (string, error)
	Destroy(id string, changedAt time.Time) (string, error)
}

// UserEntity ....
type UserEntity struct {
	ID        string         `db:"id"`
	Email     sql.NullString `db:"email"`
	Password  sql.NullString `db:"password"`
	UserName  sql.NullString `db:"user_name"`
	RoleID    sql.NullString `db:"role_id"`
	Role      RoleEntity     `db:"role"`
	Status    sql.NullBool   `db:"status"`
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at"`
}

// NewAdminModel ...
func NewAdminModel(db *sql.DB) IAdmin {
	return &adminModel{DB: db}
}

// FindAll ...
func (model adminModel) FindAll(search string, offset, limit int, by, sort string) (res []UserEntity, count int, err error) {
	query := adminSelectString + ` WHERE def."deleted_at" IS NULL AND (
	LOWER (def."data" ->> 'email' ) LIKE $1 
	OR LOWER ( def."data" ->> 'username' ) LIKE $1 
	) ORDER BY ` + by + ` ` + sort + ` OFFSET $2 LIMIT $3`
	rows, err := model.DB.Query(query, `%`+strings.ToLower(search)+`%`, offset, limit)
	if err != nil {
		return res, count, err
	}
	defer rows.Close()

	for rows.Next() {
		d, err := model.scanRows(rows)
		if err != nil {
			return res, count, err
		}
		res = append(res, d)
	}
	err = rows.Err()
	if err != nil {
		return res, count, err
	}

	query = `SELECT COUNT(def."id") FROM "users" def
		LEFT JOIN "roles" r ON r."id" = def."role_id"
		WHERE def."deleted_at" IS NULL AND (
			LOWER (def."data" ->> 'email' ) LIKE $1 
			OR LOWER ( def."data" ->> 'username' ) like $1 )`
	err = model.DB.QueryRow(query, `%`+strings.ToLower(search)+`%`).Scan(&count)

	return res, count, err
}

// FindByID ...
func (model adminModel) FindByID(id string) (res UserEntity, err error) {
	query := adminSelectString + ` WHERE def."deleted_at" IS NULL AND def."id" = $1
		ORDER BY def."created_at" DESC LIMIT 1`
	row := model.DB.QueryRow(query, id)
	res, err = model.scanRow(row)

	return res, err
}

// FindByEmail ...
func (model adminModel) FindByEmail(email string) (res UserEntity, err error) {
	query := adminSelectString + ` WHERE def."deleted_at" IS NULL  AND LOWER (def."data" ->> 'email' ) = $1 ORDER BY def."created_at" DESC  LIMIT 1`
	row := model.DB.QueryRow(query, strings.ToLower(email))
	res, err = model.scanRow(row)

	return res, err
}

// Store ...
func (model adminModel) Store(body viewmodel.UserVM, changedAt time.Time) (res string, err error) {
	sql := `INSERT INTO "users" (
		"data", "role_id", "created_at", "updated_at"
		) VALUES($1, $2, $3, $3) RETURNING "id"`
	err = model.DB.QueryRow(sql, body.Data, body.RoleID, changedAt).Scan(&res)

	return res, err
}

// Update ...
func (model adminModel) Update(id string, body viewmodel.UserVM, changedAt time.Time) (res string, err error) {
	sql := `UPDATE "users" SET "data" = $1, "role_id" = $2, "updated_at" = $3 WHERE "deleted_at" IS NULL
		AND "id" = $4 RETURNING "id"`
	err = model.DB.QueryRow(sql, body.Data, body.RoleID, changedAt, id).Scan(&res)

	return res, err
}

// Destroy ...
func (model adminModel) Destroy(id string, changedAt time.Time) (res string, err error) {
	sql := `UPDATE "users" SET "updated_at" = $1, "deleted_at" = $1
		WHERE "deleted_at" IS NULL AND "id" = $2 RETURNING "id"`
	err = model.DB.QueryRow(sql, changedAt, id).Scan(&res)

	return res, err
}
