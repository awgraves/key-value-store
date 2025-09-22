package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockStore struct {
	mock.Mock
}

func (m *mockStore) Get(key string) any {
	args := m.Called(key)
	return args.Get(0)
}

func (m *mockStore) Set(key string, value any) {
	m.Called(key, value)
}

func (m *mockStore) Delete(key string) {
	m.Called(key)
}

type routerTestSuite struct {
	suite.Suite
	mockStore *mockStore
	router    *gin.Engine
}

func (s *routerTestSuite) SetupTest() {
	s.mockStore = new(mockStore)
	s.router = setupRouter(s.mockStore)
}
func (s *routerTestSuite) TestGetKey() {
	call := s.mockStore.On("Get", "foo").Return("bar")

	req, _ := http.NewRequest("GET", "/api/v1/keys/foo", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), `{"value":"bar"}`, resp.Body.String())

	call.Unset()
}

func (s *routerTestSuite) TestSetKey() {
	call := s.mockStore.On("Set", "foo", "bar").Return()
	req, _ := http.NewRequest("POST", "/api/v1/keys/foo", strings.NewReader(`{"value":"bar"}`))
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), `{"message":"Key set."}`, resp.Body.String())

	s.mockStore.AssertCalled(s.T(), "Set", "foo", "bar")

	call.Unset()
}

func (s *routerTestSuite) TestDeleteKey() {
	call := s.mockStore.On("Delete", "foo").Return()
	req, _ := http.NewRequest("DELETE", "/api/v1/keys/foo", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), `{"message":"Key deleted."}`, resp.Body.String())

	s.mockStore.AssertCalled(s.T(), "Delete", "foo")

	call.Unset()
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(routerTestSuite))
}
