package router

import (
	"context"
	"net/http"
)

func (r *Router) RegistryPrivateRoute(ctx context.Context, path string, handler func(http.ResponseWriter, *http.Request)) {
	r.mux.HandleFunc(r.prefix+path+r.suffix, handler)
}
