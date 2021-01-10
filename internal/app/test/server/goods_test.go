package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"test-avito-merchant-experience/internal/app/test/getter/prodgetter"
	"test-avito-merchant-experience/internal/app/test/store/model"
	"test-avito-merchant-experience/internal/app/test/store/teststore"
	"test-avito-merchant-experience/internal/app/test/tasks"
	"test-avito-merchant-experience/internal/pkg/forms"
)

func Test_EmptyDB(t *testing.T) {
	//Arrange
	st := teststore.New()
	logger, _ := zap.NewDevelopment()
	getterSvc := prodgetter.NewGetter(logger)
	taskSvc := tasks.NewTasks(st, logger, getterSvc)
	s := NewServer(st, logger, taskSvc, getterSvc)
	var merchantID int64 = 129
	var url = "http://localhost:8090/api/load-file"

	//Act
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/goods?"+
			"url="+url+"&"+
			"merchant-id="+strconv.FormatInt(merchantID, 10),
		nil,
	)
	s.ServeHTTP(rec, req)

	//Assert
	assert.Equal(t, 200, rec.Code)

	task, _ := st.Tasks().GetTasksByMerchantIDAndUrlAndStatus(merchantID, url, "")

	require.NotNil(t, task)
}

func Test_AdditionalCreateTaskToSameMerchantIDAndUrl(t *testing.T) {
	//Arrange
	st := teststore.New()
	logger, _ := zap.NewDevelopment()
	getterSvc := prodgetter.NewGetter(logger)
	taskSvc := tasks.NewTasks(st, logger, getterSvc)
	s := NewServer(st, logger, taskSvc, getterSvc)
	var merchantID int64 = 129
	var url = "http://localhost:8090/api/load-file"
	m := model.NewTask(merchantID, url)
	_ = st.Tasks().Save(m)

	//Act
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/goods?"+
			"url="+url+"&"+
			"merchant-id="+strconv.FormatInt(merchantID, 10),
		nil,
	)
	s.ServeHTTP(rec, req)

	//Assert
	assert.Equal(t, 500, rec.Code)

	err := rec.Body.String()

	require.EqualValues(t, "\"such a task has already been created\"", err)
}

func Test_GetGoodsListSuccess(t *testing.T) {
	//Arrange
	st := teststore.New()

	logger, _ := zap.NewDevelopment()
	getterSvc := prodgetter.NewGetter(logger)
	taskSvc := tasks.NewTasks(st, logger, getterSvc)
	s := NewServer(st, logger, taskSvc, getterSvc)
	var merchantID, offerID int64 = 1, 1
	namePrefix := "теле"
	name := "телевизор"
	var price, quantity int32 = 100, 10
	m := &model.Goods{
		ID:         0,
		MerchantID: merchantID,
		OfferID:    offerID,
		Name:       name,
		Price:      price,
		Quantity:   quantity,
		Available:  true,
	}
	_ = st.Goods().Save(m)

	//Act
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/list?"+
			"merchant-id="+strconv.FormatInt(merchantID, 10)+"&"+
			"offer-id="+strconv.FormatInt(offerID, 10)+"&"+
			"name="+namePrefix,
		nil,
	)
	s.ServeHTTP(rec, req)

	//Assert
	assert.Equal(t, 200, rec.Code)

	var goods []*forms.Goods
	_ = json.NewDecoder(rec.Body).Decode(&goods)

	require.EqualValues(t, offerID, goods[0].OfferID)
	require.EqualValues(t, name, goods[0].Name)
	require.EqualValues(t, price, goods[0].Price)
	require.EqualValues(t, quantity, goods[0].Quantity)
	require.EqualValues(t, true, goods[0].Available)
}

func Test_GetGoodsListFailed(t *testing.T) {
	//Arrange
	st := teststore.New()
	logger, _ := zap.NewDevelopment()
	getterSvc := prodgetter.NewGetter(logger)
	taskSvc := tasks.NewTasks(st, logger, getterSvc)
	s := NewServer(st, logger, taskSvc, getterSvc)
	var merchantID, offerID int64 = 222, 2
	namePrefix := "теле"

	//Act
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/list?"+
			"merchant-id="+strconv.FormatInt(merchantID, 10)+"&"+
			"offer-id="+strconv.FormatInt(offerID, 10)+"&"+
			"name="+namePrefix,
		nil,
	)
	s.ServeHTTP(rec, req)

	//Assert
	assert.Equal(t, 404, rec.Code)

	var goods []*forms.Goods
	_ = json.NewDecoder(rec.Body).Decode(&goods)

	require.Nil(t, goods)
}
