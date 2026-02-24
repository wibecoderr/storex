package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/wibecoderr/storex"
	"github.com/wibecoderr/storex/database"
	"github.com/wibecoderr/storex/database/dbhelper"
	"github.com/wibecoderr/storex/model"
)

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var (
		device  model.CreateAssetRequest
		assetID string
	)

	if err := utils.ParseBody(r.Body, &device); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}

	if errs := utils.ValidateStruct(device); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	} // tx
	err := database.Tx(func(tx *sqlx.Tx) error {

		assetID, err := dbhelper.CreateAsset(tx, device.Brand, device.Model, device.Serial, device.Type, device.Owner, device.PurchasedAt, device.WarrantyStart, device.WarrantyEnd, device.Note)
		if err != nil {
			return err
		}

		switch device.Type {
		case "laptop":
			err = dbhelper.CreateLaptop(tx, assetID, device.Laptop.Processor, device.Laptop.Os, device.Laptop.Charger, device.Laptop.Ram, device.Laptop.Storage)
		case "mouse":
			err = dbhelper.CreateMouse(tx, assetID, device.Mouse.Dpi, device.Mouse.IsWireless)
		case "keyboard":
			err = dbhelper.CreateKeyBoard(tx, assetID, device.Keyboard.Layout)
		case "mobile":
			err = dbhelper.CreateMobile(tx, assetID, device.Mobile.Os, device.Mobile.Ram, device.Mobile.Storage, device.Mobile.Charger)
		case "hardware":
			err = dbhelper.CreateHardware(tx, assetID, device.Hardware.Storage)
		}
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to create asset details")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, map[string]string{"id": assetID})
}

func AssignAsset(w http.ResponseWriter, r *http.Request) {
	var req model.AssignAssetRequest
	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to parse body")
		return
	}
	if errs := utils.ValidateStruct(req); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	err := dbhelper.AssignAsset(req.AssetID, req.EmpID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to assign asset")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "asset assigned successfully"})
}

func DeleteAsset(w http.ResponseWriter, r *http.Request) {
	var device model.DeleteAssetRequest
	if err := utils.ParseBody(r.Body, &device); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to parse body")
		return
	}
	if errs := utils.ValidateStruct(device); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	err := dbhelper.RemoveAsset(device.AssetID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to remove asset")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "asset removed successfully"})
}

func ListAssetsByEmployee(w http.ResponseWriter, r *http.Request) {
	empID := chi.URLParam(r, "id")

	assets, count, err := dbhelper.ListAssetsByEmployee(empID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to get assets")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"assets": assets,
		"count":  count,
	})
}

func GetAssetByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "Missing asset id")
		return
	}

	asset, err := dbhelper.GetAssetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to get asset")
		return
	}

	utils.RespondJSON(w, http.StatusOK, asset)
}

// get assests all , get dashboard ,

func DisplayAsset(w http.ResponseWriter, r *http.Request) {
	limit := 10
	page := 1
	offset := (page - 1) * limit
	assetId, err := dbhelper.ListAssets(limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to get assets")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"assets": assetId,
	})

}
