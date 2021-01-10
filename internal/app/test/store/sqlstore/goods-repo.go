package sqlstore

import (
	"database/sql"
	"fmt"

	"test-avito-merchant-experience/internal/app/test/store/model"
	"test-avito-merchant-experience/internal/app/test/utils"
)

type GoodsRepo struct {
	db *sql.DB
}

func (r *GoodsRepo) GetListGoods(merchantID, offerID int64, name string) ([]*model.Goods, error) {
	query, args := utils.BuildQueryByParams(merchantID, offerID, name)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	//noinspection GoUnhandledErrorResult
	defer rows.Close()

	result := make([]*model.Goods, 0)
	for rows.Next() {
		goods := new(model.Goods)
		err := rows.Scan(
			&goods.ID, &goods.MerchantID, &goods.OfferID, &goods.Name, &goods.Price, &goods.Quantity, &goods.Available,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, goods)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *GoodsRepo) FindByMerchantIDAndOfferID(merchantID, offerID int64) (*model.Goods, error) {
	var goods model.Goods
	err := r.db.QueryRow(
		"SELECT * FROM goods WHERE merchant_id = $1 AND offer_id = $2",
		merchantID,
		offerID,
	).Scan(&goods.ID, &goods.MerchantID, &goods.OfferID, &goods.Name, &goods.Price, &goods.Quantity, &goods.Available)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &goods, nil
}

func (r *GoodsRepo) Save(m *model.Goods) error {
	return r.db.QueryRow(
		"INSERT INTO goods ("+
			"merchant_id, offer_id, price, name, quantity, available) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		&m.MerchantID,
		&m.OfferID,
		&m.Price,
		&m.Name,
		&m.Quantity,
		&m.Available,
	).Scan(&m.ID)
}

func (r *GoodsRepo) Update(m *model.Goods) error {
	result, err := r.db.Exec(
		"UPDATE goods "+
			"SET price = $1, "+
			"name = $2, "+
			"quantity = $3, "+
			"available = $4 "+
			"WHERE merchant_id = $5 AND offer_id = $6;",
		&m.Price,
		&m.Name,
		&m.Quantity,
		&m.Available,
		&m.MerchantID,
		&m.OfferID,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return fmt.Errorf("updated records %d instead of 1", affected)
	}
	return nil
}

func (r *GoodsRepo) Delete(merchantID, offerID int64) error {
	_, err := r.db.Exec(
		"DELETE FROM goods WHERE merchant_id = $1 AND offer_id = $2",
		&merchantID,
		&offerID,
	)
	return err
}
