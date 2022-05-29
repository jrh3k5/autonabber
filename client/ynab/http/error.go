package http

import "fmt"

type ErrorContainer struct {
	APIError *APIError `json:"error"`
}

type APIError struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

func (ec *ErrorContainer) Error() string {
	if ec.APIError == nil {
		return "<no error details available from API response>"
	}

	return fmt.Sprintf("error ID: %s; error name: %s; error detail: %s", ec.APIError.ID, ec.APIError.Name, ec.APIError.Detail)
}
