package handler

import (
	"context"
	"errors"
	"fmt"

	errWrappers "github.com/Alexeyts0Y/TEST_TASK_AVITO/internal/errors"
	"github.com/Alexeyts0Y/TEST_TASK_AVITO/pkg/api"
)

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

func (s *Server) PostTeamTeamNameDeactivateMembers(ctx context.Context, request api.PostTeamTeamNameDeactivateMembersRequestObject) (api.PostTeamTeamNameDeactivateMembersResponseObject, error) {
	count, err := s.Repository.DeactivateTeamMembers(ctx, request.TeamName)
	if err != nil && errors.Is(err, errWrappers.ErrNotFound) {
		return api.PostTeamTeamNameDeactivateMembers404JSONResponse(newErrorResponse(api.NOTFOUND, "Команда с таким именем не найдена")), nil
	} else if err != nil {
		return nil, err
	}
	return api.PostTeamTeamNameDeactivateMembers200JSONResponse{TeamName: request.TeamName, DeactivatedUsersCount: int(count)}, nil
}

func (s *Server) PostTeamReassignPrs(ctx context.Context, request api.PostTeamReassignPrsRequestObject) (api.PostTeamReassignPrsResponseObject, error) {
	summary, err := s.Repository.ReassignPRsForTeam(ctx, request.TeamName)

	if err != nil && errors.Is(err, errWrappers.ErrNotFound) {
		return api.PostTeamReassignPrs404JSONResponse(newErrorResponse(api.NOTFOUND, "Команда с таким именем не найдена")), nil
	} else if err != nil {
		return nil, err
	}

	return api.PostTeamReassignPrs200JSONResponse(summary), nil
}
