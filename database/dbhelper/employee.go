package dbhelper

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/wibecoderr/storex/database"
	"github.com/wibecoderr/storex/model"
)

func AddEmployee(tx *sqlx.Tx, name, email, role, phoneNo, password string) (string, error) {
	// adding employees into table
	sql := `INSERT INTO employee (name, email, role, phone_no , password_hash)
            VALUES ($1, lower(trim($2)), $3, $4, $5)
            RETURNING id`
	var id string
	err := tx.Get(&id, sql, name, email, role, phoneNo, password)
	return id, err
}

func UserExist(email string) (bool, error) {
	// for checking at time of creating admin or user whether user have id or not
	sql := `SELECT count(*) > 0 FROM employee 
            WHERE email = lower(trim($1)) AND archived_at IS NULL`
	var exists bool
	err := database.DB.Get(&exists, sql, email)
	return exists, err
}

func GetEmployeeByEmail(email string) (model.Getuser, error) {
	// for login process
	sql := `select id , password_hash from employee where email = lower(trim($1))`
	var detail model.Getuser
	err := database.DB.Get(&detail, sql, email)
	return detail, err
}

func CreateSession(tx *sqlx.Tx, empID string) (string, error) {
	sql := `INSERT INTO user_sessions (emp_id) VALUES ($1) RETURNING id`
	var sessionID string
	err := tx.Get(&sessionID, sql, empID)
	return sessionID, err
}

func GetUserIDBySession(sessionID string) (string, error) {
	var userID string
	err := database.DB.Get(&userID, `SELECT emp_id FROM user_sessions WHERE id = $1 AND archived_at IS NULL`, sessionID)
	return userID, err
}

func GetEmployeeRole(id string) (string, error) {
	var role string
	err := database.DB.Get(&role, `SELECT role FROM employee WHERE id = $1 AND archived_at IS NULL`, id)
	return role, err
}
func LogoutSession(sessionID string) error {
	sql := `UPDATE user_sessions SET archived_at = now() WHERE id = $1`
	_, err := database.DB.Exec(sql, sessionID)
	return err
}
func CreateEmployee(tx *sqlx.Tx, name, email, role, phoneNo, password string) (string, error) {
	sql := `
INSERT INTO employee (name, email, role, phone_no, password_hash)
VALUES ($1, trim(lower($2)), $3, $4, $5)
RETURNING id
`
	var id string
	err := tx.Get(&id, sql, name, email, role, phoneNo, password)
	return id, err
}
func DeleteSession(sessionID string) error {
	sql := `DELETE FROM user_sessions WHERE emp_id = $1`
	_, err := database.DB.Exec(sql, sessionID)
	return err
}
func ArchiveEmployee(empID string) error {
	return database.Tx(func(tx *sqlx.Tx) error {

		sql := `
			UPDATE employee
			SET archived_at = NOW()
			WHERE id = $1
			AND archived_at IS NULL
		`

		result, err := tx.Exec(sql, empID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`DELETE FROM user_sessions WHERE emp_id = $1`, empID)
		if err != nil {
			return err
		}

		return nil
	})
}
func ListEmployee(Type, status string) ([]model.EmployeeListResponse, error) {
	sql := `SELECT
    e.id,
    e.name,
    e.email,
    e.phone_no,
    e.role,
    COUNT(a.id) AS asset_count
FROM employee e
LEFT JOIN assets a ON a.emp_id = e.id
    AND a.archived_at IS NULL
    AND ($1 = '' OR a.type = $1)
AND ($2 = '' OR a.status = $2)
WHERE e.archived_at IS NULL
GROUP BY e.id, e.name, e.email, e.phone_no, e.role
ORDER BY e.name;`
	var user []model.EmployeeListResponse
	err := database.DB.Select(&user, sql, Type, status)
	return user, err
}
