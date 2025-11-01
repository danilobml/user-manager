package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/danilobml/user-manager/internal/errs"
	"github.com/danilobml/user-manager/internal/helpers"
	"github.com/danilobml/user-manager/internal/user/dtos"
	"github.com/danilobml/user-manager/internal/user/service"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	registerReq := dtos.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&registerReq)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	validate := validator.New()
	err = validate.Struct(registerReq)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	registerReq.Password = strings.TrimSpace(registerReq.Password)
	registerReq.Email = strings.TrimSpace(registerReq.Email)

	resp, err := uh.UserService.Register(ctx, registerReq)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}

	helpers.WriteJsonResponse(w, http.StatusCreated, resp)
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	loginReq := dtos.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	validate := validator.New()
	err = validate.Struct(loginReq)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	loginReq.Password = strings.TrimSpace(loginReq.Password)
	loginReq.Email = strings.TrimSpace(loginReq.Email)

	resp, err := uh.UserService.Login(ctx, loginReq)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			helpers.WriteJSONError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, errs.ErrInvalidCredentials) {
			helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}

	helpers.WriteJsonResponse(w, http.StatusOK, resp)
}

func (uh *UserHandler) UnregisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	unregisterReq := dtos.UnregisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&unregisterReq)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	validate := validator.New()
	err = validate.Struct(unregisterReq)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	unregisterReq.Email = strings.TrimSpace(unregisterReq.Email)
	unregisterReq.Token = strings.TrimSpace(unregisterReq.Token)

	err = uh.UserService.Unregister(ctx, unregisterReq)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			helpers.WriteJSONError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, errs.ErrInvalidCredentials) {
			helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}

	helpers.WriteJsonResponse(w, http.StatusNoContent, "unregistered")
}

func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := uh.UserService.ListAllUsers(ctx)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}

	helpers.WriteJsonResponse(w, http.StatusOK, users)
}
