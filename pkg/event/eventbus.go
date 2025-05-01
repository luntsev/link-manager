package event

const (
	EventLinkGet int = iota
	EventLinkCreate
	EventLinkUpdate
	EventLinkDelete
	EventListGet
	EventRegisterUser
	EventLoginUser
)

type Event struct {
	Type  int
	Event any
}

type EventBus struct {
	bus chan Event
}

func (eb *EventBus) Pubish(e Event) {
	eb.bus <- e
}

func (eb *EventBus) Subscribe() <-chan Event {
	return eb.bus
}

func NewEventBus() *EventBus {
	return &EventBus{
		bus: make(chan Event),
	}
}
