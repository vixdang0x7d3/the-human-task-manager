package mock

import (
	"context"

	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
)

var _ domain.TaskService = (*TaskService)(nil)

type TaskService struct {
	fnCreate   func(ctx context.Context, cmd domain.CreateTaskCmd) (domain.Task, error)
	fnDelete   func(ctx context.Context, id string) (domain.Task, error)
	fnUpdate   func(ctx context.Context, id string, cmd domain.UpdateTaskCmd) (domain.Task, error)
	fnComplete func(ctx context.Context, id string) (domain.Task, error)
}

func (s *TaskService) Create(ctx context.Context, cmd domain.CreateTaskCmd) (domain.Task, error) {
	return s.fnCreate(ctx, cmd)
}

func (s *TaskService) Delete(ctx context.Context, id string) (domain.Task, error) {
	return s.fnDelete(ctx, id)
}

func (s *TaskService) Update(ctx context.Context, id string, cmd domain.UpdateTaskCmd) (domain.Task, error) {
	return s.fnUpdate(ctx, id, cmd)
}

func (s *TaskService) Complete(ctx context.Context, id string) (domain.Task, error) {
	return s.fnComplete(ctx, id)
}
