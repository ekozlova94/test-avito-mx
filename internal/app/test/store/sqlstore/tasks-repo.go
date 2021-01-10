package sqlstore

import (
	"database/sql"
	"fmt"

	"test-avito-merchant-experience/internal/app/test/constants"
	"test-avito-merchant-experience/internal/app/test/store/model"
)

type TasksRepo struct {
	db *sql.DB
}

func (r *TasksRepo) GetTasksByMerchantIDAndUrlAndStatus(
	merchantID int64, url string, status constants.Status,
) (*model.Tasks, error) {
	var tasks model.Tasks
	err := r.db.QueryRow(
		"SELECT * FROM tasks "+
			"WHERE merchant_id = $1 AND url = $2 AND status = $3",
		merchantID,
		url,
		status,
	).Scan(
		&tasks.ID,
		&tasks.CreatedAt,
		&tasks.Status,
		&tasks.MerchantID,
		&tasks.Url,
		&tasks.TotalRows,
		&tasks.UpdatedRows,
		&tasks.SavedRows,
		&tasks.DeletedRows,
		&tasks.BadRows,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &tasks, nil
}

func (r *TasksRepo) GetTasksByStatus(created constants.Status) ([]*model.Tasks, error) {
	rows, err := r.db.Query(
		"SELECT * FROM tasks "+
			"WHERE status = $1 "+
			"ORDER BY created_at",
		created,
	)
	if err != nil {
		return nil, err
	}
	//noinspection GoUnhandledErrorResult
	defer rows.Close()

	result := make([]*model.Tasks, 0)
	for rows.Next() {
		tasks := new(model.Tasks)
		err := rows.Scan(&tasks.ID,
			&tasks.CreatedAt,
			&tasks.Status,
			&tasks.MerchantID,
			&tasks.Url,
			&tasks.TotalRows,
			&tasks.UpdatedRows,
			&tasks.SavedRows,
			&tasks.DeletedRows,
			&tasks.BadRows,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, tasks)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *TasksRepo) GetStatistics(taskID int64) (*model.Tasks, error) {
	var tasks model.Tasks
	err := r.db.QueryRow(
		"SELECT * FROM tasks "+
			"WHERE id = $1",
		taskID,
	).Scan(
		&tasks.ID,
		&tasks.CreatedAt,
		&tasks.Status,
		&tasks.MerchantID,
		&tasks.Url,
		&tasks.TotalRows,
		&tasks.UpdatedRows,
		&tasks.SavedRows,
		&tasks.DeletedRows,
		&tasks.BadRows,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &tasks, nil
}

func (r *TasksRepo) Save(m *model.Tasks) error {
	if m.ID > 0 {
		return r.update(m)
	}
	return r.insert(m)
}

func (r *TasksRepo) update(m *model.Tasks) error {
	result, err := r.db.Exec(
		"UPDATE tasks "+
			"SET updated_rows = $1, "+
			"total_rows = $2, "+
			"saved_rows = $3, "+
			"deleted_rows = $4, "+
			"bad_rows = $5, "+
			"status = $6 "+
			"WHERE id = $7;",
		&m.UpdatedRows,
		&m.TotalRows,
		&m.SavedRows,
		&m.DeletedRows,
		&m.BadRows,
		&m.Status,
		&m.ID,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return fmt.Errorf("updated records %d instead of 1", affected)
	}
	return nil
}

func (r *TasksRepo) insert(m *model.Tasks) error {
	return r.db.QueryRow(
		"INSERT INTO tasks ("+
			"created_at, status, merchant_id, url, total_rows, updated_rows, saved_rows, deleted_rows, bad_rows"+
			") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		&m.CreatedAt,
		&m.Status,
		&m.MerchantID,
		&m.Url,
		&m.TotalRows,
		&m.UpdatedRows,
		&m.SavedRows,
		&m.DeletedRows,
		&m.BadRows,
	).Scan(&m.ID)
}

func (r *TasksRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE id = $1", &id)
	return err
}
