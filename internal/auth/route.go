package auth

import "net/http"

type UserRoute interface {
	NewRouter() *http.ServeMux
}

type UserRouteImpl struct {
	handler UserHandler
}

func NewUseRoute(handler UserHandler) UserRoute {
	return &UserRouteImpl{handler: handler}
}

func (r *UserRouteImpl) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/users", r.handler.List)
	mux.HandleFunc("GET /api/v1/users/{id}", r.handler.Get)
	mux.HandleFunc("POST /api/v1/users", r.handler.Create)
	mux.HandleFunc("DELETE /api/v1/users/{id}", r.handler.Delete)

	return mux
}
