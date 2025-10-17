package event

import (
	"context"
	"event-api/internal/domain/model"

	"github.com/google/uuid"
)

type Repository interface {
	RetrieveById(ctx context.Context, id string) (*model.Event, error)
	RetrieveAll(ctx context.Context) ([]model.Event, error)
	Save(ctx context.Context, event *model.Event) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetEvent(context context.Context, id string) (*model.Event, error) {
	if id == "" {
		return nil, model.NewValidationErrorWithTime("invalid event id")
	}

	return s.repository.RetrieveById(context, id)
}

func (s *Service) GetAllEvents(context context.Context) ([]model.Event, error) {
	return s.repository.RetrieveAll(context)
}

func (s *Service) CreateEvent(context context.Context, event *model.Event) error {
	uuidNew, err := uuid.NewV6()
	if err != nil {
		return err
	}
	event.ID = uuidNew.String()
	return s.repository.Save(context, event)
}
