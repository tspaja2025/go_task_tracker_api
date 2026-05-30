package tasks

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTaskRouter_InvalidIDFormat(t *testing.T) {
	// Create nil-initialized repository
	handler := NewHandler(nil)

	// Mock request to a garbage URL id path: /task/abc instead of /task/123
	req, err := http.NewRequest(http.MethodGet, "/task/abc", nil)
	if err != nil {
		t.Fatalf("Could not create HTTP request framework: %v", err)
	}

	// Reponse recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.Router(rr, req)

	// Assertion on expectations
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, rr.Code)
	}

	expectedCleanBody := `{"error": "Invalid task ID formatting"}`
	actualCleanBody := strings.TrimSpace(rr.Body.String())

	if actualCleanBody != expectedCleanBody {
		t.Errorf("Expected body payload %q, but instead got %q", expectedCleanBody, actualCleanBody)
	}
}
