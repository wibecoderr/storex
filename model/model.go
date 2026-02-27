package model

import (
	"time"
)

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
type Role string

const (
	RoleIntern     Role = "intern"
	RoleEmployee   Role = "employee"
	RoleManager    Role = "manager"
	RoleFreelancer Role = "freelancer"
)

type Employee1Request struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required"`
	PhoneNo  string `json:"phone_no"`
}

func (r Role) Iscorrect() bool {
	switch r {
	case RoleIntern, RoleEmployee, RoleManager, RoleFreelancer:
		return true

	}
	return false
}

type Getuser struct {
	ID            string `db:"id"`
	Password_hash string `db:"password_hash"`
}

type AssetRequest struct {
	ID            string     `db:"id"`
	EmpID         *string    `db:"emp_id"`
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
	Asset    AssetRequest
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
type RegisterRequest struct {
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
type Device string

const (
	DeviceLaptop   Device = "laptop"
	DeviceMouse    Device = "mouse"
	DeviceKeyboard Device = "keyboard"
	DeviceMobile   Device = "mobile"
	DeviceHardware Device = "hardware"
)

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
	Laptop        *Laptop    `json:"laptop,omitempty"`
	Mouse         *Mouse     `json:"mouse,omitempty"`
	Keyboard      *Keyboard  `json:"keyboard,omitempty"`
	Mobile        *Mobile    `json:"mobile,omitempty"`
	Hardware      *Hardware  `json:"hardware,omitempty"`
}

func (d Device) Istype() bool {
	switch d {
	case DeviceKeyboard, DeviceLaptop, DeviceHardware, DeviceMouse, DeviceMobile:
		return true
	}
	return false

}

type AssignAssetRequest struct {
	AssetID string `json:"asset_id" validate:"required"`
	EmpID   string `json:"emp_id" validate:"required"`
}
type DeleteAssetRequest struct {
	AssetID string `json:"asset_id" validate:"required"`
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
type ReturnRequest struct {
	EmpID          string    `db:"emp_id" json:"emp_id" validate:"required"`
	Status         string    `db:"status" json:"satus"`
	Type           string    `db:"type" json:"type"`
	AssetId        string    `db:"asset_id" json:"asset_id"`
	ReturnedOn     time.Time `db:"returned_on" json:"returned_on"`
	ReturnedStatus string    `db:"returned_status" json:"returned_status"`
	Note           string    `db:"note" json:"note"`
}

type UpdateAssetRequest struct {
	Brand         string     `json:"brand" validate:"required"`
	Model         string     `json:"model" validate:"required"`
	Serial        string     `json:"serial" validate:"required"`
	Type          string     `json:"type" validate:"required"`
	Owner         string     `json:"owner" validate:"required"`
	PurchasedAt   time.Time  `json:"purchased_at" validate:"required"`
	WarrantyStart *time.Time `json:"warranty_start"`
	WarrantyEnd   *time.Time `json:"warranty_end"`
	Note          *string    `json:"note"`
	Laptop        *Laptop    `json:"laptop,omitempty"`
	Mouse         *Mouse     `json:"mouse,omitempty"`
	Keyboard      *Keyboard  `json:"keyboard,omitempty"`
	Mobile        *Mobile    `json:"mobile,omitempty"`
	Hardware      *Hardware  `json:"hardware,omitempty"`
}
type DashboardCount struct {
	Total            int `db:"total"  json:"total"`
	Available        int `db:"available" json:"available"`
	Assigned         int `db:"assigned"   json:"assigned"`
	InService        int `db:"in_service" json:"in_service"`
	WaitingForRepair int `db:"waiting_for_repair" json:"waiting_for_repair"`
	Damaged          int `db:"damaged"   json:"damaged"`
}
type EmployeeListRequest struct {
	AssetType   *string `json:"asset_type"`
	AssetStatus *string `json:"asset_status"`
}

type EmployeeListResponse struct {
	ID    int    `json:"id"    db:"id"`
	Name  string `json:"name"  db:"name"`
	Email string `json:"email" db:"email"`
	Phone string `json:"phone" db:"phone_no"`
	Role  string `json:"role"  db:"role"`
}
