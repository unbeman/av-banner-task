package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/unbeman/av-banner-task/docs"
	"github.com/unbeman/av-banner-task/internal/controller"
	"github.com/unbeman/av-banner-task/internal/models"
	"github.com/unbeman/av-banner-task/internal/storage"
	"github.com/unbeman/av-banner-task/internal/utils"
)

const BannerIDParam = "id"

type HttpHandler struct {
	*chi.Mux
	controller *controller.Controller
	jwtManager *utils.JWTManager
}

func NewHttpHandler(ctrl *controller.Controller, jwtManager *utils.JWTManager) (*HttpHandler, error) {
	h := &HttpHandler{
		Mux:        chi.NewMux(),
		controller: ctrl,
		jwtManager: jwtManager,
	}
	h.Use(logger.Logger("router", log.StandardLogger()))
	h.Get("/swagger/*", httpSwagger.Handler()) // todo: переместить
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

// GetUserBanner godoc
// @Summary Получение баннера
// @Description Возвращает баннер по заданному feature_id и tag_id
// @Produce json
// @Param feature_id query integer true "Идентификатор фичи"
// @Param tag_id query integer true "Идентификатор тэга"
// @Success 200 {object} models.GetBannerOutput
// @Failure 400 {object} models.ErrResponse
// @Failure 401 {object} models.ErrResponse
// @Failure 403 {object} models.ErrResponse
// @Failure 404 {object} models.ErrResponse
// @Failure 500 {object} models.ErrResponse
// @Security Bearer
// @Router /user_banner [get]
func (h HttpHandler) GetUserBanner(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	accessLevel := h.getAccessLevelFromContext(ctx)

	input := &models.GetBannerInput{}
	if err := input.FromURI(request); err != nil {
		render.Render(writer, request, models.ErrBadRequest(err))
		return
	}

	var out *models.GetBannerOutput
	var err error

	switch accessLevel {
	case ADMIN:
		out, err = h.controller.GetBanner(ctx, input, nil)
	case USER:
		isActive := true
		out, err = h.controller.GetBanner(ctx, input, &isActive)
	}

	if errors.Is(err, storage.ErrNotFound) {
		render.Render(writer, request, models.ErrNotFound(err))
		return
	}
	if err != nil {
		render.Render(writer, request, models.ErrInternalServerError(err))
		return
	}
	render.JSON(writer, request, json.RawMessage(*out))
}

// GetBanners godoc
// @Summary Получение списка баннеров
// @Description Возвращает список баннеров по заданной фильтрации feature_id и/или tag_id
// @Produce json
// @Param feature_id query integer false "Идентификатор фичи"
// @Param tag_id query integer false "Идентификатор тэга"
// @Param limit query integer false "Лимит выдачи"
// @Param offset query integer false "Сдвиг выдачи"
// @Success 200 {object} models.Banners
// @Failure 400 {object} models.ErrResponse
// @Failure 401 {object} models.ErrResponse
// @Failure 403 {object} models.ErrResponse
// @Failure 500 {object} models.ErrResponse
// @Security Bearer
// @Router /banner [get]
func (h HttpHandler) GetBanners(writer http.ResponseWriter, request *http.Request) {
	input := &models.GetBannersInput{}

	if err := input.FromURI(request); err != nil {
		render.Render(writer, request, models.ErrBadRequest(err))
		return
	}

	out, err := h.controller.GetBanners(request.Context(), input)
	if err != nil {
		render.Render(writer, request, models.ErrInternalServerError(err))
		return
	}
	render.JSON(writer, request, out)
}

// CreateBanner godoc
// @Summary Создание баннера
// @Description Заводит новый баннер с заданными полями
// @Accept json
// @Produce json
// @Param input body models.CreateBannerInput true "Информация о добавляемом баннере"
// @Success 201 {object} models.CreateBannerOutput
// @Failure 400 {object} models.ErrResponse
// @Failure 401 {object} models.ErrResponse
// @Failure 403 {object} models.ErrResponse
// @Failure 404 {object} models.ErrResponse
// @Failure 409 {object} models.ErrResponse
// @Failure 500 {object} models.ErrResponse
// @Security Bearer
// @Router /banner [post]
func (h HttpHandler) CreateBanner(writer http.ResponseWriter, request *http.Request) {
	input := &models.CreateBannerInput{}
	if err := render.Bind(request, input); err != nil {
		render.Render(writer, request, models.ErrBadRequest(err))
		return
	}
	out, err := h.controller.CreateBanner(request.Context(), input)
	if errors.Is(err, storage.ErrConflict) {
		render.Render(writer, request, models.ErrBadRequest(err))
		return
	}
	if err != nil {
		render.Render(writer, request, models.ErrInternalServerError(err))
		return
	}
	render.Status(request, http.StatusCreated)
	render.JSON(writer, request, out)
}

// UpdateBanner godoc
// @Summary Обновление баннера
// @Description Обновляет параметры существующего баннера
// @Accept json
// @Produce json
// @Param id path integer true "Идентификатор баннера"
// @Param input body models.UpdateBannerInput true "Информация об обновлении баннера"
// @Success 200
// @Failure 400 {object} models.ErrResponse
// @Failure 401 {object} models.ErrResponse
// @Failure 403 {object} models.ErrResponse
// @Failure 404 {object} models.ErrResponse
// @Failure 409 {object} models.ErrResponse
// @Failure 500 {object} models.ErrResponse
// @Security Bearer
// @Router /banner/{id} [patch]
func (h HttpHandler) UpdateBanner(writer http.ResponseWriter, request *http.Request) {
	input := &models.UpdateBannerInput{}

	bannerId, err := getBannerIDFromURI(request)
	if err != nil {
		render.Render(writer, request, models.ErrBadRequest(err))
		return
	}

	if err = render.Bind(request, input); err != nil {
		render.Render(writer, request, models.ErrBadRequest(err))
		return
	}

	input.Id = bannerId

	err = h.controller.UpdateBanner(request.Context(), input)
	if errors.Is(err, storage.ErrNotFound) {
		render.Render(writer, request, models.ErrBadRequest(err))
		return
	}
	if errors.Is(err, storage.ErrConflict) {
		render.Render(writer, request, models.ErrBadRequest(err))
		return
	}
	if err != nil {
		render.Render(writer, request, models.ErrInternalServerError(err))
		return
	}

	render.Status(request, http.StatusOK)
}

// DeleteBanner godoc
// @Summary Удаление баннера баннера
// @Description Удаляет баннер по заданному идентификатору
// @Produce json
// @Param id path integer true "Идентификатор баннера"
// @Success 204
// @Failure 400 {object} models.ErrResponse
// @Failure 401 {object} models.ErrResponse
// @Failure 403 {object} models.ErrResponse
// @Failure 404 {object} models.ErrResponse
// @Failure 500 {object} models.ErrResponse
// @Security Bearer
// @Router /banner/{id} [delete]
func (h HttpHandler) DeleteBanner(writer http.ResponseWriter, request *http.Request) {
	bannerId, err := getBannerIDFromURI(request)
	if err != nil {
		render.Render(writer, request, models.ErrBadRequest(err))
		return
	}
	err = h.controller.DeleteBanner(request.Context(), bannerId)
	if errors.Is(err, storage.ErrNotFound) {
		render.Render(writer, request, models.ErrBadRequest(err))
		return
	}
	if err != nil {
		render.Render(writer, request, models.ErrInternalServerError(err))
		return
	}

	render.Status(request, http.StatusNoContent)
}

func (h HttpHandler) getAccessLevelFromContext(ctx context.Context) int {
	return ctx.Value(AccessContextKey).(int)
}

func getBannerIDFromURI(request *http.Request) (int, error) {
	rawID := chi.URLParam(request, BannerIDParam)
	return strconv.Atoi(rawID)
}
