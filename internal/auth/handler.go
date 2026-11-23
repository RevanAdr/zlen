package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/RevanAdr/zlen/internal/httpx"
)

type UserHandler struct {
	service *UserService
	jwt     []byte
}

func NewUserHandler(service *UserService, token []byte) *UserHandler {
	return &UserHandler{service: service, jwt: token}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	responses, err := h.service.List(r.Context())
	if err != nil {
		httpx.WriteToJSON(w, http.StatusInternalServerError, "Internal server error", map[string]string{"error": string(err.Error())})
		return
	}
	httpx.WriteToJSON(w, http.StatusOK, "OK", responses)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, cerr := strconv.Atoi(r.PathValue("id"))
	if cerr != nil {
		httpx.WriteToJSON(w, http.StatusBadRequest, "Bad Request", map[string]string{"error": string(cerr.Error())})
		return
	}
	response, err := h.service.Get(r.Context(), id)
	switch {
	case errors.Is(ErrNotFound, err):
		httpx.WriteToJSON(w, http.StatusNotFound, "Not Found", map[string]string{"error": string(err.Error())})
		return
	case err != nil:
		httpx.WriteToJSON(w, http.StatusInternalServerError, "Internal Server Error", map[string]string{"error": "Internal Server Error"})
		return
	default:
		httpx.WriteToJSON(w, http.StatusOK, "OK", response)
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	request := UserRequest{}
	if rerr := httpx.ReadBody(r, &request); rerr != nil {
		httpx.WriteToJSON(w, http.StatusInternalServerError, "Internal Server Error", map[string]string{"error": string(rerr.Error())})
		return
	}

	hashed, hErr := HashPassword(request.Password)
	if hErr != nil {
		httpx.WriteToJSON(w, http.StatusInternalServerError, "Internal Server Error", map[string]string{"error": string(hErr.Error())})
		return
	}

	data := UserRequest{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashed,
	}
	response, err := h.service.Create(r.Context(), data)
	if err != nil {
		httpx.WriteToJSON(w, http.StatusInternalServerError, "Internal Server Error", map[string]string{"error": string(err.Error())})
		return
	}

	httpx.WriteToJSON(w, http.StatusCreated, "Created", response)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.WriteToJSON(w, http.StatusBadRequest, "Bad Request", map[string]string{"error": err.Error()})
		return
	}

	rErr := h.service.Delete(r.Context(), id)
	switch {
	case errors.Is(rErr, ErrNotFound):
		httpx.WriteToJSON(w, http.StatusNotFound, "Not Found", map[string]string{"error": "user not found"})
		return
	case rErr != nil:
		httpx.WriteToJSON(w, http.StatusInternalServerError, "Internal Server Error", map[string]string{"error": rErr.Error()})
		return
	default:
		httpx.WriteToJSON(w, http.StatusOK, "OK", map[string]string{"messages": "Succed delete user"})
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	request := LoginRequest{}
	if rErr := httpx.ReadBody(r, &request); rErr != nil {
		httpx.WriteToJSON(w, http.StatusBadRequest, "Bad Request", map[string]string{"error": string(rErr.Error())})
		return
	}

	data := LoginRequest{
		Email:    request.Email,
		Password: request.Password,
	}

	response, err := h.service.FindByEmail(r.Context(), data.Email)

	switch {
	case errors.Is(ErrNotFound, err):
		httpx.WriteToJSON(w, http.StatusNotFound, "Not Found", map[string]string{"error": string(err.Error())})
		return
	case err != nil:
		httpx.WriteToJSON(w, http.StatusInternalServerError, "Internal Server Error", map[string]string{"error": string(err.Error())})
		return
	}

	pwErr := MatchPassword(data.Password, response.Password)
	if pwErr != nil {
		httpx.WriteToJSON(w, http.StatusUnauthorized, "Unauthorized", map[string]string{"error": string(pwErr.Error())})
		return
	}

	formattedId := strconv.Itoa(response.ID)
	token, tokenErr := httpx.MakeToken(formattedId, h.jwt)

	if tokenErr != nil {
		slog.Error("make token", "err", tokenErr)
		httpx.WriteToJSON(w, http.StatusInternalServerError, "Internal Server Error", nil)
		return
	}

	appear := UserDisplay{
		Name:  response.Name,
		Email: response.Email,
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60,
	})
	httpx.WriteToJSON(w, http.StatusOK, "OK", appear)
}
