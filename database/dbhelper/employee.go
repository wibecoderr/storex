package dbhelper

import (
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func AddEmployee(name, email, role, phoneNo, passwordHash string) (string, error) {
	sql := `INSERT INTO employee (name, email, role, phone_no, password_hash)
            VALUES ($1, lower(trim($2)), $3, $4, $5)
            RETURNING id`
	var id string
	err := DB.Get(&id, sql, name, email, role, phoneNo, passwordHash)
	return id, err
}

func GetUserByEmail(email string) (bool, error) {
	sql := `SELECT count(*) > 0 FROM employee 
            WHERE email = lower(trim($1)) AND archived_at IS NULL`
	var exists bool
	err := DB.Get(&exists, sql, email)
	return exists, err
}

func ArchiveEmployee(id string) error {
	sql := ` UPDATE employee
SET archived_at = NOW() WHERE id = $1
and archived_at IS NULL`
	_, err := DB.Exec(sql, id)
	return err
}
