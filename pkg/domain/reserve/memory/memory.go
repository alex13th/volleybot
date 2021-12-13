package memory

import (
	"fmt"
	"sync"
	"volleybot/pkg/domain/reserve"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	reserves map[uuid.UUID]reserve.Reserve
	sync.Mutex
}

func New() *MemoryRepository {
	return &MemoryRepository{
		reserves: make(map[uuid.UUID]reserve.Reserve),
	}
}

func (mr *MemoryRepository) Get(id uuid.UUID) (reserve.Reserve, error) {
	if reserve, ok := mr.reserves[id]; ok {
		return reserve, nil
	}

	return reserve.Reserve{}, reserve.ErrReserveNotFound
}

func (mr *MemoryRepository) Add(memr reserve.Reserve) error {
	if mr.reserves == nil {
		mr.Lock()
		mr.reserves = make(map[uuid.UUID]reserve.Reserve)
		mr.Unlock()
	}

	if _, ok := mr.reserves[memr.Id]; ok {
		return fmt.Errorf("reserve already exists: %w", reserve.ErrFailedToAddReserve)
	}
	mr.Lock()
	mr.reserves[memr.Id] = memr
	mr.Unlock()
	return nil
}

func (mr *MemoryRepository) Update(memr reserve.Reserve) error {
	if _, ok := mr.reserves[memr.Id]; !ok {
		return fmt.Errorf("reserve does not exist: %w", reserve.ErrUpdateReserve)
	}
	mr.Lock()
	mr.reserves[memr.Id] = memr
	mr.Unlock()
	return nil
}
