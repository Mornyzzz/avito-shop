package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slog"

	"avito-shop/internal/controller/worker"
	worker_mocks "avito-shop/internal/controller/worker/mocks"
	"avito-shop/internal/entity"
	auth_mocks "avito-shop/internal/usecase/auth/mocks"
)

func TestAuthRoute_Auth(t *testing.T) {
	mockAuthUC := new(auth_mocks.Auth)
	mockWorkerPool := new(worker_mocks.PoolI)
	log := slog.Default()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := `{"username": "testuser", "password": "testpass"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	mockWorkerPool.On("Submit", mock.AnythingOfType("worker.Task")).Run(func(args mock.Arguments) {
		task := args.Get(0).(worker.Task)
		task()
	}).Return()
	mockAuthUC.On("Login", mock.Anything, entity.User{Username: "testuser", Password: "testpass"}).Return("testtoken", nil)

	authRoute := &AuthRoute{authUC: mockAuthUC, wp: mockWorkerPool, log: log}
	authRoute.Auth(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"token": "testtoken"}`, w.Body.String())

	mockAuthUC.AssertExpectations(t)
	mockWorkerPool.AssertExpectations(t)
}

func TestAuthRoute_Auth_InvalidCredentials(t *testing.T) {
	mockAuthUC := new(auth_mocks.Auth)
	mockWorkerPool := new(worker_mocks.PoolI)
	log := slog.Default()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	reqBody := `{"username": "testuser", "password": "wrongpass"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/auth", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	mockWorkerPool.On("Submit", mock.AnythingOfType("worker.Task")).Run(func(args mock.Arguments) {
		task := args.Get(0).(worker.Task)
		task()
	}).Return()
	mockAuthUC.On("Login", mock.Anything, entity.User{Username: "testuser", Password: "wrongpass"}).Return("", errors.New("invalid credentials"))

	authRoute := &AuthRoute{authUC: mockAuthUC, wp: mockWorkerPool, log: log}
	authRoute.Auth(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"error": "Invalid credentials"}`, w.Body.String())

	mockAuthUC.AssertExpectations(t)
	mockWorkerPool.AssertExpectations(t)
}
