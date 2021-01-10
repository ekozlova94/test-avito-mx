package model

type Goods struct {
	ID         int64
	MerchantID int64
	OfferID    int64
	Name       string
	Price      int32
	Quantity   int32
	Available  bool
}
