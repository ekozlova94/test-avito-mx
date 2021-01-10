package forms

import (
	"test-avito-merchant-experience/internal/app/test/constants"
	"test-avito-merchant-experience/internal/app/test/store/model"
)

type Statistics struct {
	Status      constants.Status `json:"Status"`
	TotalRows   int32            `json:"Total"`
	UpdatedRows int32            `json:"Updated"`
	SavedRows   int32            `json:"Saved"`
	DeletedRows int32            `json:"Deleted"`
	BadRows     int32            `json:"Bad"`
}

func NewStatistics(m *model.Tasks) *Statistics {
	return &Statistics{
		Status:      m.Status,
		TotalRows:   m.TotalRows,
		UpdatedRows: m.UpdatedRows,
		SavedRows:   m.SavedRows,
		DeletedRows: m.DeletedRows,
		BadRows:     m.BadRows,
	}
}
