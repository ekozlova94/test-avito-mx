package utils

import (
	"strconv"
	"strings"

	"test-avito-merchant-experience/internal/app/test/constants"
)

func ParseAndValidateUrlAndMerchantID(valueUrl string, valueMerchantID string) (string, int64, error) {
	url := strings.Trim(valueUrl, " ")
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		return "", -1, constants.ErrIncorrect
	}
	merchantID, err := strconv.ParseInt(valueMerchantID, 10, 64)
	if err != nil || merchantID < 0 {
		return "", -1, constants.ErrNotNumber
	}
	return url, merchantID, nil
}

func BuildQueryByParams(merchantID int64, offerID int64, name string) (string, []interface{}) {
	if merchantID < 0 && offerID < 0 && name == "" {
		return "SELECT * FROM goods", nil
	}
	query := "SELECT * FROM goods WHERE"
	i := 1
	args := make([]interface{}, 0, 3)
	if merchantID > 0 {
		query += " merchant_id = $" + strconv.Itoa(i)
		args = append(args, merchantID)
		i++
	}
	if offerID > 0 {
		if i > 1 {
			query += " AND"
		}
		query += " offer_id = $" + strconv.Itoa(i)
		args = append(args, offerID)
		i++
	}
	if name != "" {
		if i > 1 {
			query += " AND"
		}
		query += " name LIKE '%' || $" + strconv.Itoa(i) + " || '%'"
		args = append(args, name)
	}
	return query, args
}
