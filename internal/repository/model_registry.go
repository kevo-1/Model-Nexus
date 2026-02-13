package repository

import (
	"sync"

	"github.com/kevo-1/model-serving-platform/internal/domain"
	"github.com/kevo-1/model-serving-platform/internal/metrics"
)


type ModelRegistry  struct {
	mu sync.RWMutex
	models map[string]domain.ModelPredictor
}


func NewModelRegistery() *ModelRegistry  {
	return &ModelRegistry {
		models: make(map[string]domain.ModelPredictor),
	}
}

func (r *ModelRegistry ) Get(id string) (domain.ModelPredictor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	model, ok := r.models[id]
	if !ok {
		return nil, &domain.ModelNotFoundError{ModelID: id}
	}
	return model, nil
}


func (r *ModelRegistry) Register(id string, model domain.ModelPredictor) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.models[id]; exists {
        return &domain.ModelAlreadyExistsError{ModelID: id}
    }

    r.models[id] = model
	metrics.SetModelsLoaded(len(r.models))
    
	return nil
}


func (r *ModelRegistry) List() []string {
    r.mu.RLock()
    defer r.mu.RUnlock()

    ids := make([]string, 0, len(r.models))
    for id := range r.models {
        ids = append(ids, id)
    }
    return ids
}


func (r *ModelRegistry ) Count() int {
	r.mu.RLock()
	defer r.mu.Unlock()

	return len(r.models)
}