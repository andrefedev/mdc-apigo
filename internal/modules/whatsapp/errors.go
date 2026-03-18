package whatsapp

import (
	"fmt"
	"net/http"
)

// RequestError represents a typed error response from the Graph API.
type RequestError struct {
	StatusCode int
	Code       int
	Message    string
	Type       string
	Subcode    int
	TraceID    string
	Body       string
}

func (e *RequestError) Error() string {
	if e == nil {
		return "<nil>"
	}

	switch {
	case e.Code != 0 && e.Message != "":
		return fmt.Sprintf("whatsapp api status %d code %d: %s", e.StatusCode, e.Code, e.Message)
	case e.Message != "":
		return fmt.Sprintf("whatsapp api status %d: %s", e.StatusCode, e.Message)
	default:
		return fmt.Sprintf("whatsapp api status %d", e.StatusCode)
	}
}

func (e *RequestError) Temporary() bool {
	if e == nil {
		return false
	}

	if e.RateLimited() {
		return true
	}

	return e.StatusCode >= http.StatusInternalServerError
}

func (e *RequestError) Unauthorized() bool {
	if e == nil {
		return false
	}

	return e.StatusCode == http.StatusUnauthorized || e.StatusCode == http.StatusForbidden
}

func (e *RequestError) RateLimited() bool {
	if e == nil {
		return false
	}

	return e.StatusCode == http.StatusTooManyRequests
}

type graphErrorBody struct {
	Type         string `json:"type"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
	FBTraceID    string `json:"fbtrace_id"`
	ErrorSubcode int    `json:"error_subcode"`
}

type graphErrorEnvelope struct {
	Error graphErrorBody `json:"error"`
}
