package scheduler

import (
	"context"
	"fmt"
	"wms-server/constants"
	"wms-server/helpers"
	"wms-server/usecases/v1"
	"time"

	"github.com/astaxie/beego/toolbox"
	logs "github.com/sirupsen/logrus"
)

// Scheduler ...
type (
	scheduler struct {
		Usecaase usecases.Usecase
		Logs     *logs.Logger
	}
	Scheduler interface {
		StartScheduler()
	}
)

// InitializeScheduler ...
func InitializeScheduler(
	u usecases.Usecase,
	l *logs.Logger,
) Scheduler {
	return &scheduler{
		Usecaase: u,
		Logs:     l,
	}
}

// StartScheduller ...
func (s *scheduler) StartScheduler() {
	// task := helpers.GetEnv("SCHEDULER_TASK")
	// timeTask := helpers.GetEnv("SCHEDULER_TIME_TASK")
	ctx := context.WithValue(context.Background(), constants.SPAN_ID, fmt.Sprintf("%d-%d", constants.SPAN_ID_LOGIN, time.Now().UnixNano()))
	// funcTask := toolbox.NewTask(task, timeTask, func() error {
	// 	s.Logs.WithContext(ctx).WithFields(logs.Fields{"Message": "Task"}).Info(task)
	// 	return s.Usecaase.ExScheduller(ctx)
	// })

	// toolbox.AddTask(task, funcTask)
	toolbox.StartTask()

	s.Logs.WithContext(ctx).WithFields(logs.Fields{"Message": "Start Toolbox"}).Info(helpers.GetCaller())
}
