package app

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/av-banner-task/internal/controller"
	"github.com/unbeman/av-banner-task/internal/handlers"
	"github.com/unbeman/av-banner-task/internal/utils"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(ctrl *controller.Controller, jwtManager *utils.JWTManager) (*HTTPServer, error) {
	handler, err := handlers.NewHttpHandler(ctrl, jwtManager)
	if err != nil {
		return nil, fmt.Errorf("couldn't setup HTTP server: %w", err)
	}

	hs := &HTTPServer{
		server: &http.Server{
			Handler: handler,
			Addr:    ":8080",
		},
	}
	return hs, nil
}

func (h *HTTPServer) GetServer() *http.Server {
	return h.server
}

func (h *HTTPServer) GetAddress() string {
	return h.server.Addr
}

func (h *HTTPServer) Run() {
	log.Info("starting HTTP server")
	h.server.ListenAndServe()
}

func (h *HTTPServer) Close() {
	log.Info("http server closed")
	h.server.Shutdown(context.TODO())
}
