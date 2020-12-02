package model

import (
	"database/sql"
	"strings"
)

var (
	// RoleCodeSuperadmin ...
	RoleCodeSuperadmin = "Admin"
	// RoleCodeAdmin ...
	RoleCodeAdmin = "Member"

	// DefaultRoleBy ...
	DefaultRoleBy = "def.updated_at"
	// RoleBy ...
	RoleBy = []string{"def.created_at", "def.updated_at", "def.code", "def.name"}

	roleSelectString = `SELECT def."id", def."code", def."name", def."status", def."created_at",
		def."updated_at", def."deleted_at" FROM "roles" def`
)

func (model roleModel) scanRows(rows *sql.Rows) (d RoleEntity, err error) {
	err = rows.Scan(
		&d.ID, &d.Code, &d.Name, &d.Status, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

func (model roleModel) scanRow(row *sql.Row) (d RoleEntity, err error) {
	err = row.Scan(
		&d.ID, &d.Code, &d.Name, &d.Status, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

// roleModel ...
type roleModel struct {
	DB *sql.DB
}

// IRole ...
type IRole interface {
	SelectAll(search, by, sort string) ([]RoleEntity, error)
	FindByID(id string) (RoleEntity, error)
	FindByCode(code string) (RoleEntity, error)
}

// RoleEntity ....
type RoleEntity struct {
	ID        string         `db:"id"`
	Code      sql.NullString `db:"code"`
	Name      sql.NullString `db:"name"`
	Status    sql.NullBool   `db:"status"`
	CreatedAt string         `db:"created_at"`
	UpdatedAt string         `db:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at"`
}

// NewRoleModel ...
func NewRoleModel(db *sql.DB) IRole {
	return &roleModel{DB: db}
}

// SelectAll ...
func (model roleModel) SelectAll(search, by, sort string) (res []RoleEntity, err error) {
	query := roleSelectString + ` WHERE def."deleted_at" IS NULL AND (
		LOWER(def."code") LIKE $1 OR LOWER(def."name") LIKE $1
	) ORDER BY ` + by + ` ` + sort
	rows, err := model.DB.Query(query, `%`+strings.ToLower(search)+`%`)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		d, err := model.scanRows(rows)
		if err != nil {
			return res, err
		}
		res = append(res, d)
	}
	err = rows.Err()

	return res, err
}

// FindByID ...
func (model roleModel) FindByID(id string) (res RoleEntity, err error) {
	query := roleSelectString + ` WHERE def."deleted_at" IS NULL AND def."id" = $1
		ORDER BY def."created_at" DESC LIMIT 1`
	row := model.DB.QueryRow(query, id)
	res, err = model.scanRow(row)

	return res, err
}

// FindByCode ...
func (model roleModel) FindByCode(code string) (res RoleEntity, err error) {
	query := roleSelectString + ` WHERE def."deleted_at" IS NULL AND def."code" = $1
		ORDER BY def."created_at" DESC LIMIT 1`
	row := model.DB.QueryRow(query, code)
	res, err = model.scanRow(row)

	return res, err
}
