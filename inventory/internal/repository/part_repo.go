package repository

import (
	"errors"
	"sync"

	"inventory/pkg/model"
)

var (
	ErrPartNotFound = errors.New("part not found")
)

type PartRepository interface {
	Save(part *model.Part) error
	FindByUUID(uuid string) (*model.Part, error)
	FindAll() ([]*model.Part, error)
	Update(part *model.Part) error
}

type InMemoryPartRepository struct {
	mu    sync.RWMutex
	parts map[string]*model.Part
}

func NewInMemoryPartRepository() *InMemoryPartRepository {
	return &InMemoryPartRepository{
		parts: make(map[string]*model.Part),
	}
}

func (r *InMemoryPartRepository) Save(part *model.Part) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.parts[part.UUID] = part
	return nil
}

func (r *InMemoryPartRepository) FindByUUID(uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, exists := r.parts[uuid]
	if !exists {
		return nil, ErrPartNotFound
	}
	return part, nil
}

func (r *InMemoryPartRepository) FindAll() ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := make([]*model.Part, 0, len(r.parts))
	for _, part := range r.parts {
		parts = append(parts, part)
	}
	return parts, nil
}

func (r *InMemoryPartRepository) Update(part *model.Part) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.parts[part.UUID]; !exists {
		return ErrPartNotFound
	}
	r.parts[part.UUID] = part
	return nil
}