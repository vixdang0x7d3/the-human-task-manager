package mock

import (
	"context"

	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
)

var _ domain.TaskService = (*TaskService)(nil)

type TaskService struct {
	fnCreateTask func(ctx context.Context, cmd domain.CreateTaskCmd) (domain.Task, error)
}

func (s *TaskService) CreateTask(ctx context.Context, cmd domain.CreateTaskCmd) (domain.Task, error) {
	return s.fnCreateTask(ctx, cmd)
}
