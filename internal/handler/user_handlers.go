package handler

import (
	"context"
	"errors"
	"fmt"

	errWrappers "github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/errors"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/pkg/api"
)

func (s *Server) PostUsersSetIsActive(ctx context.Context, request api.PostUsersSetIsActiveRequestObject) (api.PostUsersSetIsActiveResponseObject, error) {
	userId := request.Body.UserId
	isActive := request.Body.IsActive

	user, err := s.Repository.SetUserIsActive(ctx, userId, isActive)
	if err != nil && errors.Is(err, errWrappers.ErrNotFound) {
		return api.PostUsersSetIsActive404JSONResponse(newErrorResponse(api.NOTFOUND, fmt.Sprintf("Пользователь %s не найден", userId))), nil
	} else if err != nil {
		return nil, err
	}

	return api.PostUsersSetIsActive200JSONResponse{User: &user}, nil
}

func (s *Server) GetUsersGetReview(ctx context.Context, request api.GetUsersGetReviewRequestObject) (api.GetUsersGetReviewResponseObject, error) {
	userId := request.Params.UserId

	_, err := s.Repository.GetUser(ctx, userId)
	if err != nil && errors.Is(err, errWrappers.ErrNotFound) {
		return api.GetUsersGetReview200JSONResponse{
			UserId:       userId,
			PullRequests: []api.PullRequestShort{},
		}, nil
	} else if err != nil {
		return nil, err
	}

	pullRequests, err := s.Repository.FindUserPullRequests(ctx, userId)
	if err != nil {
		return nil, err
	}

	return api.GetUsersGetReview200JSONResponse{
		PullRequests: pullRequests,
		UserId:       userId,
	}, nil
}
