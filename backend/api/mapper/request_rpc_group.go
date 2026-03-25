package mapper

import (
	"fmt"
	"personal_bot/api/dto"
	rpcgroups "personal_bot/internal/core/rpc_groups"
	"regexp"
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

func parseGroupString(groupString string) (rpcgroups.Group, error) {
	//group string will be formatted as "http:http://abc.com, ws:ws://abc.com\nhttp:http://abc.com, ws:ws://abc.com\n"
	// we will need to return []rpcgroups.GroupItem
	rpcLineRegex, err := regexp.Compile(`http:([^,]+),\s*ws:([^\n]+)`)
	if err != nil {
		return rpcgroups.Group{}, err
	}
	matches := rpcLineRegex.FindAllStringSubmatch(groupString, -1)

	group := rpcgroups.Group{}
	for _, match := range matches {
		http := strings.TrimSpace(match[1])
		ws := strings.TrimSpace(match[2])

		groupItem := rpcgroups.GroupItem{
			Http: http,
			WS:   ws,
		}

		group = append(group, groupItem)
	}

	return group, nil
}
