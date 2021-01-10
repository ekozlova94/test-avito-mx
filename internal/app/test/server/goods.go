package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"test-avito-merchant-experience/internal/app/test/constants"
	"test-avito-merchant-experience/internal/app/test/getter"
	"test-avito-merchant-experience/internal/app/test/store"
	"test-avito-merchant-experience/internal/app/test/store/model"
	"test-avito-merchant-experience/internal/app/test/tasks"
	"test-avito-merchant-experience/internal/app/test/utils"
	"test-avito-merchant-experience/internal/pkg/forms"
)

type server struct {
	Router *gin.Engine
	Store  store.Store
	Log    *zap.Logger
	Tasks  tasks.Tasks
	Getter getter.Getter
}

func NewServer(store store.Store, logger *zap.Logger, tasks tasks.Tasks, getter getter.Getter) *server {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	s := &server{
		Router: r,
		Store:  store,
		Log:    logger,
		Tasks:  tasks,
		Getter: getter,
	}
	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.Router.GET("/api/v1/goods", s.getLinkAndMerchantID)
	s.Router.GET("/api/v1/list", s.getGoodsList)
	s.Router.GET("/api/v1/progress", s.getProgress)
}

func (s *server) getLinkAndMerchantID(c *gin.Context) {
	valueUrl := c.Query("url")
	valueMerchantID := c.Query("merchant-id")

	url, merchantID, err := utils.ParseAndValidateUrlAndMerchantID(valueUrl, valueMerchantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("%s", err))
		return
	}
	taskID, err := s.Tasks.CreateTask(merchantID, url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	c.JSON(http.StatusOK, taskID)
	return
}

func (s *server) getGoodsList(c *gin.Context) {
	valueMerchantID := c.Query("merchant-id")
	valueOfferID := c.Query("offer-id")
	name := c.Query("name")

	var merchantID, offerID int64 = -1, -1
	var err error
	if valueMerchantID != "" {
		merchantID, err = strconv.ParseInt(valueMerchantID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("%s", err))
			return
		}
	}
	if valueOfferID != "" {
		offerID, err = strconv.ParseInt(valueOfferID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("%s", err))
			return
		}
	}
	list, err := s.getList(merchantID, offerID, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	if len(list) == 0 {
		c.JSON(http.StatusNotFound, constants.ErrNotFound)
		return
	}
	c.JSON(http.StatusOK, forms.NewGoods(list))
}

func (s *server) getList(merchantID, offerID int64, name string) ([]*model.Goods, error) {
	list, err := s.Store.Goods().GetListGoods(merchantID, offerID, name)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *server) getProgress(c *gin.Context) {
	valueTaskID := c.Query("task-id")
	taskID, err := strconv.ParseInt(valueTaskID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, constants.ErrNotNumber)
		return
	}
	task, err := s.Tasks.GetStatisticsByTask(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	if task == nil {
		c.JSON(http.StatusNotFound, constants.ErrNotFound)
		return
	}
	c.JSON(http.StatusOK, forms.NewStatistics(task))
}
