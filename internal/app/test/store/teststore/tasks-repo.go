package teststore

import (
	"test-avito-merchant-experience/internal/app/test/constants"
	"test-avito-merchant-experience/internal/app/test/store/model"
)

type TasksRepo struct {
	tasks map[int64]*model.Tasks

	counter int64
}

func (r *TasksRepo) GetTasksByMerchantIDAndUrlAndStatus(merchantID int64, url string, status constants.Status) (*model.Tasks, error) {
	for _, v := range r.tasks {
		if v.MerchantID == merchantID && v.Url == url && v.Status == status {
			return v, nil
		}
	}
	return nil, nil
}

func (r *TasksRepo) GetTasksByStatus(created constants.Status) ([]*model.Tasks, error) {
	tasks := make([]*model.Tasks, 0)
	for _, v := range r.tasks {
		if v.Status == created {
			tasks = append(tasks, v)
		}
	}
	return tasks, nil
}

func (r *TasksRepo) GetStatistics(taskID int64) (*model.Tasks, error) {
	for _, v := range r.tasks {
		if v.ID == taskID {
			return v, nil
		}
	}
	return nil, nil
}

func (r *TasksRepo) Save(m *model.Tasks) error {
	if m.ID > 0 {
		return r.update(m)
	}
	return r.insert(m)
}

func (r *TasksRepo) update(m *model.Tasks) error {
	for _, v := range r.tasks {
		if v.ID == m.ID {
			v.UpdatedRows = m.UpdatedRows
			v.SavedRows = m.SavedRows
			v.DeletedRows = m.DeletedRows
			v.BadRows = m.BadRows
			v.Status = m.Status
			v.ID = m.ID
		}
	}
	return nil
}

func (r *TasksRepo) insert(m *model.Tasks) error {
	if len(r.tasks) == 0 {
		m.ID = 1
		r.tasks[0] = m
		return nil
	}
	r.tasks[int64(len(r.tasks)-1)] = m
	m.ID = r.tasks[int64(len(r.tasks)-2)].ID + 1
	return nil
}

func (r *TasksRepo) Delete(id int64) error {
	for k, v := range r.tasks {
		if v.ID == id {
			delete(r.tasks, k)
		}
	}
	return nil
}
