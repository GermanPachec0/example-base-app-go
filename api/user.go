package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/GermanPachec0/app-go/domain"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CreateUserRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

func (a *APIServer) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user, err := a.userRepo.GetAll(ctx)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}
	WriteJson(w, 200, user)
}

func (a *APIServer) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	uidStr := mux.Vars(r)["uid"]
	if uidStr == "" {
		a.errorResponse(w, r, 500, errors.New("no params found"))
		return
	}
	uid, err := uuid.ParseBytes([]byte(uidStr))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}
	user, err := a.userRepo.GetById(ctx, uid)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}
	WriteJson(w, 200, user)
}

func (a *APIServer) createUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	req := CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}
	user, err := domain.NewUser(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	if err := a.userRepo.Create(ctx, user); err != nil {
		fmt.Println(err.Error())
		a.errorResponse(w, r, 500, err)
		return
	}
	WriteJson(w, http.StatusCreated, user)
}
