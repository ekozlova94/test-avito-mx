package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"test-avito-merchant-experience/internal/app/test/constants"
)

func Test_BuildQueryByParams(t *testing.T) {
	testCases := []struct {
		testName   string
		merchantID int64
		offerID    int64
		name       string
		query      string
		args       []interface{}
	}{
		{
			testName:   "Correct_1",
			merchantID: -1,
			offerID:    -1,
			name:       "",
			query:      "SELECT * FROM goods",
			args:       nil,
		},
		{
			testName:   "Correct_2",
			merchantID: 1,
			offerID:    -1,
			name:       "",
			query:      "SELECT * FROM goods WHERE merchant_id = $1",
			args:       []interface{}{int64(1)},
		},
		{
			testName:   "Correct_3",
			merchantID: -1,
			offerID:    1,
			name:       "",
			query:      "SELECT * FROM goods WHERE offer_id = $1",
			args:       []interface{}{int64(1)},
		},
		{
			testName:   "Correct_4",
			merchantID: -1,
			offerID:    -1,
			name:       "теле",
			query:      "SELECT * FROM goods WHERE name LIKE '%' || $1 || '%'",
			args:       []interface{}{"теле"},
		},
		{
			testName:   "Correct_5",
			merchantID: 1,
			offerID:    1,
			name:       "",
			query:      "SELECT * FROM goods WHERE merchant_id = $1 AND offer_id = $2",
			args:       []interface{}{int64(1), int64(1)},
		},
		{
			testName:   "Correct_6",
			merchantID: 1,
			offerID:    1,
			name:       "теле",
			query:      "SELECT * FROM goods WHERE merchant_id = $1 AND offer_id = $2 AND name LIKE '%' || $3 || '%'",
			args:       []interface{}{int64(1), int64(1), "теле"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			query, args := BuildQueryByParams(tc.merchantID, tc.offerID, tc.name)
			require.Equal(t, tc.query, query)
			require.Equal(t, tc.args, args)
		})
	}
}

func Test_ParseAndValidateUrlAndMerchantID(t *testing.T) {
	testCases := []struct {
		testName           string
		url                string
		merchantID         string
		urlExpected        string
		merchantIDExpected int64
		err                error
	}{
		{
			testName:           "Correct_1",
			url:                " http://localhost:8090/api/load-file ",
			merchantID:         "1",
			urlExpected:        "http://localhost:8090/api/load-file",
			merchantIDExpected: 1,
			err:                nil,
		},
		{
			testName:           "Correct_2",
			url:                "localhost:8090/api/load-file",
			merchantID:         "1",
			urlExpected:        "",
			merchantIDExpected: -1,
			err:                constants.ErrIncorrect,
		},
		{
			testName:           "Correct_3",
			url:                "http://localhost:8090/api/load-file",
			merchantID:         "test",
			urlExpected:        "",
			merchantIDExpected: -1,
			err:                constants.ErrNotNumber,
		},
		{
			testName:           "Correct_3",
			url:                "http://localhost:8090/api/load-file",
			merchantID:         "-1",
			urlExpected:        "",
			merchantIDExpected: -1,
			err:                constants.ErrNotNumber,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			url, merchantID, err := ParseAndValidateUrlAndMerchantID(tc.url, tc.merchantID)
			require.Equal(t, tc.urlExpected, url)
			require.Equal(t, tc.merchantIDExpected, merchantID)
			require.Equal(t, tc.err, err)
		})
	}
}
