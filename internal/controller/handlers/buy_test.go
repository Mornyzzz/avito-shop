package handlers

import (
	"avito-shop/internal/controller/worker"
	worker_mocks "avito-shop/internal/controller/worker/mocks"
	buy_mocks "avito-shop/internal/usecase/buy/mocks"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuyRoute_Buy_Success(t *testing.T) {
	mockBuyUC := new(buy_mocks.Buy)
	mockWorkerPool := new(worker_mocks.PoolI)
	log := slog.Default()

	mockWorkerPool.On("Submit", mock.AnythingOfType("worker.Task")).Run(func(args mock.Arguments) {
		task := args.Get(0).(worker.Task)
		task()
	}).Return()

	mockBuyUC.On("BuyItem", mock.Anything, "testuser", "testitem").Return(nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/buy/testitem", nil)
	c.Params = gin.Params{gin.Param{Key: "item", Value: "testitem"}}
	c.Set("username", "testuser")

	buyRoute := &BuyRoute{
		buyUC: mockBuyUC,
		wp:    mockWorkerPool,
		log:   log,
	}

	buyRoute.Buy(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `"Item purchased successfully"`, w.Body.String())

	mockBuyUC.AssertExpectations(t)
	mockWorkerPool.AssertExpectations(t)
}

func TestBuyRoute_Buy_InvalidRequest(t *testing.T) {
	mockBuyUC := new(buy_mocks.Buy)
	mockWorkerPool := new(worker_mocks.PoolI)
	log := slog.Default()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/buy/", nil)
	c.Set("username", "testuser") // Устанавливаем username в контексте

	buyRoute := &BuyRoute{
		buyUC: mockBuyUC,
		wp:    mockWorkerPool,
		log:   log,
	}

	buyRoute.Buy(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"Invalid request"}`, w.Body.String())
}

func TestBuyRoute_Buy_Unauthorized(t *testing.T) {
	mockBuyUC := new(buy_mocks.Buy)
	mockWorkerPool := new(worker_mocks.PoolI)
	log := slog.Default()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/buy/testitem", nil)
	c.Params = gin.Params{gin.Param{Key: "item", Value: "testitem"}}

	buyRoute := &BuyRoute{
		buyUC: mockBuyUC,
		wp:    mockWorkerPool,
		log:   log,
	}

	buyRoute.Buy(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"error":"User not authenticated"}`, w.Body.String())
}

func TestBuyRoute_Buy_InternalError(t *testing.T) {
	mockBuyUC := new(buy_mocks.Buy)
	mockWorkerPool := new(worker_mocks.PoolI)
	log := slog.Default()

	mockWorkerPool.On("Submit", mock.AnythingOfType("worker.Task")).Run(func(args mock.Arguments) {
		task := args.Get(0).(worker.Task)
		task()
	}).Return()

	mockBuyUC.On("BuyItem", mock.Anything, "testuser", "testitem").Return(errors.New("internal error"))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/buy/testitem", nil)
	c.Params = gin.Params{gin.Param{Key: "item", Value: "testitem"}}
	c.Set("username", "testuser")

	buyRoute := &BuyRoute{
		buyUC: mockBuyUC,
		wp:    mockWorkerPool,
		log:   log,
	}

	buyRoute.Buy(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":"Failed to buy item"}`, w.Body.String())

	mockBuyUC.AssertExpectations(t)
	mockWorkerPool.AssertExpectations(t)
}
