package handler

import (
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/repository"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/pkg/api"
)

type Server struct {
	Repository repository.Repository
}

func NewServer(repository repository.Repository) *Server {
	return &Server{Repository: repository}
}

func newErrorResponse(code api.ErrorResponseErrorCode, message string) api.ErrorResponse {
	return api.ErrorResponse{
		Error: struct {
			Code    api.ErrorResponseErrorCode "json:\"code\""
			Message string                     "json:\"message\""
		}{
			Code:    code,
			Message: message,
		},
	}
}
