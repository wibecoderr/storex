package dbhelper

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/wibecoderr/storex/database"
	"github.com/wibecoderr/storex/model"
)

/*
CreateAsset(...)       // insert into assets table -- Done
GetAssetByID(id)       // fetch one asset -- Done
ListAssets(filters...) // list with optional filter by type/status - Done
AssignAsset(assetID, empID) // set emp_id + status = 'assigned' -- Done
ReturnAsset(assetID)        // set emp_id = null + status = 'available'
ArchiveAsset(id)            // soft delete
*/
func CreateAsset(tx *sqlx.Tx, request model.CreateAssetRequest) (string, error) {

	sql := `INSERT INTO assets (brand, model, serial_no, type, owner, purchased_at, warranty_start, warranty_end, note)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            RETURNING id`
	var id string
	err := tx.Get(&id, sql, request.Brand, request.Model, request.Serial, request.Type, request.Owner, request.PurchasedAt, request.WarrantyStart, request.WarrantyEnd, request.Note)
	return id, err
}

func CreateLaptop(tx *sqlx.Tx, assetID, processor string, ram, storage int, os, charger string) error {
	sql := `INSERT INTO laptop (asset_id, processor, ram, storage, os, charger)
            VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := tx.Exec(sql, assetID, processor, ram, storage, os, charger)
	return err
}
func CreateKeyBoard(tx *sqlx.Tx, assestId, layout string) error {
	sql := ` insert into keyboard (asset_id, layout) values ($1, $2)`
	_, err := tx.Exec(sql, assestId, layout)
	return err
}
func CreateMouse(tx *sqlx.Tx, assestid string, dpi int, wireless bool) error {
	sql := `insert into mouse ( asset_id , dpi, is_wireless)
VALUES ($1, $2, $3)`
	_, err := tx.Exec(sql, assestid, dpi, wireless)
	return err
}
func CreateMobile(tx *sqlx.Tx, assetID, os string, ram, storage int, charger string) error {
	sql := `INSERT INTO mobile (asset_id, os, ram, storage, charger)
            VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(sql, assetID, os, ram, storage, charger)
	return err
}
func CreateHardware(tx *sqlx.Tx, assestId string, storage int) error {
	sql := `insert into hardware ( asset_id , storage)
values ($1, $2)`
	_, err := tx.Exec(sql, assestId, storage)
	return err
}

func GetAssetByID(id string) (model.AssetDetail, error) {
	sql := `SELECT id, brand, model, serial_no, type, status, owner, purchased_at, warranty_start, warranty_end, note, archived_at 
            FROM assets 
            WHERE id = $1 AND archived_at IS NULL`
	var details model.AssetDetail
	err := database.DB.Get(&details.Asset, sql, id)
	if err != nil {
		return details, err
	}

	// todo use making of const - remanin, omitempty -done

	switch details.Asset.Type {
	case "laptop":
		var laptop model.Laptop
		err = database.DB.Get(&laptop, `SELECT processor, ram, storage, os, charger FROM laptop WHERE asset_id = $1`, id)
		details.Laptop = &laptop
	case "mobile":
		var mobile model.Mobile
		err = database.DB.Get(&mobile, `SELECT os, ram, storage, charger FROM mobile WHERE asset_id = $1`, id)
		details.Mobile = &mobile
	case "keyboard":
		var keyboard model.Keyboard
		err = database.DB.Get(&keyboard, `SELECT layout FROM keyboard WHERE asset_id = $1`, id)
		details.Keyboard = &keyboard
	case "hardware":
		var hardware model.Hardware
		err = database.DB.Get(&hardware, `SELECT storage FROM hardware WHERE asset_id = $1`, id)
		details.Hardware = &hardware
	case "mouse":
		var mouse model.Mouse
		err = database.DB.Get(&mouse, `SELECT dpi, is_wireless FROM mouse WHERE asset_id = $1`, id)
		details.Mouse = &mouse
	}

	return details, err
}

func ListAssets(limit, offset int) ([]model.DisplayAssetResponse, error) {
	sql := `SELECT a.brand, a.model, a.serial_no, a.type, a.status,
                   e.name as employee_name, e.id as employee_id
            FROM assets a
            LEFT JOIN employee e ON e.id = a.emp_id
            WHERE a.archived_at IS NULL
            ORDER BY a.type
            LIMIT $1 OFFSET $2`
	var devices []model.DisplayAssetResponse
	err := database.DB.Select(&devices, sql, limit, offset)
	return devices, err
}

func AssignAsset(assetID, empID string) error {
	return database.Tx(func(tx *sqlx.Tx) error {

		sql := `SELECT status FROM assets WHERE id = $1 AND archived_at IS NULL`
		var status string
		err := tx.Get(&status, sql, assetID)
		if err != nil {
			return fmt.Errorf("asset not found or archived: %w", err)
		}
		if status == "assigned" {
			return fmt.Errorf("asset %s already assigned to %s", assetID, empID)
		}

		sql = `UPDATE assets SET emp_id = $1, status = 'assigned' WHERE id = $2 AND archived_at IS NULL`
		_, err = tx.Exec(sql, empID, assetID)
		if err != nil {
			return err
		}

		sql = `INSERT INTO asset_history (type, asset_id, assigned_to, assigned_on) VALUES ('assigned', $1, $2, now())`
		_, err = tx.Exec(sql, assetID, empID)
		return err
	})
}

func ReturnAsset(assetID, note string) error {
	return database.Tx(func(tx *sqlx.Tx) error {

		// check whom assigned
		var empID string
		err := tx.Get(&empID, `SELECT emp_id FROM assets WHERE id = $1 AND archived_at IS NULL AND status = 'assigned'`, assetID)
		if err != nil {
			return fmt.Errorf("asset not found or not currently assigned: %w", err)
		}

		_, err = tx.Exec(`UPDATE assets SET emp_id = NULL, status = 'available' WHERE id = $1 AND archived_at IS NULL`, assetID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO asset_history (type, asset_id, assigned_to, assigned_on, returned_on, return_status)
			VALUES ('available', $1, $2, now(), now(), $3)`,
			assetID, empID, note)
		return err
	})
}

func RemoveAsset(assetID string) error {
	//  todo checked if remove only when procuct is not assigned
	sql := `UPDATE assets SET archived_at = now() WHERE id = $1 AND archived_at IS NULL and status !='assigned'`
	_, err := database.DB.Exec(sql, assetID)
	return err
}
func ListAssetsByEmployee(empID string) ([]model.Asset, int, error) {
	var assets []model.Asset
	err := database.DB.Select(&assets, `SELECT * FROM assets WHERE emp_id = $1 AND archived_at IS NULL ORDER BY type`, empID)
	if err != nil {
		return nil, 0, err
	}
	return assets, len(assets), nil
}

func CheckStatus(assestId string) bool {
	sql := `select status from assets where id = $1 and   archived_at IS NULL`
	var status string
	err := database.DB.Get(&status, sql, assestId)
	if err != nil {
		return false
	}
	return status != "assigned"
}
