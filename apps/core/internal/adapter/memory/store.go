// Package memory provides in-memory implementations of port.* repositories.
// It's used by the CLI and by tests; no persistence between runs.
package memory

import (
	"sync"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
)

// Store is the shared in-memory store. Safe for concurrent use.
type Store struct {
	mu sync.RWMutex

	projects   map[uuid.UUID]model.Project
	specs      map[uuid.UUID]model.Spec
	operations map[uuid.UUID]model.Operation
	runs       map[uuid.UUID]model.Run
	runResults map[uuid.UUID]model.StoredRunResult
	comments   map[uuid.UUID]model.Comment
	artifacts  map[uuid.UUID]model.Artifact
}

func NewStore() *Store {
	return &Store{
		projects:   make(map[uuid.UUID]model.Project),
		specs:      make(map[uuid.UUID]model.Spec),
		operations: make(map[uuid.UUID]model.Operation),
		runs:       make(map[uuid.UUID]model.Run),
		runResults: make(map[uuid.UUID]model.StoredRunResult),
		comments:   make(map[uuid.UUID]model.Comment),
		artifacts:  make(map[uuid.UUID]model.Artifact),
	}
}
