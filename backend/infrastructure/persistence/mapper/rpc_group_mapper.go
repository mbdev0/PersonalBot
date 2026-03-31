package mapper

import (
	"encoding/json"
	"personal_bot/infrastructure/persistence/models"
	rpcgroups "personal_bot/internal/core/rpc_groups"
)

func MapRpcGroupPostToRepository(src rpcgroups.RPCGroupPost) (models.RpcGroupRepository, error) {
	dst := models.RpcGroupRepository{}
	dst.GroupName = src.Name
	rpcGroups, err := json.Marshal(src.Group)
	if err != nil {
		return dst, err
	}
	dst.RpcGroups = string(rpcGroups)
	dst.CreationTime = src.CreationTime

	return dst, nil
}

func MapRpcGroupPutToRepository(src rpcgroups.RPCGroupPut) (models.RpcGroupRepository, error) {
	dst := models.RpcGroupRepository{}
	dst.GroupName = src.Name
	rpcGroups, err := json.Marshal(src.Group)
	if err != nil {
		return dst, err
	}
	dst.RpcGroups = string(rpcGroups)

	return dst, nil
}

func MapRepositoryToRpcGroup(src models.RpcGroupRepository) (rpcgroups.RPCGroup, error) {
	dst := rpcgroups.RPCGroup{}
	dst.Name = src.GroupName
	var groups rpcgroups.Group
	err := json.Unmarshal([]byte(src.RpcGroups), &groups)
	if err != nil {
		return dst, err
	}

	dst.Group = groups
	dst.CreationTime = src.CreationTime
	dst.Id = src.Id

	return dst, nil
}
