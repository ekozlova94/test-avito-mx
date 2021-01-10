package store

import (
	"test-avito-merchant-experience/internal/app/test/constants"
	"test-avito-merchant-experience/internal/app/test/store/model"
)

type TasksRepo interface {
	GetTasksByMerchantIDAndUrlAndStatus(merchantID int64, url string, status constants.Status) (*model.Tasks, error)
	GetTasksByStatus(created constants.Status) ([]*model.Tasks, error)
	GetStatistics(taskID int64) (*model.Tasks, error)
	Save(task *model.Tasks) error
	Delete(id int64) error
}
