package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		return fmt.Errorf("failed to invoke GET request for '%s': %w", requestPath, err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		errorContainer, err := handleError(httpResponse)
		if err != nil {
			return fmt.Errorf("failed to extract error details for GET request to '%s': %w", requestPath, err)
		} else if errorContainer != nil {
			return fmt.Errorf("failed to issue GET request to '%s': %w", requestPath, errorContainer)
		}

		return fmt.Errorf("unexpected status in response to GET request to '%s': %d", requestPath, httpResponse.StatusCode)
	}

	if err := json.NewDecoder(httpResponse.Body).Decode(result); err != nil {
		return fmt.Errorf("error decoding response to GET request to '%s': %w", requestPath, err)
	}

	return nil
}

// Patch issues a PATCH request to the given request path with the given request body
func (c *Client) Patch(requestPath string, requestBody io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("PATCH", baseURL+"/"+requestPath, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PATCH HTTP request to path '%s': %w", requestPath, err)
	}

	request.Header.Set("Authorization", "Bearer "+c.accessToken)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(request)
}

// PatchJSON submits a PATCH request formatted as JSON
func (c *Client) PatchJSON(requestPath string, requestBody interface{}) error {
	jsonBytesBuffer := &bytes.Buffer{}
	if err := json.NewEncoder(jsonBytesBuffer).Encode(requestBody); err != nil {
		return fmt.Errorf("failed to encode request body to JSON: %w", err)
	}

	httpResponse, err := c.Patch(requestPath, jsonBytesBuffer)
	if err != nil {
		return fmt.Errorf("failed to invoke PATCH request for '%s': %w", requestPath, err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != 200 {
		errorContainer, err := handleError(httpResponse)
		if err != nil {
			return fmt.Errorf("failed to extract error details for PATCH request to '%s': %w", requestPath, err)
		} else if errorContainer != nil {
			return fmt.Errorf("failed to issue PATCH request to '%s': %w", requestPath, errorContainer)
		}

		return fmt.Errorf("unexpected status in response to PATCH request to '%s': %d", requestPath, httpResponse.StatusCode)
	}

	return nil
}

// handleError is used for handling a non-OK HTTP response and attempting to get the error details out of it
// The returned *ErrorContainer can be nil if the given response is nil if the response has no body
func handleError(httpResponse *http.Response) (*ErrorContainer, error) {
	var errorContainer *ErrorContainer
	if err := json.NewDecoder(httpResponse.Body).Decode(&errorContainer); err != nil {
		return nil, err
	}
	return errorContainer, nil
}
