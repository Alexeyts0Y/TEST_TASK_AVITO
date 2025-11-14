package errors

import (
	"errors"
	"fmt"

	"github.com/Alexeyts0Y/TEST_TASK_AVITO/pkg/api"
)

type ApiError struct {
	Code    api.ErrorResponseErrorCode
	Message string
	Err     error
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%s: %s:", e.Code, e.Message)
}

func (e *ApiError) Is(target error) bool {
	if targetError, ok := target.(*ApiError); ok {
		return e.Code == targetError.Code
	}
	return errors.Is(e.Err, target)
}

var (
	ErrNotFound    = &ApiError{Code: api.NOTFOUND, Message: "resource not found"}
	ErrNotAssigned = &ApiError{Code: api.NOTASSIGNED, Message: "reviewer is not assigned to this PR"}
	ErrNoCandidate = &ApiError{Code: api.NOCANDIDATE, Message: "no active replacement candidate in team"}
	ErrPrExists    = &ApiError{Code: api.PREXISTS, Message: "PR id already exists"}
	ErrPrMerged    = &ApiError{Code: api.PRMERGED, Message: "cannot reassign on merged PR"}
	ErrTeamExists  = &ApiError{Code: api.TEAMEXISTS, Message: "team_name already exists"}
)
