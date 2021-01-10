package store

import "test-avito-merchant-experience/internal/app/test/store/model"

type GoodsRepo interface {
	GetListGoods(merchantID, offerID int64, name string) ([]*model.Goods, error)
	FindByMerchantIDAndOfferID(merchantID, offerID int64) (*model.Goods, error)
	Update(m *model.Goods) error
	Save(goods *model.Goods) error
	Delete(merchantID, offerID int64) error
}
