package handler

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/jmoiron/sqlx"
	utils "github.com/wibecoderr/storex"
	"github.com/wibecoderr/storex/database"
	"github.com/wibecoderr/storex/database/dbhelper"
	"github.com/wibecoderr/storex/model"
	"google.golang.org/api/option"
)

var firebaseAuth *auth.Client

func InitFirebase() {
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("failed to init firebase: " + err.Error())
	}
	firebaseAuth, err = app.Auth(context.Background())
	if err != nil {
		panic("failed to init firebase auth: " + err.Error())
	}
}

func FirebaseLogin(w http.ResponseWriter, r *http.Request) {
	var req model.FirebaseLoginRequest
	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}

	// verify firebase token
	token, err := firebaseAuth.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "invalid firebase token")
		return
	}

	email, _ := token.Claims["email"].(string)
	name, _ := token.Claims["name"].(string)

	// check if user exists
	exists, err := dbhelper.UserExist(email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "database error")
		return
	}

	var empID string
	if !exists {
		// new user - auto register as employee
		err = database.Tx(func(tx *sqlx.Tx) error {
			empID, err = dbhelper.AddEmployee(tx, name, email, "employee", "0000000000", "")
			return err
		})
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to create user")
			return
		}
	} else {
		// existing user - get their ID
		empID, err = dbhelper.GetEmployeeIDByEmail(email)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "failed to get user")
			return
		}
	}

	// create session + JWT (same as your existing login)
	var jwtToken string
	err = database.Tx(func(tx *sqlx.Tx) error {
		sessionID, err := dbhelper.CreateSession(tx, empID)
		if err != nil {
			return err
		}
		jwtToken, err = utils.GenerateJWT(empID, sessionID)
		return err
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create session")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"token": jwtToken})
}
