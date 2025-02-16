package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slog"

	"avito-shop/internal/controller/worker"
	workermocks "avito-shop/internal/controller/worker/mocks"
	"avito-shop/internal/entity"
	infomocks "avito-shop/internal/usecase/info/mocks"
)

func TestInfoRoute_Info_Success(t *testing.T) {
	mockInfoUC := new(infomocks.Info)
	mockWorkerPool := new(workermocks.PoolI)
	log := slog.Default()

	mockWorkerPool.On("Submit", mock.AnythingOfType("worker.Task")).Run(func(args mock.Arguments) {
		task := args.Get(0).(worker.Task)
		task()
	}).Return()

	expectedInfo := &entity.Info{
		Coins: 100,
		Inventory: []entity.InventoryItem{
			{Name: "Item1", Quantity: 1},
			{Name: "Item2", Quantity: 2},
		},
		CoinHistory: entity.CoinHistory{
			Received: []entity.ReceivedTransaction{
				{FromUser: "user1", Amount: 50},
			},
			Sent: []entity.SentTransaction{
				{ToUser: "user2", Amount: 30},
			},
		},
	}
	mockInfoUC.On("GetInfo", mock.Anything, "testuser").Return(expectedInfo, nil)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/info", http.NoBody)
	c.Set("username", "testuser")

	infoRoute := &InfoRoute{
		infoUC: mockInfoUC,
		wp:     mockWorkerPool,
		log:    log,
	}

	infoRoute.Info(c)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{
		"coins": 100,
		"inventory": [
			{"name": "Item1", "quantity": 1},
			{"name": "Item2", "quantity": 2}
		],
		"coinHistory": {
			"received": [
				{"fromUser": "user1", "amount": 50}
			],
			"sent": [
				{"toUser": "user2", "amount": 30}
			]
		}
	}`, w.Body.String())

	mockInfoUC.AssertExpectations(t)
	mockWorkerPool.AssertExpectations(t)
}

func TestInfoRoute_Info_Unauthorized(t *testing.T) {
	mockInfoUC := new(infomocks.Info)
	mockWorkerPool := new(workermocks.PoolI)
	log := slog.Default()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/info", http.NoBody)

	infoRoute := &InfoRoute{
		infoUC: mockInfoUC,
		wp:     mockWorkerPool,
		log:    log,
	}

	infoRoute.Info(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"error":"User not authenticated"}`, w.Body.String())
}

func TestInfoRoute_Info_InternalError(t *testing.T) {
	mockInfoUC := new(infomocks.Info)
	mockWorkerPool := new(workermocks.PoolI)
	log := slog.Default()

	mockWorkerPool.On("Submit", mock.AnythingOfType("worker.Task")).Run(func(args mock.Arguments) {
		task := args.Get(0).(worker.Task)
		task()
	}).Return()

	mockInfoUC.On("GetInfo", mock.Anything, "testuser").Return(nil, errors.New("internal error"))

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/info", http.NoBody)
	c.Set("username", "testuser")

	infoRoute := &InfoRoute{
		infoUC: mockInfoUC,
		wp:     mockWorkerPool,
		log:    log,
	}

	infoRoute.Info(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":"Failed to get info"}`, w.Body.String())

	mockInfoUC.AssertExpectations(t)
	mockWorkerPool.AssertExpectations(t)
}
