package v1

import (
	"net/http"

	"github.com/bogdanovds/rocket_factory/order/internal/service"
	orderV1 "github.com/bogdanovds/rocket_factory/shared/pkg/openapi/order/v1"
)

type Handler struct {
	service service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{service: svc}
}

func badRequest(msg string) *orderV1.BadRequestError {
	return &orderV1.BadRequestError{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

func notFound(msg string) *orderV1.NotFoundError {
	return &orderV1.NotFoundError{
		Code:    http.StatusNotFound,
		Message: msg,
	}
}

func conflict(msg string) *orderV1.ConflictError {
	return &orderV1.ConflictError{
		Code:    http.StatusConflict,
		Message: msg,
	}
}
