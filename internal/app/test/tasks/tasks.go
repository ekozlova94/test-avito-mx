package tasks

import "test-avito-merchant-experience/internal/app/test/store/model"

type Tasks interface {
	GetStatisticsByTask(taskID int64) (*model.Tasks, error)
	CreateTask(merchantID int64, url string) (int64, error)
	StartBackgroundTask()
}
