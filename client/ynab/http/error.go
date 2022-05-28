package http

import "fmt"

type ErrorContainer struct {
	apiError *APIError `json:"error"`
}

type APIError struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

func (ec *ErrorContainer) Error() string {
	if ec.Error == nil {
		return "<no error details available from API response>"
	}

	return fmt.Sprint("error ID: %s; error name: %s; error detail: %s", ec.apiError.ID, ec.apiError.Name, ec.apiError.Detail)
}
