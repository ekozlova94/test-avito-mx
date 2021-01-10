package tasks

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/tealeg/xlsx/v3"
	"go.uber.org/zap"
	"test-avito-merchant-experience/internal/app/test/constants"
	"test-avito-merchant-experience/internal/app/test/getter"
	"test-avito-merchant-experience/internal/app/test/store"
	"test-avito-merchant-experience/internal/app/test/store/model"
)

type TasksImpl struct {
	log    *zap.Logger
	store  store.Store
	getter getter.Getter
}

func NewTasks(store store.Store, log *zap.Logger, getter getter.Getter) *TasksImpl {
	return &TasksImpl{
		log:    log,
		store:  store,
		getter: getter,
	}
}

func (s *TasksImpl) CreateTask(merchantID int64, url string) (int64, error) {
	task, err := s.store.Tasks().GetTasksByMerchantIDAndUrlAndStatus(merchantID, url, constants.StatusCreated)
	if err != nil {
		return -1, err
	}
	if task != nil {
		return -1, constants.ErrCreateTask
	}
	task = model.NewTask(merchantID, url)
	if err := s.store.Tasks().Save(task); err != nil {
		return -1, err
	}
	s.log.Info("Task created successfully", zap.Int64("id", task.ID))
	return task.ID, nil
}

func (s *TasksImpl) GetStatisticsByTask(taskID int64) (*model.Tasks, error) {
	return s.store.Tasks().GetStatistics(taskID)
}

func (s *TasksImpl) StartBackgroundTask() {
	go func() {
		for {
			s.executeTask()
			time.Sleep(5 * time.Second)
		}
	}()
}

func (s *TasksImpl) executeTask() {
	tasks, err := s.store.Tasks().GetTasksByStatus(constants.StatusCreated)
	if err != nil {
		s.log.Error("Error while getting task", zap.Error(err))
		return
	}
	if tasks != nil {
		for _, task := range tasks {
			s.downloadFileAndUpdateGoods(task)
			if errSave := s.store.Tasks().Save(task); errSave != nil {
				s.log.Error("Error while saving task status", zap.Error(errSave))
				continue
			}
		}
	}
}

func (s *TasksImpl) downloadFileAndUpdateGoods(task *model.Tasks) {
	task.Status = constants.StatusInProgress
	if errSave := s.store.Tasks().Save(task); errSave != nil {
		s.log.Error("Error while saving task status", zap.Error(errSave))
		return
	}
	file, err := ioutil.TempFile("", "goods.*.xlsx")
	if err != nil {
		s.log.Error("Error while creating temp file", zap.Error(err))
		return
	}
	fileName := file.Name()

	//noinspection GoUnhandledErrorResult
	defer os.Remove(fileName)

	if err := s.getter.DownloadFile(task.Url, file); err != nil {
		task.Status = constants.StatusFailed
		s.log.Error("Error while downloading file", zap.Error(err))
		return
	}
	if err := s.UpdateGoods(task, fileName); err != nil {
		task.UpdatedRows = 0
		task.SavedRows = 0
		task.DeletedRows = 0
		task.BadRows = 0
		s.log.Error("Error while updating file", zap.Error(err))
		return
	}
	task.Status = constants.StatusSuccess
	s.log.Info("Task executed successfully", zap.Int64("id", task.ID))
	return
}

func (s *TasksImpl) UpdateGoods(task *model.Tasks, filename string) error {
	wb, err := xlsx.OpenFile(filename)
	if err != nil {
		return err
	}
	for _, sh := range wb.Sheets {
		task.TotalRows += int32(sh.MaxRow)
	}
	for _, sh := range wb.Sheets {
		for i := 0; i < sh.MaxRow; i++ {
			if err := s.processRow(task, sh, i); err != nil {
				return err
			}
			if errSave := s.store.Tasks().Save(task); errSave != nil {
				return errSave
			}
		}
		sh.Close()
	}
	s.log.Info("Goods updated successfully", zap.Int64("task id", task.ID), zap.Int32("total rows", task.TotalRows))
	return nil
}

func (s *TasksImpl) processRow(task *model.Tasks, sh *xlsx.Sheet, i int) error {
	goods, errPars := parsRow(sh, i)
	if errPars != nil {
		task.BadRows += 1
		return nil
	}
	result, err := s.store.Goods().FindByMerchantIDAndOfferID(task.MerchantID, goods.OfferID)
	if err != nil {
		return err
	}
	if result != nil {
		if goods.Available {
			result.Price = goods.Price
			result.Name = goods.Name
			result.Quantity = goods.Quantity
			result.Available = goods.Available
			if err := s.store.Goods().Update(result); err != nil {
				return err
			}
			task.UpdatedRows += 1
			return nil
		}
		if errDelete := s.store.Goods().Delete(result.MerchantID, result.OfferID); errDelete != nil {
			return errDelete
		}
		task.DeletedRows += 1
		return nil
	}
	if !goods.Available {
		return nil
	}
	goods.MerchantID = task.MerchantID
	if errSave := s.store.Goods().Save(goods); errSave != nil {
		return errSave
	}
	task.SavedRows += 1
	return nil
}

func parsRow(sh *xlsx.Sheet, i int) (*model.Goods, error) {
	valueOfferID, err := sh.Cell(i, 0)
	if err != nil {
		return nil, err
	}
	offerID, err := validate(valueOfferID)
	if err != nil {
		return nil, err
	}
	valueName, err := sh.Cell(i, 1)
	if err != nil {
		return nil, err
	}
	if valueName == nil {
		return nil, fmt.Errorf("name field is empty")
	}
	valuePrice, err := sh.Cell(i, 2)
	if err != nil {
		return nil, err
	}
	price, err := validate(valuePrice)
	if err != nil {
		return nil, err
	}
	valueQuantity, err := sh.Cell(i, 3)
	if err != nil {
		return nil, err
	}
	quantity, err := validate(valueQuantity)
	if err != nil {
		return nil, err
	}
	valueAvailable, err := sh.Cell(i, 4)
	if err != nil {
		return nil, err
	}
	if valueAvailable == nil {
		return nil, fmt.Errorf("available field is empty")
	}
	goods := &model.Goods{
		OfferID:   offerID,
		Name:      valueName.String(),
		Price:     int32(price),
		Quantity:  int32(quantity),
		Available: valueAvailable.Bool(),
	}
	return goods, nil
}

func validate(value *xlsx.Cell) (int64, error) {
	res, err := value.Int()
	if err != nil {
		return -1, err
	}
	if res < 0 {
		return -1, fmt.Errorf("number negative")
	}
	return int64(res), nil
}
