package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://api.youneedabudget.com/v1"

type Client struct {
	httpClient  *http.Client
	accessToken string
}

func NewClient(accessToken string) (*Client, error) {
	return &Client{
		accessToken: accessToken,
		httpClient:  &http.Client{},
	}, nil
}

func (c *Client) Get(requestPath string) (*http.Response, error) {
	request, err := http.NewRequest("GET", baseURL+"/"+requestPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate GET HTTP request to path '%s': %w", requestPath, err)
	}

	request.Header.Set("Authorization", "Bearer "+c.accessToken)
	request.Header.Set("Accept", "application/json")
	return c.httpClient.Do(request)
}

// handleError is used for handling a non-OK HTTP response and attempting to get the error details out of it
func handleError(httpResponse *http.Response) (*ErrorContainer, error) {
	var errorContainer *ErrorContainer
	if err := json.NewDecoder(httpResponse.Body).Decode(&errorContainer); err != nil {
		return nil, err
	}
	return errorContainer, nil
}
