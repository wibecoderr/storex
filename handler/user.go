package handler

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/wibecoderr/storex"
	"github.com/wibecoderr/storex/database"
	"github.com/wibecoderr/storex/database/dbhelper"
	"github.com/wibecoderr/storex/middleware"
	"github.com/wibecoderr/storex/model"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.Register
	err := utils.ParseBody(r.Body, &user)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Fail to parse body")
		return
	}
	if errs := utils.ValidateStruct(user); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}
	// check if user exist
	exists, err := dbhelper.UserExist(user.Email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Fail to create user")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusConflict, nil, "User already exists")
		return
	}

	hashpassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "Fail to hash password")
		return
	}
	// add employee
	var (
		jwtToken  string
		empID     string
		sessionId string
	)
	//  err = database2.Tx(func(tx *sqlx.Tx) error {
	err = database.Tx(func(tx *sqlx.Tx) error {
		empID, err = dbhelper.AddEmployee(tx, user.Name, user.Email, user.Role, user.PhoneNo, hashpassword)
		if err != nil {
			return err
		}

		sessionId, err = dbhelper.CreateSession(tx, empID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Fail to create user")
		return
	}
	jwtToken, err = utils.GenerateJWT(empID, sessionId)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Fail to create user")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, map[string]string{"token": jwtToken})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user model.Login
	if err := utils.ParseBody(r.Body, &user); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to parse body")
		return
	}
	if errs := utils.ValidateStruct(user); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	emp, err := dbhelper.GetEmployeeByEmail(user.Email)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "Invalid credentials")
		return
	}
	if !utils.CheckPasswordHash(user.Password, emp.Password_hash) {
		utils.RespondError(w, http.StatusUnauthorized, nil, "Invalid credentials")
		return
	}
	var jwtToken string
	err = database.Tx(func(tx *sqlx.Tx) error {
		sessionID, err := dbhelper.CreateSession(tx, emp.ID)
		if err != nil {
			return err
		}
		jwtToken, err = utils.GenerateJWT(emp.ID, sessionID)
		return err
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to login")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"token": jwtToken})
}
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	err := dbhelper.LogoutSession(user.SessionId)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to logout")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}
