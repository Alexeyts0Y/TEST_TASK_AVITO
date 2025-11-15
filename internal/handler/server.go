package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	errWrappers "github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/errors"
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

func (s *Server) PostTeamAdd(ctx context.Context, request api.PostTeamAddRequestObject) (api.PostTeamAddResponseObject, error) {
	teamToAdd := *request.Body
	_, err := s.Repository.SaveTeam(ctx, teamToAdd)
	if err != nil && errors.Is(err, errWrappers.ErrTeamExists) {
		return api.PostTeamAdd400JSONResponse(newErrorResponse(api.TEAMEXISTS, fmt.Sprintf("Команда с именем %s уже существует", teamToAdd.TeamName))), nil
	} else if err != nil {
		return nil, err
	}

	return api.PostTeamAdd201JSONResponse{Team: &teamToAdd}, nil
}

func (s *Server) GetTeamGet(ctx context.Context, request api.GetTeamGetRequestObject) (api.GetTeamGetResponseObject, error) {
	teamName := request.Params.TeamName
	team, err := s.Repository.GetTeam(ctx, teamName)
	if err != nil && errors.Is(err, errWrappers.ErrNotFound) {
		return api.GetTeamGet404JSONResponse(newErrorResponse(api.NOTFOUND, fmt.Sprintf("Команда с именем %s не найдена", teamName))), nil
	} else if err != nil {
		return nil, err
	}

	return api.GetTeamGet200JSONResponse(team), nil
}

func (s *Server) PostUsersSetIsActive(ctx context.Context, request api.PostUsersSetIsActiveRequestObject) (api.PostUsersSetIsActiveResponseObject, error) {
	userID := request.Body.UserId
	isActive := request.Body.IsActive

	user, err := s.Repository.SetUserIsActive(ctx, userID, isActive)
	if err != nil && errors.Is(err, errWrappers.ErrNotFound) {
		return api.PostUsersSetIsActive404JSONResponse(newErrorResponse(api.NOTFOUND, fmt.Sprintf("User %s not found", userID))), nil
	} else if err != nil {
		return nil, err
	}

	return api.PostUsersSetIsActive200JSONResponse{User: &user}, nil
}

func (s *Server) GetUsersGetReview(ctx context.Context, request api.GetUsersGetReviewRequestObject) (api.GetUsersGetReviewResponseObject, error) {
	// TODO: implement method
	return nil, nil
}

func (s *Server) PostPullRequestCreate(ctx context.Context, request api.PostPullRequestCreateRequestObject) (api.PostPullRequestCreateResponseObject, error) {
	body := request.Body

	author, err := s.Repository.GetUser(ctx, body.AuthorId)
	if err != nil && errors.Is(err, errWrappers.ErrNotFound) {
		return api.PostPullRequestCreate404JSONResponse(newErrorResponse(api.NOTFOUND, "Автор с таким ID не найден")), nil
	} else if err != nil {
		return nil, err
	}

	candidates, err := s.Repository.FindActiveCandidates(ctx, author.TeamName, []string{author.UserId})
	if err != nil && errors.Is(err, errWrappers.ErrNotFound) {
		return api.PostPullRequestCreate404JSONResponse(newErrorResponse(api.NOTFOUND, "Команда с таким автором не найдена")), nil
	}

	assignedReviewers := s.Repository.ChooseRandomCandidates(candidates, 2)

	newPullRequest := api.PullRequest{
		PullRequestId:     body.PullRequestId,
		PullRequestName:   body.PullRequestName,
		AuthorId:          body.AuthorId,
		AssignedReviewers: assignedReviewers,
		Status:            api.PullRequestStatusOPEN,
	}

	savedPullRequest, err := s.Repository.SavePullRequest(ctx, newPullRequest)
	if err != nil && errors.Is(err, errWrappers.ErrPrExists) {
		return api.PostPullRequestCreate409JSONResponse(newErrorResponse(api.PREXISTS, "Такой пул реквест уже существует")), nil
	} else if err != nil {
		return nil, err
	}

	return api.PostPullRequestCreate201JSONResponse{Pr: &savedPullRequest}, nil
}

func (s *Server) PostPullRequestMerge(ctx context.Context, request api.PostPullRequestMergeRequestObject) (api.PostPullRequestMergeResponseObject, error) {
	pullRequestId := request.Body.PullRequestId

	pullRequest, err := s.Repository.GetPullRequest(ctx, pullRequestId)
	if err != nil && errors.Is(err, errWrappers.ErrNotFound) {
		return api.PostPullRequestMerge404JSONResponse(newErrorResponse(api.NOTFOUND, "Пул реквест с таким ID не найден")), nil
	} else if err != nil {
		return nil, err
	}

	if pullRequest.Status == api.PullRequestStatusMERGED {
		return api.PostPullRequestMerge200JSONResponse{Pr: &pullRequest}, nil
	}

	currentTime := time.Now()
	pullRequest.Status = api.PullRequestStatusMERGED
	pullRequest.MergedAt = &currentTime

	mergedPullRequest, err := s.Repository.UpdatePullRequest(ctx, pullRequest)
	if err != nil {
		return nil, err
	}

	return api.PostPullRequestMerge200JSONResponse{Pr: &mergedPullRequest}, nil
}

func (s *Server) PostPullRequestReassign(ctx context.Context, request api.PostPullRequestReassignRequestObject) (api.PostPullRequestReassignResponseObject, error) {
	// TODO: implement method
	return nil, nil
}
