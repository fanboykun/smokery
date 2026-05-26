package port

type EventType string

const (
	EventRunStarted  EventType = "run_started"
	EventStepDone    EventType = "step_done"
	EventRunFinished EventType = "run_finished"
)

type Event struct {
	Type  EventType `json:"type"`
	RunID string    `json:"run_id"`
	Data  any       `json:"data,omitempty"`
}

// EventBus is a simple pub/sub for run events.
type EventBus interface {
	Subscribe(runID string) <-chan Event
	Unsubscribe(runID string, ch <-chan Event)
	Publish(event Event)
}
