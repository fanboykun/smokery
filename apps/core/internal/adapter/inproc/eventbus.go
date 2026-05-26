package inproc

import (
	"sync"

	"github.com/fanboykun/smokery/apps/core/internal/port"
)

// EventBus is a simple in-memory pub/sub implementing port.EventBus.
type EventBus struct {
	mu   sync.RWMutex
	subs map[string][]chan port.Event
}

func NewEventBus() *EventBus {
	return &EventBus{subs: make(map[string][]chan port.Event)}
}

var _ port.EventBus = (*EventBus)(nil)

func (b *EventBus) Subscribe(runID string) <-chan port.Event {
	ch := make(chan port.Event, 64)
	b.mu.Lock()
	b.subs[runID] = append(b.subs[runID], ch)
	b.mu.Unlock()
	return ch
}

func (b *EventBus) Unsubscribe(runID string, ch <-chan port.Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	subs := b.subs[runID]
	for i, s := range subs {
		if s == ch {
			b.subs[runID] = append(subs[:i], subs[i+1:]...)
			close(s)
			return
		}
	}
}

func (b *EventBus) Publish(event port.Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subs[event.RunID] {
		select {
		case ch <- event:
		default:
		}
	}
}
