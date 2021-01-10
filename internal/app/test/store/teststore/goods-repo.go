package teststore

import (
	"strings"

	"test-avito-merchant-experience/internal/app/test/store/model"
)

type GoodsRepo struct {
	goods map[int64]*model.Goods

	counter int64
}

func (r *GoodsRepo) GetListGoods(merchantID, offerID int64, name string) ([]*model.Goods, error) {
	goods := make([]*model.Goods, 0)
	for _, v := range r.goods {
		if v.MerchantID == merchantID && v.OfferID == offerID && strings.HasPrefix(v.Name, name) {
			goods = append(goods, v)
		}
	}
	return goods, nil
}

func (r *GoodsRepo) FindByMerchantIDAndOfferID(merchantID, offerID int64) (*model.Goods, error) {
	for _, v := range r.goods {
		if v.MerchantID == merchantID && v.OfferID == offerID {
			return v, nil
		}
	}
	return nil, nil
}

func (r *GoodsRepo) Update(m *model.Goods) error {
	r.counter++
	r.goods[r.counter] = m
	return nil
}

func (r *GoodsRepo) Save(m *model.Goods) error {
	r.counter++
	m.ID = r.counter
	r.goods[r.counter] = m
	return nil
}

func (r *GoodsRepo) Delete(merchantID, offerID int64) error {
	for k, v := range r.goods {
		if v.MerchantID == merchantID && v.OfferID == offerID {
			delete(r.goods, k)
		}
	}
	return nil
}
