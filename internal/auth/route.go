package auth

import (
	"net/http"

	"github.com/RevanAdr/zlen/internal/httpx"
)

type UserRoute interface {
	NewRouter() *http.ServeMux
}

type UserRouteImpl struct {
	handler UserHandler
	jwt     []byte
}

func NewUseRoute(handler UserHandler, jwt []byte) UserRoute {
	return &UserRouteImpl{handler: handler, jwt: jwt}
}

func (r *UserRouteImpl) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	auth := httpx.Authentication(r.jwt)

	// authentication requeired
	mux.Handle("GET /api/v1/employees", auth(http.HandlerFunc(r.handler.List)))
	mux.Handle("GET /api/v1/employees/{id}", auth(http.HandlerFunc(r.handler.Get)))
	mux.Handle("POST /api/v1/employees", auth(http.HandlerFunc(r.handler.Create)))
	mux.Handle("DELETE /api/v1/employees/{id}", auth(http.HandlerFunc(r.handler.Delete)))

	// public
	mux.HandleFunc("POST /employees/login", r.handler.Login)

	return mux
}
