package app

import (
	"context"
	"fmt"
	"github.com/unbeman/av-banner-task/internal/controller"
	"github.com/unbeman/av-banner-task/internal/handlers"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(addr string, ctrl *controller.Controller) (*HTTPServer, error) {
	handler, err := handlers.NewHttpHandler(ctrl)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup HTTP server: %w", err)
	}

	hs := &HTTPServer{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
	return hs, nil
}

func (h *HTTPServer) GetAddress() string {
	return h.server.Addr
}

func (h *HTTPServer) Run() error {
	return h.server.ListenAndServe()
}

func (h *HTTPServer) Close() error {
	return h.server.Shutdown(context.TODO())
}
