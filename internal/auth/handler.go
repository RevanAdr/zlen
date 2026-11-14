package auth

import "net/http"

type UserHandler interface {
	List(w http.ResponseWriter, r *http.Request)
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
