package stat

import (
	"link-manager/pkg/event"
	"log"
)

type StatServiceDeps struct {
	StatRepository *StatRepository
	EventBus       *event.EventBus
}

type StatService struct {
	StatRepository *StatRepository
	EventBus       *event.EventBus
}

func NewStatService(deps StatServiceDeps) *StatService {
	return &StatService{
		StatRepository: deps.StatRepository,
		EventBus:       deps.EventBus,
	}
}

func (s *StatService) AddClick() {
	for msg := range s.EventBus.Subscribe() {
		if msg.Type == event.EventLinkGet {
			id, ok := msg.Event.(uint)
			if !ok {
				log.Fatalln("Bad EventLinkGet:", msg.Event)
				continue
			}
			s.StatRepository.AddClick(id)
		}
	}
}
