package model

import (
	"time"

	"test-avito-merchant-experience/internal/app/test/constants"
)

type Tasks struct {
	ID          int64
	CreatedAt   time.Time
	Status      constants.Status
	MerchantID  int64
	Url         string
	TotalRows   int32
	UpdatedRows int32
	SavedRows   int32
	DeletedRows int32
	BadRows     int32
}

func NewTask(merchantID int64, url string) *Tasks {
	return &Tasks{
		CreatedAt:  time.Now(),
		Status:     constants.StatusCreated,
		MerchantID: merchantID,
		Url:        url,
	}
}
