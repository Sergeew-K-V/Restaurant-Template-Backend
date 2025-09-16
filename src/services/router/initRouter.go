package router

import (
	"context"
	"fmt"
	"net/http"
)

func InitRouter(ctx context.Context, mux *http.ServeMux, prefix string, suffix string) (RouterInterface, error) {
	if ctx == nil {
		return nil, fmt.Errorf("error to init router")
	}

	if mux == nil {
		return nil, fmt.Errorf("error to init router")
	}

	return &Router{ctx: ctx, mux: mux, suffix: suffix, prefix: prefix}, nil
}
