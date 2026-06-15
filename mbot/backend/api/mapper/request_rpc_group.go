package mapper

import (
	"fmt"
	"personal_bot/backend/api/dto"
	rpcgroups "personal_bot/backend/internal/core/rpc_groups"
	"strings"
	"time"
)

func MapRPCGroupToDto(src rpcgroups.RPCGroup) dto.RPCGroup {
	returnVal := dto.RPCGroup{
		Name:         src.Name,
		CreationTime: src.CreationTime,
		Id:           src.Id,
	}

	dtoGroup := []dto.Group{}
	for _, srcGroupItem := range src.Group {
		dtoGroupItem := dto.Group{
			Http: srcGroupItem.Http,
			Ws:   srcGroupItem.WS,
		}

		dtoGroup = append(dtoGroup, dtoGroupItem)
	}

	returnVal.Group = dtoGroup
	return returnVal
}

func MapRpcGroupPushToRpcGroupPut(src dto.RPCGroupPush) (mappedRcpGroup rpcgroups.RPCGroupPut, err error) {
	mappedRcpGroup.Name = src.Name

	group, err := parseGroupString(src.Group)
	if err != nil {
		return
	}

	mappedRcpGroup.Group = group
	return
}

func MapRpcGroupPostDtoToRpcGroupPost(src dto.RPCGroupPush) (rpcgroups.RPCGroupPost, error) {
	rpcgroup := rpcgroups.RPCGroupPost{}
	rpcgroup.CreationTime = fmt.Sprint(time.Now().Unix())

	//we need to parse the string to create groups
	group, err := parseGroupString(src.Group)
	if err != nil {
		return rpcgroups.RPCGroupPost{}, err
	}

	rpcgroup.Group = group
	rpcgroup.Name = src.Name

	return rpcgroup, nil
}

func MapRPCGroupToResponseDto(src rpcgroups.RPCGroup) dto.RPCGroupResponse {
	returnVal := dto.RPCGroupResponse{
		Name:         src.Name,
		CreationTime: src.CreationTime,
		Id:           src.Id,
	}

	var groupString strings.Builder
	for _, group := range src.Group {
		fmt.Fprintf(&groupString, "%s, %s\n", group.Http, group.WS)
	}

	returnVal.Group = groupString.String()

	return returnVal

}

func parseGroupString(groupString string) (rpcgroups.Group, error) {
	lines := strings.Split(strings.TrimSpace(groupString), "\n")
	group := make(rpcgroups.Group, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ",", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid rpc line format: %q", line)
		}

		http := strings.TrimSpace(parts[0])
		ws := strings.TrimSpace(parts[1])

		if !strings.HasPrefix(http, "http") {
			return nil, fmt.Errorf("invalid http url: %q", http)
		}
		if !strings.HasPrefix(ws, "ws") {
			return nil, fmt.Errorf("invalid ws url: %q", ws)
		}

		group = append(group, rpcgroups.GroupItem{Http: http, WS: ws})
	}

	return group, nil
}
