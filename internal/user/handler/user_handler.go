package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/danilobml/user-manager/internal/helpers"
	"github.com/danilobml/user-manager/internal/user/dtos"
	"github.com/danilobml/user-manager/internal/user/service"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
	apiKey      string
}

func NewUserHandler(userService service.UserService, apiKey string) *UserHandler {
	return &UserHandler{
		userService: userService,
		apiKey:      apiKey,
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

	if !uh.isInputValid(w, registerReq) {
		return
	}

	registerReq.Password = strings.TrimSpace(registerReq.Password)
	registerReq.Email = strings.TrimSpace(registerReq.Email)

	resp, err := uh.userService.Register(ctx, registerReq)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusCreated, resp)
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

	if !uh.isInputValid(w, loginReq) {
		return
	}

	loginReq.Password = strings.TrimSpace(loginReq.Password)
	loginReq.Email = strings.TrimSpace(loginReq.Email)

	resp, err := uh.userService.Login(ctx, loginReq)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, resp)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "no valid user id supplied")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	updateReq := dtos.UpdateUserRequest{}
	err = json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !uh.isInputValid(w, updateReq) {
		return
	}

	updateReq.ID = userId
	updateReq.Email = strings.TrimSpace(updateReq.Email)

	err = uh.userService.UpdateUserData(ctx, updateReq)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, "updated successfully")
}

func (uh *UserHandler) UnregisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "no valid user id supplied")
		return
	}

	user, err := uh.userService.GetUser(ctx, userId)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	err = uh.userService.Unregister(ctx, dtos.UnregisterRequest{Email: user.Email})
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusNoContent, "unregistered")
}

func (uh *UserHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	requestPassResetReq := dtos.RequestPasswordResetRequest{}
	err := json.NewDecoder(r.Body).Decode(&requestPassResetReq)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !uh.isInputValid(w, requestPassResetReq) {
		return
	}

	requestPassResetReq.Email = strings.TrimSpace(requestPassResetReq.Email)

	err = uh.userService.RequestPasswordReset(ctx, requestPassResetReq)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusNoContent, "")
}

func (uh *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	resetPassReq := dtos.ResetPasswordRequest{}
	err := json.NewDecoder(r.Body).Decode(&resetPassReq)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !uh.isInputValid(w, resetPassReq) {
		return
	}

	resetPassReq.Password = strings.TrimSpace(resetPassReq.Password)
	resetPassReq.Email = strings.TrimSpace(resetPassReq.Email)

	err = uh.userService.ResetPassword(ctx, resetPassReq)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusNoContent, "")
}

func (uh *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := uh.userService.ListAllUsers(ctx)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, users)
}

func (uh *UserHandler) RemoveUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")
	userId, err := uuid.Parse(idString)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "no valid user id supplied")
		return
	}

	err = uh.userService.RemoveUser(ctx, userId)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusNoContent, "removed")
}

func (uh *UserHandler) CheckUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	apiKey := strings.TrimSpace(r.Header.Get("User-Api-Key"))
	if apiKey != uh.apiKey {
		helpers.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	checkUserReq := dtos.CheckUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&checkUserReq)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !uh.isInputValid(w, checkUserReq) {
		return
	}

	resp, err := uh.userService.CheckUser(ctx, checkUserReq)
	if err != nil {
		helpers.WriteErrorsResponse(w, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, resp)
}

// Validation Helper:
func (uh *UserHandler) isInputValid(w http.ResponseWriter, structToValidate any) bool {
	validate := validator.New()
	err := validate.Struct(structToValidate)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		helpers.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", errors))
		return false
	}

	return true
}
