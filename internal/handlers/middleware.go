package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"

	"github.com/unbeman/av-banner-task/internal/models"
	"github.com/unbeman/av-banner-task/internal/utils"
)

// Access levels for requests.
const (
	ADMIN = iota
	USER
)

var AccessContextKey = "access"

func (h HttpHandler) userAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		accessToken := getTokenFromRequest(request)
		userClaims, err := h.jwtManager.Verify(accessToken)
		if errors.Is(err, utils.ErrInvalidToken) {
			render.JSON(writer, request, models.ErrUnauthorized(err))
			return
		}
		if err != nil {
			render.JSON(writer, request, models.ErrInternalServerError(err))
			return
		}

		if userClaims.UserRole != ADMIN && userClaims.UserRole != USER {
			render.JSON(writer, request, models.ErrForbidden(fmt.Errorf("invalid user role")))
			return
		}
		contextWithAccess := context.WithValue(request.Context(), AccessContextKey, userClaims.UserRole)
		next.ServeHTTP(writer, request.WithContext(contextWithAccess))
	})
}

func (h HttpHandler) adminAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		accessToken := getTokenFromRequest(request)
		userClaims, err := h.jwtManager.Verify(accessToken)
		if errors.Is(err, utils.ErrInvalidToken) {
			render.JSON(writer, request, models.ErrUnauthorized(err))
			return
		}
		if err != nil {
			render.JSON(writer, request, models.ErrInternalServerError(err))
			return
		}

		if userClaims.UserRole != ADMIN {
			render.JSON(writer, request, models.ErrForbidden(fmt.Errorf("no access with given token")))
			return
		}

		next.ServeHTTP(writer, request)
	})
}

func getTokenFromRequest(request *http.Request) string {
	bearerToken := request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
