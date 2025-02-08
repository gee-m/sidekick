package registry

import (
	"context"
	"fmt"
	"sync"
)

// Appgent defines the interface that all appgents must implement
type Appgent interface {
	ID() string
	Initialize(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// Registry manages all registered appgents
type Registry struct {
	mu       sync.RWMutex
	appgents map[string]Appgent
}

func New() *Registry {
	return &Registry{
		appgents: make(map[string]Appgent),
	}
}

func (r *Registry) Register(appgent Appgent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := appgent.ID()
	if _, exists := r.appgents[id]; exists {
		return fmt.Errorf("appgent with ID %s already registered", id)
	}

	r.appgents[id] = appgent
	return nil
}

func (r *Registry) Get(id string) (Appgent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	appgent, exists := r.appgents[id]
	if !exists {
		return nil, fmt.Errorf("appgent with ID %s not found", id)
	}

	return appgent, nil
}

func (r *Registry) List() []Appgent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]Appgent, 0, len(r.appgents))
	for _, appgent := range r.appgents {
		list = append(list, appgent)
	}
	return list
}
