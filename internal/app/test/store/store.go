package store

type Store interface {
	Goods() GoodsRepo
	Tasks() TasksRepo
}
