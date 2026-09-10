package openai

import (
	"testing"

	"github.com/yaoapp/gou/http"
)

func TestIsError_TransportError(t *testing.T) {
	o := OpenAI{}
	res := &http.Response{
		Status:  0,
		Code:    0,
		Data:    nil,
		Message: "Post \"https://api.openai.com/v1/audio/transcriptions\": context deadline exceeded",
	}
	err := o.isError(res)
	if err == nil {
		t.Fatal("expected error for Status=0")
	}
	if err.Message != res.Message {
		t.Errorf("expected message %q, got %q", res.Message, err.Message)
	}
	t.Logf("isError returned: code=%d message=%q", err.Code, err.Message)
}

func TestIsError_APIError(t *testing.T) {
	o := OpenAI{}
	res := &http.Response{
		Status: 400,
		Code:   400,
		Data: map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Invalid file format.",
				"code":    "invalid_request",
			},
		},
		Message: "transport-level fallback",
	}
	err := o.isError(res)
	if err == nil {
		t.Fatal("expected error for Status=400")
	}
	// API-level error.message should take priority over res.Message.
	expected := "OpenAI invalid_request Invalid file format."
	if err.Message != expected {
		t.Errorf("expected %q, got %q", expected, err.Message)
	}
}

func TestIsError_StringData(t *testing.T) {
	o := OpenAI{}
	res := &http.Response{
		Status:  500,
		Code:    500,
		Data:    "Internal Server Error",
		Message: "",
	}
	err := o.isError(res)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Message != "Internal Server Error" {
		t.Errorf("expected %q, got %q", "Internal Server Error", err.Message)
	}
}

func TestIsError_Success(t *testing.T) {
	o := OpenAI{}
	res := &http.Response{
		Status: 200,
		Code:   200,
		Data:   map[string]interface{}{"text": "hello"},
	}
	err := o.isError(res)
	if err != nil {
		t.Errorf("expected nil error for Status=200, got %v", err)
	}
}

func TestExtractHost(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"https://api.openai.com", "api.openai.com"},
		{"https://api.openai.com:443/v1", "api.openai.com"},
		{"http://localhost:8080", "localhost"},
		{"not-a-url", "not-a-url"},
	}
	for _, c := range cases {
		got := extractHost(c.input)
		if got != c.want {
			t.Errorf("extractHost(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestCheckAndRetry_TransportErrorRetries(t *testing.T) {
	o := OpenAI{host: "https://api.openai.com"}

	callCount := 0
	retryFn := func() *http.Response {
		callCount++
		return &http.Response{
			Status: 200,
			Code:   200,
			Data:   map[string]interface{}{"text": "ok"},
		}
	}

	// First response is a transport error.
	firstRes := &http.Response{
		Status:  0,
		Code:    0,
		Data:    nil,
		Message: "connection refused",
	}

	data, err := o.checkAndRetry(firstRes, retryFn)
	if err != nil {
		t.Fatalf("expected retry to succeed, got error: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected retryFn called once, got %d", callCount)
	}
	if data == nil {
		t.Error("expected non-nil data from retry")
	}
}

func TestCheckAndRetry_APIErrorNoRetry(t *testing.T) {
	o := OpenAI{host: "https://api.openai.com"}

	callCount := 0
	retryFn := func() *http.Response {
		callCount++
		return &http.Response{Status: 200, Data: "ok"}
	}

	// API error (non-zero status) should NOT trigger retry.
	firstRes := &http.Response{
		Status: 401,
		Code:   401,
		Data:   "Unauthorized",
	}

	_, err := o.checkAndRetry(firstRes, retryFn)
	if err == nil {
		t.Fatal("expected error for 401")
	}
	if callCount != 0 {
		t.Errorf("retryFn should not be called for API errors, got %d", callCount)
	}
}

func TestCheckAndRetry_Success(t *testing.T) {
	o := OpenAI{host: "https://api.openai.com"}

	retryFn := func() *http.Response {
		t.Error("retryFn should not be called on success")
		return nil
	}

	firstRes := &http.Response{
		Status: 200,
		Code:   200,
		Data:   map[string]interface{}{"result": "ok"},
	}

	data, err := o.checkAndRetry(firstRes, retryFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data == nil {
		t.Error("expected non-nil data")
	}
}
