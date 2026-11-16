package auth

import (
	"errors"
	"net/http"
	"strconv"
)

type UserHandler interface {
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
}

type UserHandlerImpl struct {
	service UserService
}

func NewUserHandler(service UserService) UserHandler {
	return &UserHandlerImpl{service: service}
}

func (h *UserHandlerImpl) List(w http.ResponseWriter, r *http.Request) {
	responses, err := h.service.List(r.Context())
	if err != nil {
		WriteToJSON(w, http.StatusInternalServerError, "Internal server error", map[string]string{"error": string(err.Error())})
		return
	}
	WriteToJSON(w, http.StatusOK, "OK", responses)
}

func (h *UserHandlerImpl) Get(w http.ResponseWriter, r *http.Request) {
	id, cerr := strconv.Atoi(r.PathValue("id"))
	if cerr != nil {
		WriteToJSON(w, http.StatusBadRequest, "Bad Request", map[string]string{"error": string(cerr.Error())})
		return
	}
	response, err := h.service.Get(r.Context(), id)
	switch {
	case errors.Is(ErrNotFound, err):
		WriteToJSON(w, http.StatusNotFound, "Not Found", map[string]string{"error": string(err.Error())})
		return
	case err != nil:
		WriteToJSON(w, http.StatusInternalServerError, "Internal Server Error", map[string]string{"error": "Internal Server Error"})
		return
	default:
		WriteToJSON(w, http.StatusOK, "OK", response)
	}
}

func (h *UserHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	request := User{}
	if rerr := ReadBody(r, &request); rerr != nil {
		WriteToJSON(w, http.StatusInternalServerError, "Internal Server Error", map[string]string{"error": string(rerr.Error())})
		return
	}

	data := User{
		Name: request.Name,
		Role: request.Role,
	}
	response, err := h.service.Create(r.Context(), data)
	if err != nil {
		WriteToJSON(w, http.StatusInternalServerError, "Internal Server Error", map[string]string{"error": string(err.Error())})
		return
	}

	WriteToJSON(w, http.StatusCreated, "Created", response)
}
