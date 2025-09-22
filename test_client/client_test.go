package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type clientTestSuite struct {
	suite.Suite
}

func (s *clientTestSuite) TestNewAPIv1Client() {
	baseURL := "http://localhost:8080"
	client := NewAPIv1Client(baseURL)

	assert.NotNil(s.T(), client)
	assert.Equal(s.T(), baseURL, client.baseURL)
	assert.Implements(s.T(), (*APIv1Client)(nil), client)
}

func (s *clientTestSuite) TestSetKey_Success() {
	// Create a test server that responds with 200 OK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "POST", r.Method)
		assert.Equal(s.T(), "/keys/testkey", r.URL.Path)
		assert.Equal(s.T(), "application/json", r.Header.Get("Content-Type"))

		// Verify the request body
		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "testvalue", body["value"])

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewAPIv1Client(server.URL)
	err := client.SetKey("testkey", "testvalue")

	assert.NoError(s.T(), err)
}

func (s *clientTestSuite) TestSetKey_WithComplexValue() {
	// Test with a complex value (map)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(s.T(), err)

		expectedValue := map[string]interface{}{
			"nested": "data",
			"number": float64(42), // JSON numbers are decoded as float64
		}
		assert.Equal(s.T(), expectedValue, body["value"])

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewAPIv1Client(server.URL)
	complexValue := map[string]interface{}{
		"nested": "data",
		"number": 42,
	}
	err := client.SetKey("complexkey", complexValue)

	assert.NoError(s.T(), err)
}

func (s *clientTestSuite) TestSetKey_ServerError() {
	// Create a test server that responds with 500 Internal Server Error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewAPIv1Client(server.URL)
	err := client.SetKey("testkey", "testvalue")

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to set key")
	assert.Contains(s.T(), err.Error(), "500")
}

func (s *clientTestSuite) TestSetKey_InvalidJSON() {
	client := NewAPIv1Client("http://localhost:8080")

	// Test with a value that cannot be marshaled to JSON
	invalidValue := make(chan int)
	err := client.SetKey("testkey", invalidValue)

	assert.Error(s.T(), err)
}

func (s *clientTestSuite) TestSetKey_NetworkError() {
	// Use an invalid URL to simulate network error
	client := NewAPIv1Client("http://invalid-url-that-does-not-exist:9999")
	err := client.SetKey("testkey", "testvalue")

	assert.Error(s.T(), err)
}

func (s *clientTestSuite) TestGetKey_Success() {
	expectedValue := "testvalue"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "GET", r.Method)
		assert.Equal(s.T(), "/keys/testkey", r.URL.Path)

		response := map[string]interface{}{
			"value": expectedValue,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAPIv1Client(server.URL)
	value, err := client.GetKey("testkey")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedValue, value)
}

func (s *clientTestSuite) TestGetKey_ComplexValue() {
	expectedValue := map[string]interface{}{
		"nested": "data",
		"number": float64(42),
		"array":  []interface{}{"item1", "item2"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"value": expectedValue,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAPIv1Client(server.URL)
	value, err := client.GetKey("complexkey")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedValue, value)
}

func (s *clientTestSuite) TestGetKey_ServerError() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewAPIv1Client(server.URL)
	value, err := client.GetKey("testkey")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), value)
	assert.Contains(s.T(), err.Error(), "failed to get key")
	assert.Contains(s.T(), err.Error(), "500")
}

func (s *clientTestSuite) TestGetKey_InvalidJSON() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Return invalid JSON
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewAPIv1Client(server.URL)
	value, err := client.GetKey("testkey")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), value)
}

func (s *clientTestSuite) TestGetKey_NetworkError() {
	client := NewAPIv1Client("http://invalid-url-that-does-not-exist:9999")
	value, err := client.GetKey("testkey")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), value)
}

func (s *clientTestSuite) TestDeleteKey_Success() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(s.T(), "DELETE", r.Method)
		assert.Equal(s.T(), "/keys/testkey", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewAPIv1Client(server.URL)
	err := client.DeleteKey("testkey")

	assert.NoError(s.T(), err)
}

func (s *clientTestSuite) TestDeleteKey_ServerError() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewAPIv1Client(server.URL)
	err := client.DeleteKey("testkey")

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to delete key")
	assert.Contains(s.T(), err.Error(), "500")
}

func (s *clientTestSuite) TestDeleteKey_NetworkError() {
	client := NewAPIv1Client("http://invalid-url-that-does-not-exist:9999")
	err := client.DeleteKey("testkey")

	assert.Error(s.T(), err)
}

func (s *clientTestSuite) TestDeleteKey_RequestCreationError() {
	// Test with an invalid URL that would cause NewRequest to fail
	client := NewAPIv1Client("ht tp://invalid-url-with-space")
	err := client.DeleteKey("testkey")

	assert.Error(s.T(), err)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(clientTestSuite))
}
