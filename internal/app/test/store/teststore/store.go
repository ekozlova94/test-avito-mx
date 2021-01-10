package teststore

import (
	"test-avito-merchant-experience/internal/app/test/store"
	"test-avito-merchant-experience/internal/app/test/store/model"
)

type Store struct {
	goodsRepo *GoodsRepo
	tasksRepo *TasksRepo
}

func New() *Store {
	return &Store{
		goodsRepo: &GoodsRepo{
			goods: make(map[int64]*model.Goods, 0),
		},
		tasksRepo: &TasksRepo{
			tasks: make(map[int64]*model.Tasks, 0),
		},
	}
}

func (s *Store) Goods() store.GoodsRepo {
	return s.goodsRepo
}

func (s *Store) Tasks() store.TasksRepo {
	return s.tasksRepo
}
