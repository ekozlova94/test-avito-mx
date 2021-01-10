package forms

import "test-avito-merchant-experience/internal/app/test/store/model"

type Goods struct {
	OfferID   int64
	Name      string
	Price     int32
	Quantity  int32
	Available bool
}

func NewGoods(goods []*model.Goods) []*Goods {
	result := make([]*Goods, 0, len(goods))
	for _, v := range goods {
		r := &Goods{
			OfferID:   v.OfferID,
			Name:      v.Name,
			Price:     v.Price,
			Quantity:  v.Quantity,
			Available: v.Available,
		}
		result = append(result, r)
	}
	return result
}
