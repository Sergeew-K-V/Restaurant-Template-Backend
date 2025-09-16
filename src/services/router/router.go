package router

import (
	"context"
	"net/http"
)

type Router struct {
	ctx    context.Context
	mux    *http.ServeMux
	suffix string
	prefix string
}

type RouterInterface interface {
	RegistryPublicRoute(ctx context.Context, path string, handler func(http.ResponseWriter, *http.Request))
	RegistryPrivateRoute(ctx context.Context, path string, handler func(http.ResponseWriter, *http.Request))
}
