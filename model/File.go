package model

import (
	"kriyapeople/usecase/viewmodel"
	"time"

	"database/sql"
)

var (
	// FileAdminProfile ...
	FileAdminProfile = "admin_profile"
	// FileUserProfile ...
	FileUserProfile = "user_profile"
	// FileWhitelist ...
	FileWhitelist = []string{FileAdminProfile, FileUserProfile}

	fileSelectString = `SELECT f."id", f."type", f."url", f."user_upload", f."created_at", f."updated_at",
	f."deleted_at" FROM "files" f
	LEFT JOIN "admins" admins ON admins."profile_image_id" = f."id"
	LEFT JOIN "users" users ON users."profile_image_id" = f."id"`
	unassignedQueryString = `AND admins."id" IS NULL AND users."id" IS NULL`
)

func (model fileModel) scanRows(rows *sql.Rows) (d FileEntity, err error) {
	err = rows.Scan(
		&d.ID, &d.Type, &d.URL, &d.UserUpload, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

func (model fileModel) scanRow(row *sql.Row) (d FileEntity, err error) {
	err = row.Scan(
		&d.ID, &d.Type, &d.URL, &d.UserUpload, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)

	return d, err
}

// IFile ...
type IFile interface {
	FindAllUnassignedByUserID(userUpload, types string) (data []FileEntity, err error)
	FindByID(id string) (FileEntity, error)
	FindUnassignedByID(id, types, userUpload string) (FileEntity, error)
	Store(body viewmodel.FileVM, changedAt time.Time) (string, error)
	Destroy(id string, changedAt time.Time) (string, error)
}

// FileEntity ....
type FileEntity struct {
	ID         string         `db:"id"`
	Type       sql.NullString `db:"type"`
	URL        sql.NullString `db:"url"`
	UserUpload sql.NullString `db:"user_upload"`
	CreatedAt  string         `db:"created_at"`
	UpdatedAt  string         `db:"updated_at"`
	DeletedAt  sql.NullString `db:"deleted_at"`
}

// fileModel ...
type fileModel struct {
	DB *sql.DB
}

// NewFileModel ...
func NewFileModel(db *sql.DB) IFile {
	return &fileModel{DB: db}
}

// FindAllUnassignedByUserID ...
func (model fileModel) FindAllUnassignedByUserID(userUpload, types string) (res []FileEntity, err error) {
	query := fileSelectString + ` WHERE f."deleted_at" IS NULL AND f."user_upload" = $1 AND f."type" = $2
		` + unassignedQueryString + ` ORDER BY f."created_at"`

	rows, err := model.DB.Query(query, userUpload, types)
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
func (model fileModel) FindByID(id string) (res FileEntity, err error) {
	query := `SELECT "id", "type", "url", "user_upload", "created_at", "updated_at", "deleted_at"
		FROM "files" WHERE "deleted_at" IS NULL AND "id" = $1
		ORDER BY "created_at" DESC LIMIT 1`
	row := model.DB.QueryRow(query, id)
	res, err = model.scanRow(row)

	return res, err
}

// FindUnassignedByID ...
func (model fileModel) FindUnassignedByID(id, types, userUpload string) (res FileEntity, err error) {
	query := fileSelectString + ` WHERE f."deleted_at" IS NULL AND f."id" = $1 AND f."type" = $2
		AND f."user_upload" = $3 ` + unassignedQueryString + ` ORDER BY f."created_at" DESC LIMIT 1`
	row := model.DB.QueryRow(query, id, types, userUpload)
	res, err = model.scanRow(row)

	return res, err
}

// Store ...
func (model fileModel) Store(body viewmodel.FileVM, changedAt time.Time) (res string, err error) {
	sql :=
		`INSERT INTO "files" ("type", "url", "user_upload", "created_at", "updated_at")
		VALUES($1, $2, $3, $4, $4) RETURNING "id"`
	err = model.DB.QueryRow(sql, body.Type, body.URL, body.UserUpload, changedAt).Scan(&res)

	return res, err
}

// Destroy ...
func (model fileModel) Destroy(id string, changedAt time.Time) (res string, err error) {
	sql := `UPDATE "files" SET deleted_at = $1 WHERE id = $2 RETURNING "id"`
	err = model.DB.QueryRow(sql, changedAt, id).Scan(&res)

	return res, err
}
