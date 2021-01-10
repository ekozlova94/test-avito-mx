package sqlstore

import (
	"database/sql"

	"test-avito-merchant-experience/internal/app/test/store"
)

type Store struct {
	db        *sql.DB
	goodsRepo *GoodsRepo
	tasksRepo *TasksRepo
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
		goodsRepo: &GoodsRepo{
			db: db,
		},
		tasksRepo: &TasksRepo{
			db: db,
		},
	}
}

func (s *Store) Goods() store.GoodsRepo {
	return s.goodsRepo
}

func (s *Store) Tasks() store.TasksRepo {
	return s.tasksRepo
}
