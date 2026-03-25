package rpcgroups

import (
	"context"
	"personal_bot/infrastructure/persistence/repository"
	rpcgroups "personal_bot/internal/core/rpc_groups"
)

type Service struct {
	repo *repository.RPCGroup
}

func NewRPCGroupService(repo *repository.RPCGroup) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetDashboard(ctx context.Context) (rpcgroups.RPCGroupDashboard, error) {
	dashboard, err := s.repo.GetDashboard(ctx)
	if err != nil {
		return nil, err
	}

	return dashboard, nil
}

func (s *Service) Post(ctx context.Context, rpcGroup rpcgroups.RPCGroupPost) (rpcgroups.RPCGroup, error) {
	return s.repo.Add(ctx, rpcGroup)
}

func (s *Service) GetBy(ctx context.Context, id int64) (rpcgroups.RPCGroup, error) {
	rpcGroup, err := s.repo.GetById(ctx, id)
	if err != nil {
		return rpcgroups.RPCGroup{}, err
	}

	return rpcGroup, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, rcpGroupPut rpcgroups.RPCGroupPut) (rcpGroup rpcgroups.RPCGroup, err error) {
	return s.repo.Update(ctx, id, rcpGroupPut)
}

func (s *Service) Load() {
	// so when we do a getby for a rpc group, if it doesn't exist in the rpc groups map, we load it into the service so we can quikcly call it
	//TODO: Load RPC group in memory?? - do we load every one in memory? -> so when strategy tasks call GetNode it'll be a quicker call
}

func (s *Service) Unload() {
	// when the map no longer has a specific group inside of it, we will unload it from memory
	//TODO: Unload RPC group in memory -> when no more strategy tasks use the Node group, remove from memory
}

func (s *Service) GetNode() {
	//TODO: For a given node group, get the node group at index i + 1, i is the last accessed node and return it
}
