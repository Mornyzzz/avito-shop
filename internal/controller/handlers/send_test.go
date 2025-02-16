package handlers

import (
	"avito-shop/internal/controller/worker"
	workermocks "avito-shop/internal/controller/worker/mocks"
	sendmocks "avito-shop/internal/usecase/send/mocks"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendRoute_Send_Success(t *testing.T) {
	mockSendUC := new(sendmocks.Send)
	mockWorkerPool := new(workermocks.PoolI)
	log := slog.Default()

	mockWorkerPool.On("Submit", mock.AnythingOfType("worker.Task")).Run(func(args mock.Arguments) {
		task := args.Get(0).(worker.Task)
		task()
	}).Return()

	mockSendUC.On("SendCoin", mock.Anything, "senderUser", "receiverUser", 100).Return(nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := `{"toUser": "receiverUser", "amount": 100}`
	c.Request = httptest.NewRequest(http.MethodPost, "/sendCoin", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("username", "senderUser")

	sendRoute := &SendRoute{
		sendUC: mockSendUC,
		wp:     mockWorkerPool,
		log:    log,
	}

	sendRoute.Send(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `"Coins sent successfully"`, w.Body.String())

	mockSendUC.AssertExpectations(t)
	mockWorkerPool.AssertExpectations(t)
}

func TestSendRoute_Send_Unauthorized(t *testing.T) {
	mockSendUC := new(sendmocks.Send)
	mockWorkerPool := new(workermocks.PoolI)
	log := slog.Default()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := `{"toUser": "receiverUser", "amount": 100}`
	c.Request = httptest.NewRequest(http.MethodPost, "/sendCoin", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	sendRoute := &SendRoute{
		sendUC: mockSendUC,
		wp:     mockWorkerPool,
		log:    log,
	}

	sendRoute.Send(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"error":"User not authenticated"}`, w.Body.String())
}

func TestSendRoute_Send_InvalidRequest(t *testing.T) {
	mockSendUC := new(sendmocks.Send)
	mockWorkerPool := new(workermocks.PoolI)
	log := slog.Default()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := `{"toUser": "", "amount": 0}`
	c.Request = httptest.NewRequest(http.MethodPost, "/sendCoin", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("username", "senderUser")

	sendRoute := &SendRoute{
		sendUC: mockSendUC,
		wp:     mockWorkerPool,
		log:    log,
	}

	sendRoute.Send(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"Invalid request"}`, w.Body.String())
}

func TestSendRoute_Send_InternalError(t *testing.T) {
	mockSendUC := new(sendmocks.Send)
	mockWorkerPool := new(workermocks.PoolI)
	log := slog.Default()

	mockWorkerPool.On("Submit", mock.AnythingOfType("worker.Task")).Run(func(args mock.Arguments) {
		task := args.Get(0).(worker.Task)
		task()
	}).Return()

	mockSendUC.On("SendCoin", mock.Anything, "senderUser", "receiverUser", 100).Return(errors.New("internal error"))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := `{"toUser": "receiverUser", "amount": 100}`
	c.Request = httptest.NewRequest(http.MethodPost, "/sendCoin", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("username", "senderUser")

	sendRoute := &SendRoute{
		sendUC: mockSendUC,
		wp:     mockWorkerPool,
		log:    log,
	}

	sendRoute.Send(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":"Failed to send coins"}`, w.Body.String())

	mockSendUC.AssertExpectations(t)
	mockWorkerPool.AssertExpectations(t)
}
