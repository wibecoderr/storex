package model

import "time"

type Employee struct {
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	Password     string    `db:"password"`
	Role         string    `db:"role"`
	PhoneNo      string    `db:"phone_no"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Getuser struct {
	ID            string `db:"id"`
	Password_hash string `db:"password_hash"`
}

type Asset struct {
	ID            string     `db:"id"`
	Brand         string     `db:"brand"`
	Model         string     `db:"model"`
	SerialNo      string     `db:"serial_no"`
	Type          string     `db:"type"`
	Status        string     `db:"status"`
	PurchasedAt   time.Time  `db:"purchased_at"`
	WarrantyStart *time.Time `db:"warranty_start"` // add this
	WarrantyEnd   *time.Time `db:"warranty_end"`
	Owner         string     `db:"owner"`
	ArchivedAt    *time.Time `db:"archived_at"`
	Note          *string    `db:"note"`
}
type AssetDetail struct {
	Asset    Asset
	Laptop   *Laptop
	Mouse    *Mouse
	Keyboard *Keyboard
	Mobile   *Mobile
	Hardware *Hardware
}

type Laptop struct {
	Processor string `db:"processor"`
	Ram       int    `db:"ram"`
	Storage   int    `db:"storage"`
	Os        string `db:"os"`
	Charger   string `db:"charger"`
}

type Mouse struct {
	Dpi        int  `db:"dpi"`
	IsWireless bool `db:"is_wireless"`
}

type Mobile struct {
	Os      string `db:"os"`
	Ram     int    `db:"ram"`
	Storage int    `db:"storage"`
	Charger string `db:"charger"`
}

type Hardware struct {
	Storage int `db:"storage"`
}

type Keyboard struct {
	Layout string `db:"layout"`
}

type UserCxt struct {
	UserId    string `json:"user_id"`
	SessionId string `json:"session_id"`
}
type Error struct {
	Error      string
	StatusCode int
	Message    string
}
type Register struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role"`
	PhoneNo  string `json:"phone_no" validate:"required,min=10,max=10"`
	Password string `json:"password" validate:"required"`
}
type Login struct {
	Email    string `db:"email" validate:"required,email"`
	Password string `db:"password_hash" validate:"required"`
}

type CreateAssetRequest struct {
	Brand         string     `json:"brand" validate:"required"`
	Model         string     `json:"model" validate:"required"`
	Serial        string     `json:"serial" validate:"required"`
	Type          string     `json:"type" validate:"required"`
	Owner         string     `json:"owner" validate:"required"`
	PurchasedAt   time.Time  `json:"purchased_at" validate:"required"`
	WarrantyStart *time.Time `json:"warranty_start"`
	WarrantyEnd   *time.Time `json:"warranty_end"`
	Note          *string    `json:"note"`
	Laptop        *Laptop    `json:"laptop"`
	Mouse         *Mouse     `json:"mouse"`
	Keyboard      *Keyboard  `json:"keyboard"`
	Mobile        *Mobile    `json:"mobile"`
	Hardware      *Hardware  `json:"hardware"`
}
type AssignAssetRequest struct {
	AssetID string `json:"asset_id" validate:"required"`
	EmpID   string `json:"emp_id" validate:"required"`
}
type DeleteAssetRequest struct {
	AssetID string `json:"asset_id" validate:"required"`
}
type DisplayRequest struct {
	Model string `json:"model" validate:"required"`
	TYpe  string `json:"type" validate:"required"`
	empId string `json:"emp_id" validate:"required"`
	owner string `json:"owner" validate:"required"`
}
type DisplayAssetResponse struct {
	Brand    string  `db:"brand" json:"brand"`
	Model    string  `db:"model" json:"model"`
	SerialNo string  `db:"serial_no" json:"serial_no"`
	Type     string  `db:"type" json:"type"`
	Status   string  `db:"status" json:"status"`
	EmpName  *string `db:"employee_name" json:"employee_name"`
	EmpID    *string `db:"employee_id" json:"employee_id"`
}
