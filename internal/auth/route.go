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

	return mux
}
