package rpcgroups

import (
	"context"
	"fmt"
	"personal_bot/infrastructure/persistence/repository"
	rpcgroups "personal_bot/internal/core/rpc_groups"
	"sync"
)

type Service struct {
	repo            *repository.RPCGroup
	loadedRpcGroups map[int64]*rpcgroups.LoadedRPCGroup
	mu              sync.Mutex
}

func NewRPCGroupService(repo *repository.RPCGroup) *Service {
	return &Service{repo: repo, loadedRpcGroups: map[int64]*rpcgroups.LoadedRPCGroup{}}
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

func (s *Service) Load(ctx context.Context, id int64) (rpcGroup rpcgroups.RPCGroup, err error) {
	//a trading task will call this, load id
	s.mu.Lock()
	defer s.mu.Unlock()

	loadedGroup, ok := s.loadedRpcGroups[id]
	if !ok {
		rpcGroup, err = s.GetBy(ctx, id)
		if err != nil {
			return
		}

		s.loadedRpcGroups[id] = &rpcgroups.LoadedRPCGroup{
			RPCGroup:   rpcGroup,
			GroupIndex: 0,
			References: 1,
		}

		return
	}

	loadedGroup.References++

	return loadedGroup.RPCGroup, nil
}

func (s *Service) Unload(id int64) error {
	// a trading task will call this, load id
	s.mu.Lock()
	defer s.mu.Unlock()
	loadedGroup, ok := s.loadedRpcGroups[id]
	if !ok {
		return fmt.Errorf("rpc group not found with id %d", id)
	}

	loadedGroup.References--
	if loadedGroup.References == 0 {
		delete(s.loadedRpcGroups, id)
	}

	return nil
}

func (s *Service) GetNode(id int64) (rpcgroups.GroupItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rpcGroup, ok := s.loadedRpcGroups[id]
	if !ok {
		return rpcgroups.GroupItem{}, fmt.Errorf("rpc group not loaded in app - double check logic")
	}

	node := rpcGroup.RPCGroup.Group[rpcGroup.GroupIndex]
	rpcGroup.GroupIndex = (rpcGroup.GroupIndex + 1) % len(rpcGroup.RPCGroup.Group)
	return node, nil
}
