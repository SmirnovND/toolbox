package http

import (
	"testing"
	"time"
)

// Тест для NewAPIClient
func TestNewAPIClient(t *testing.T) {
	client := NewAPIClient()
	
	if client == nil {
		t.Error("NewAPIClient() returned nil")
	}
	
	if client.HTTPClient == nil {
		t.Error("NewAPIClient().HTTPClient is nil")
	}
	
	if client.HTTPClient.Timeout != 10*time.Second {
		t.Errorf("NewAPIClient().HTTPClient.Timeout = %v, want %v", client.HTTPClient.Timeout, 10*time.Second)
	}
	
	if client.BaseURL != "" {
		t.Errorf("NewAPIClient().BaseURL = %v, want empty string", client.BaseURL)
	}
}

// Тест для SetBaseURL
func TestAPIClient_SetBaseURL(t *testing.T) {
	client := NewAPIClient()
	
	testURL := "https://api.example.com"
	client.SetBaseURL(testURL)
	
	if client.BaseURL != testURL {
		t.Errorf("After SetBaseURL(%q), client.BaseURL = %q, want %q", testURL, client.BaseURL, testURL)
	}
}