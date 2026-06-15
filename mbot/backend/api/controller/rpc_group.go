package controller

import (
	"context"
	"personal_bot/backend/api/dto"
	"personal_bot/backend/api/mapper"
	models "personal_bot/backend/internal/core/rpc_groups"
	rpcgroups "personal_bot/backend/internal/services/rpc_groups"
)

type RPCGroupController struct {
	service *rpcgroups.Service
}

func NewRPCGroupController(service *rpcgroups.Service) *RPCGroupController {
	return &RPCGroupController{service: service}
}

func (rgc *RPCGroupController) GetDashboard(ctx context.Context) (models.RPCGroupDashboard, error) {
	dashboard, err := rgc.service.GetDashboard(ctx)
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}

func (rgc *RPCGroupController) Post(ctx context.Context, rpcGroup dto.RPCGroupPush) (dto.RPCGroup, error) {
	// we do not assign id in be, sql will do it - we should still send back the new rpc group
	rpcGroupPost, err := mapper.MapRpcGroupPostDtoToRpcGroupPost(rpcGroup)
	if err != nil {
		return dto.RPCGroup{}, err
	}

	createdRpcGroup, err := rgc.service.Post(ctx, rpcGroupPost)
	if err != nil {
		return dto.RPCGroup{}, err
	}

	return mapper.MapRPCGroupToDto(createdRpcGroup), nil
}

func (rgc *RPCGroupController) GetBy(ctx context.Context, id int64) (dto.RPCGroupResponse, error) {
	rpcGroup, err := rgc.service.GetBy(ctx, id)
	if err != nil {
		return dto.RPCGroupResponse{}, err
	}
	return mapper.MapRPCGroupToResponseDto(rpcGroup), nil
}

func (rgc *RPCGroupController) Delete(ctx context.Context, id int64) error {
	return rgc.service.Delete(ctx, id)
}

func (rgc *RPCGroupController) Update(ctx context.Context, id int64, rcpGroup dto.RPCGroupPush) (returnVal dto.RPCGroup, err error) {
	mappedRcpGroup, err := mapper.MapRpcGroupPushToRpcGroupPut(rcpGroup)
	if err != nil {
		return
	}
	rcpGroupUpdated, err := rgc.service.Update(ctx, id, mappedRcpGroup)
	if err != nil {
		return
	}

	return mapper.MapRPCGroupToDto(rcpGroupUpdated), nil
}
