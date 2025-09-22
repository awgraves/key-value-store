package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type APIv1Client interface {
	SetKey(key string, value any) error
	DeleteKey(key string) error
	GetKey(key string) (any, error) // returns unwrapped value from response
}

type apiV1Client struct {
	baseURL string
}

func NewAPIv1Client(baseURL string) *apiV1Client {
	return &apiV1Client{baseURL: baseURL}
}

func (c *apiV1Client) SetKey(key string, value any) error {
	fmt.Println("setting key: ", key, "value: ", value)
	bodyMap := map[string]any{
		"value": value,
	}
	body, err := json.Marshal(bodyMap)
	if err != nil {
		return err
	}
	response, err := http.Post(fmt.Sprintf("%s/keys/%s", c.baseURL, key), "application/json", strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		fmt.Println("response body: ", string(bodyBytes))
		return fmt.Errorf("failed to set key: %s, response: %s", response.Status, string(bodyBytes))
	}
	fmt.Println("response: ", response)
	return nil
}

func (c *apiV1Client) DeleteKey(key string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/keys/%s", c.baseURL, key), nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		return fmt.Errorf("failed to delete key: %s, response: %s", response.Status, string(bodyBytes))
	}
	return nil
}

func (c *apiV1Client) GetKey(key string) (any, error) {
	response, err := http.Get(fmt.Sprintf("%s/keys/%s", c.baseURL, key))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("failed to get key: %s, response: %s", response.Status, string(bodyBytes))
	}
	var valueResponse struct {
		Value any `json:"value"`
	}
	err = json.NewDecoder(response.Body).Decode(&valueResponse)
	if err != nil {
		return nil, err
	}
	return valueResponse.Value, nil
}
