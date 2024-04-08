package handlers

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/unbeman/av-banner-task/internal/controller"
	"github.com/unbeman/av-banner-task/internal/models"
	"github.com/unbeman/av-banner-task/internal/storage"
	"github.com/unbeman/av-banner-task/internal/utils"
	"net/http"
)

type HttpHandler struct {
	*chi.Mux
	controller *controller.Controller
	jwtManager *utils.JWTManager
}

func NewHttpHandler(ctrl *controller.Controller, jwtManager *utils.JWTManager) (*HttpHandler, error) {
	//todo: setup routes
	h := &HttpHandler{
		Mux:        chi.NewMux(),
		controller: ctrl,
		jwtManager: jwtManager,
	}
	h.Route("/", func(router chi.Router) {
		router.Group(func(userRouter chi.Router) {
			userRouter.Use(h.userAuthorization)
			userRouter.Get("/user_banner", h.GetUserBanner)
		})
		router.Group(func(adminRouter chi.Router) {
			adminRouter.Use(h.adminAuthorization)
			adminRouter.Get("/banner", h.GetBanners)
			adminRouter.Post("/banner", h.CreateBanner)
			adminRouter.Patch("/banner/{id}", h.UpdateBanner)
			adminRouter.Delete("/banner/{id}", h.DeleteBanner)
		})

	})
	return h, nil
}

func (h HttpHandler) GetUserBanner(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	accessLevel := h.getAccessLevelFromContext(ctx)
	input := &models.GetBannerInput{}
	if err := render.Bind(request, input); err != nil {
		render.Render(writer, request, models.ErrBadRequest(err))
		return
	}

	var out *models.GetBannerOutput
	var err error

	switch accessLevel {
	case ADMIN:
		out, err = h.controller.GetBannerForAdmin(ctx, input)
	case USER:
		out, err = h.controller.GetBannerForUser(ctx, input)
	}

	if errors.Is(err, storage.ErrNotFound) {
		render.Render(writer, request, models.ErrNotFound(err))
		return
	}
	if err != nil {
		render.Render(writer, request, models.ErrInternalServerError(err))
		return
	}
	render.Render(writer, request, out)
}

func (h HttpHandler) GetBanners(writer http.ResponseWriter, request *http.Request) {

}

func (h HttpHandler) CreateBanner(writer http.ResponseWriter, request *http.Request) {

}

func (h HttpHandler) DeleteBanner(writer http.ResponseWriter, request *http.Request) {

}

func (h HttpHandler) UpdateBanner(writer http.ResponseWriter, request *http.Request) {

}

func (h HttpHandler) getAccessLevelFromContext(ctx context.Context) int {
	return ctx.Value(AccessContextKey).(int)
}
