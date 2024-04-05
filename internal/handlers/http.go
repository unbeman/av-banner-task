package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/unbeman/av-banner-task/internal/controller"
)

type HttpHandler struct {
	*chi.Mux
	controller *controller.Controller
}

func NewHttpHandler(ctrl *controller.Controller) (*HttpHandler, error) {
	//todo: setup routes
	handler := &HttpHandler{controller: ctrl}
	return handler, nil
}
