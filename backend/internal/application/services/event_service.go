package application

import (
	"backend/internal/application/interfaces"
	"backend/internal/domain/events"
)

type EventService struct {
	repo interfaces.EventRepository
}

func NewEventService(repo interfaces.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) SaveEvent(event events.Event) error {
	return s.repo.Save(event)
}
