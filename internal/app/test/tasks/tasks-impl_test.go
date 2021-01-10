package tasks

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"test-avito-merchant-experience/internal/app/test/constants"
	"test-avito-merchant-experience/internal/app/test/getter/testgetter"
	"test-avito-merchant-experience/internal/app/test/store/model"
	"test-avito-merchant-experience/internal/app/test/store/teststore"
)

func Test_BackgroundTaskSuccess(t *testing.T) {
	//Arrange
	st := teststore.New()
	logger, _ := zap.NewDevelopment()
	getterSvc := testgetter.NewGetter()
	taskSvc := NewTasks(st, logger, getterSvc)

	var merchantID int64 = 129
	var url = "../../../../example-goods.xlsx"
	m := model.NewTask(merchantID, url)
	_ = st.Tasks().Save(m)

	//Act
	taskSvc.executeTask()

	//Assert
	result, _ := st.Tasks().GetTasksByMerchantIDAndUrlAndStatus(merchantID, url, "")

	require.Equal(t, m.ID, result.ID)
	require.Equal(t, constants.StatusSuccess, string(result.Status))
}

func Test_BackgroundTaskFailedToDownload(t *testing.T) {
	//Arrange
	st := teststore.New()
	logger, _ := zap.NewDevelopment()
	getterSvc := testgetter.NewGetter()
	taskSvc := NewTasks(st, logger, getterSvc)

	var merchantID int64 = 129
	var url = "example-goods.xlsx"
	m := model.NewTask(merchantID, url)
	_ = st.Tasks().Save(m)

	//Act
	taskSvc.executeTask()

	//Assert
	result, _ := st.Tasks().GetTasksByMerchantIDAndUrlAndStatus(merchantID, url, "")

	require.Equal(t, m.ID, result.ID)
	require.Equal(t, constants.StatusFailed, string(result.Status))
}
