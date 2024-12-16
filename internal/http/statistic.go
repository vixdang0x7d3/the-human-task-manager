package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/generic"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/models"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/components"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/pages"
)

func (s *Server) registerStatisticRoutes(r *echo.Group) {
	r.GET("/statistic", s.handleStatisticShow)

	r.POST("/statistic", s.handleStatisticDayChange)
}

func (s *Server) handleStatisticShow(c echo.Context) error {

	// tasks, total, err := s.TaskItemService.Find(c.Request().Context(), domain.TaskItemFilter{

	// })

	taskTimes := []models.TaskTimeView{
		{
			TaskItem: models.TaskItemView{
				Description: "Task1",
			},
			TaskTime: 20,
		},
		{
			TaskItem: models.TaskItemView{
				Description: "Task2",
			},
			TaskTime: 10,
		},
		{
			TaskItem: models.TaskItemView{
				Description: "Task3",
			},
			TaskTime: 30,
		},
	}

	descriptions := generic.Map(taskTimes, func(t models.TaskTimeView) string {
		return t.TaskItem.Description
	})

	taskTimeViews := generic.Map(taskTimes, func(t models.TaskTimeView) int {
		return t.TaskTime
	})

	return render(c, http.StatusOK, pages.Statistics(
		descriptions,
		taskTimeViews,
		"/logout"))
}

func (s *Server) handleStatisticDayChange(c echo.Context) error {

	taskTimes := []models.TaskTimeView{
		{
			TaskItem: models.TaskItemView{
				Description: "Task4",
			},
			TaskTime: 20,
		},
		{
			TaskItem: models.TaskItemView{
				Description: "Task5",
			},
			TaskTime: 15,
		},
		{
			TaskItem: models.TaskItemView{
				Description: "Task6",
			},
			TaskTime: 70,
		},
	}

	descriptions := generic.Map(taskTimes, func(t models.TaskTimeView) string {
		return t.TaskItem.Description
	})

	taskTimeViews := generic.Map(taskTimes, func(t models.TaskTimeView) int {
		return t.TaskTime
	})

	return render(c, http.StatusOK, components.StatisticChart(
		descriptions,
		taskTimeViews))
}
