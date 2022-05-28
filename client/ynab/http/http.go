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

// Get issues an HTTP GET request to the given request path
func (c *Client) Get(requestPath string) (*http.Response, error) {
	request, err := http.NewRequest("GET", baseURL+"/"+requestPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate GET HTTP request to path '%s': %w", requestPath, err)
	}

	request.Header.Set("Authorization", "Bearer "+c.accessToken)
	request.Header.Set("Accept", "application/json")
	return c.httpClient.Do(request)
}

// GetJSON issues a GET request for JSON data to the given request path and populates the given result with the body of the response
func (c *Client) GetJSON(requestPath string, result interface{}) error {
	httpResponse, err := c.Get(requestPath)
	if err != nil {
		return fmt.Errorf("failed to invoke request for '%s': %w", requestPath, err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		errorContainer, err := handleError(httpResponse)
		if err != nil {
			return fmt.Errorf("failed to extract error details for request to '%s': %w", requestPath, err)
		}
		return errorContainer
	}

	if err := json.NewDecoder(httpResponse.Body).Decode(result); err != nil {
		return fmt.Errorf("error decoding response from '%s': %w", requestPath, err)
	}

	return nil
}

// handleError is used for handling a non-OK HTTP response and attempting to get the error details out of it
func handleError(httpResponse *http.Response) (*ErrorContainer, error) {
	var errorContainer *ErrorContainer
	if err := json.NewDecoder(httpResponse.Body).Decode(&errorContainer); err != nil {
		return nil, err
	}
	return errorContainer, nil
}
