package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/wibecoderr/storex"
	"github.com/wibecoderr/storex/database"
	"github.com/wibecoderr/storex/database/dbhelper"
	"github.com/wibecoderr/storex/middleware"
	"github.com/wibecoderr/storex/model"
)

// this is getting too big - split later maybe? idk
// todo pass whole model in db call so dont have to write too many things - annoying rn
func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var (
		device  model.CreateAssetRequest
		assetID string
		err     error
	)

	if err := utils.ParseBody(r.Body, &device); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}

	if errs := utils.ValidateStruct(device); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	// using Istype() here bc we dont unneccessary typing for type

	if !model.Device(device.Type).Istype() {
		utils.RespondError(w, http.StatusInternalServerError, nil, "wrong type for device ")
		return
	}

	err = database.Tx(func(tx *sqlx.Tx) error {
		assetID, err = dbhelper.CreateAsset(tx, device)
		if err != nil {
			return err
		}

		switch device.Type {
		case "laptop":
			err = dbhelper.CreateLaptop(tx, assetID, device.Laptop.Processor, device.Laptop.Ram, device.Laptop.Storage, device.Laptop.Os, device.Laptop.Charger)
		case "mouse":
			err = dbhelper.CreateMouse(tx, assetID, device.Mouse.Dpi, device.Mouse.IsWireless)
		case "keyboard":
			err = dbhelper.CreateKeyBoard(tx, assetID, device.Keyboard.Layout)
		case "mobile":
			err = dbhelper.CreateMobile(tx, assetID, device.Mobile.Os, device.Mobile.Ram, device.Mobile.Storage, device.Mobile.Charger)
		case "hardware":
			err = dbhelper.CreateHardware(tx, assetID, device.Hardware.Storage)
			// if we add more device types later dont forget to add here too - happened once already
		}
		return err
	})

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create asset details")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, map[string]string{"id": assetID})
}

func AssignAsset(w http.ResponseWriter, r *http.Request) {
	var req model.AssignAssetRequest
	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	if errs := utils.ValidateStruct(req); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	// CheckStatus returns false if already assigned
	// todo - make this less confusing, the naming is backwards feeling
	if !dbhelper.CheckStatus(req.AssetID) {
		utils.RespondError(w, http.StatusBadRequest, nil, "assest already assigned to someone ")
		return
	}

	if err := dbhelper.AssignAsset(req.AssetID, req.EmpID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to assign asset")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "asset assigned successfully"})
}

// soft delete basically
// wont work if asset is currently assigned, which is intentional fo stoppong dupicay of rows in tabl e
func DeleteAsset(w http.ResponseWriter, r *http.Request) {
	var device model.DeleteAssetRequest
	if err := utils.ParseBody(r.Body, &device); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	if errs := utils.ValidateStruct(device); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	// todo := check asset is already assigned or not -- Done , do testing on postman
	if !dbhelper.CheckStatus(device.AssetID) {
		utils.RespondError(w, http.StatusBadRequest, nil, "assest status is assgined mark it available first ")
		return
	}

	if err := dbhelper.RemoveAsset(device.AssetID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to remove asset")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "asset removed successfully"})
}

// for the employee themselves to see their own stuff which they are using
func ListAssetsByEmployee(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	assets, count, err := dbhelper.ListAssetsByEmployee(user.UserId)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get assets")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"assets": assets,
		"count":  count,
	})
}

// same as above but admin can pass any empID via url param
// probably should merge these two at some point but works for now
func ListAssetsByEmployeeAdmin(w http.ResponseWriter, r *http.Request) {
	empID := chi.URLParam(r, "id")

	assets, count, err := dbhelper.ListAssetsByEmployee(empID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get assets")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{"assets": assets, "count": count})
}

func GetAssetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "missing asset id")
		return
	}

	asset, err := dbhelper.GetAssetByID(id)
	if err != nil {
		// could be sql.ErrNoRows but just returning 500 for now, fix later
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get asset")
		return
	}
	utils.RespondJSON(w, http.StatusOK, asset)
}

// get assests all , get dashboard ,
func DisplayAsset(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")

	l, err := strconv.Atoi(limit)
	if err != nil || l > 100 {
		l = 10
	}

	p, err := strconv.Atoi(page)
	if err != nil || p < 1 {
		p = 1
	}

	offset := (p - 1) * l

	Devices, err := dbhelper.ListAssets(l, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get assets")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"assets": Devices,
	})
}

// note: typo in function name is intentional to match route registration - dont fix
func ReturnAssest(w http.ResponseWriter, r *http.Request) {
	var device model.ReturnRequest

	assestId := chi.URLParam(r, "id")
	if assestId == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "missing asset id")
		return
	}

	if err := utils.ParseBody(r.Body, &device); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	if errs := utils.ValidateStruct(device); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	// note field is required - tells us why its being returned, goes into asset_history
	if err := dbhelper.ReturnAsset(assestId, device.Note); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to return asset")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{"message": "asset returned successfully"})
}
