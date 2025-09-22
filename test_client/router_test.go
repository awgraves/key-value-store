package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockAPIv1Client struct {
	mock.Mock
}

func (m *mockAPIv1Client) SetKey(key string, value any) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *mockAPIv1Client) DeleteKey(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *mockAPIv1Client) GetKey(key string) (any, error) {
	args := m.Called(key)
	return args.Get(0), args.Error(1)
}

type routerTestSuite struct {
	suite.Suite
	mockClient *mockAPIv1Client
	router     *gin.Engine
}

func (s *routerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockClient = new(mockAPIv1Client)
	s.router = setupRouter(s.mockClient)
}

func (s *routerTestSuite) TestTestDeletion_Success() {
	// Mock the expected calls for successful deletion test
	s.mockClient.On("SetKey", "test-key", "test-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return("test-value", nil).Once()
	s.mockClient.On("DeleteKey", "test-key").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return(nil, nil).Once()

	req, _ := http.NewRequest("GET", "/api/v1/test_deletion", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Test deletion successful")

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestTestDeletion_SetKeyError() {
	// Mock SetKey to return an error
	s.mockClient.On("SetKey", "test-key", "test-value").Return(assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/test_deletion", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Error setting test key")
	assert.Contains(s.T(), resp.Body.String(), assert.AnError.Error())

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestTestDeletion_GetKeyAfterSetError() {
	// Mock successful SetKey but GetKey returns error
	s.mockClient.On("SetKey", "test-key", "test-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/test_deletion", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Error getting test key after setting")

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestTestDeletion_VerificationAfterSetFailed() {
	// Mock SetKey and GetKey but GetKey returns wrong value
	s.mockClient.On("SetKey", "test-key", "test-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return("wrong-value", nil)

	req, _ := http.NewRequest("GET", "/api/v1/test_deletion", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Error verifying test key after setting")

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestTestDeletion_DeleteKeyError() {
	// Mock successful SetKey and GetKey, but DeleteKey returns error
	s.mockClient.On("SetKey", "test-key", "test-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return("test-value", nil).Once()
	s.mockClient.On("DeleteKey", "test-key").Return(assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/test_deletion", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Error deleting test key")

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestTestDeletion_GetKeyAfterDeleteError() {
	// Mock successful operations until GetKey after delete
	s.mockClient.On("SetKey", "test-key", "test-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return("test-value", nil).Once()
	s.mockClient.On("DeleteKey", "test-key").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return(nil, assert.AnError).Once()

	req, _ := http.NewRequest("GET", "/api/v1/test_deletion", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Error getting test key after deletion")

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestTestDeletion_VerificationAfterDeleteFailed() {
	// Mock all operations successful but key still exists after deletion
	s.mockClient.On("SetKey", "test-key", "test-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return("test-value", nil).Once()
	s.mockClient.On("DeleteKey", "test-key").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return("still-exists", nil).Once()

	req, _ := http.NewRequest("GET", "/api/v1/test_deletion", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Error verifying test key after deletion")

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestTestOverwrite_Success() {
	// Mock the expected calls for successful overwrite test
	s.mockClient.On("SetKey", "test-key", "test-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return("test-value", nil).Once()
	s.mockClient.On("SetKey", "test-key", "new-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return("new-value", nil).Once()

	req, _ := http.NewRequest("GET", "/api/v1/test_overwrite", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Test overwrite successful")

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestTestOverwrite_SetKeyError() {
	// Mock SetKey to return an error
	s.mockClient.On("SetKey", "test-key", "test-value").Return(assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/test_overwrite", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Error setting test key")

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestTestOverwrite_GetKeyAfterSetError() {
	// Mock successful SetKey but GetKey returns error
	s.mockClient.On("SetKey", "test-key", "test-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/test_overwrite", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Error getting test key after setting")

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestTestOverwrite_VerificationAfterSetFailed() {
	// Mock SetKey and GetKey but GetKey returns wrong value
	s.mockClient.On("SetKey", "test-key", "test-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return("wrong-value", nil)

	req, _ := http.NewRequest("GET", "/api/v1/test_overwrite", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Error verifying test key after setting")

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestTestOverwrite_SetKeyAgainError() {
	// Mock successful first operations but second SetKey fails
	s.mockClient.On("SetKey", "test-key", "test-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return("test-value", nil).Once()
	s.mockClient.On("SetKey", "test-key", "new-value").Return(assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/test_overwrite", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Error setting test key again")

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestTestOverwrite_GetKeyAfterOverwriteError() {
	// Mock successful operations until final GetKey
	s.mockClient.On("SetKey", "test-key", "test-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return("test-value", nil).Once()
	s.mockClient.On("SetKey", "test-key", "new-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return(nil, assert.AnError).Once()

	req, _ := http.NewRequest("GET", "/api/v1/test_overwrite", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Error getting test key after overwriting")

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestTestOverwrite_VerificationAfterOverwriteFailed() {
	// Mock all operations successful but final verification fails
	s.mockClient.On("SetKey", "test-key", "test-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return("test-value", nil).Once()
	s.mockClient.On("SetKey", "test-key", "new-value").Return(nil)
	s.mockClient.On("GetKey", "test-key").Return("old-value", nil).Once()

	req, _ := http.NewRequest("GET", "/api/v1/test_overwrite", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "Error verifying test key after overwriting")

	s.mockClient.AssertExpectations(s.T())
}

func (s *routerTestSuite) TestConfig() {
	req, _ := http.NewRequest("GET", "/api/v1/config", nil)
	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Contains(s.T(), resp.Body.String(), "kv_api_v1_base_url")

	// Should contain the default URL since no environment variable is set
	assert.Contains(s.T(), resp.Body.String(), "http://localhost:8080/api/v1")
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(routerTestSuite))
}
