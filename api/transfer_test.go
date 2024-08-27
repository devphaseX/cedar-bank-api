package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/devphasex/cedar-bank-api/db/mock"
	db "github.com/devphasex/cedar-bank-api/db/sqlc"
	"github.com/devphasex/cedar-bank-api/token"
	"github.com/devphasex/cedar-bank-api/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	user3, _ := randomUser(t)

	account1 := randomAccount(user1.ID)
	account2 := randomAccount(user2.ID)
	account3 := randomAccount(user3.ID)

	account1.Currency = string(util.USD)
	account2.Currency = string(util.USD)
	account3.Currency = string(util.CAD)

	testCases := []struct {
		name          string
		FromAccountID int64
		ToAccountID   int64
		Amount        float64
		Currency      string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "TransferFromAccountOneToTwo",
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        20,
			Currency:      string(util.USD),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.ID, user1.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByID(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccountByID(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:          "InsufficientFundTransfer",
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        account1.Balance + 1,
			Currency:      string(util.USD),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.ID, user1.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByID(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccountByID(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).Return(nil, db.ErrFundNotSufficient)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name:          "InvalidTransferDifferentCurrency",
			FromAccountID: account1.ID,
			ToAccountID:   account3.ID,
			Amount:        20,
			Currency:      string(util.USD),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.ID, user1.Email, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByID(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccountByID(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			b, err := json.Marshal(TransferRequest{
				FromAccountID: tc.FromAccountID,
				ToAccountID:   tc.ToAccountID,
				Amount:        tc.Amount,
				Currency:      tc.Currency,
			})

			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, "/transfer", bytes.NewBuffer(b))

			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}

}
