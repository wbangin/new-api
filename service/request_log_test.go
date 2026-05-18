package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
)

func TestWriteAndSearchRequestLog(t *testing.T) {
	// Setup: use a temporary directory for log files
	tmpDir := t.TempDir()
	requestLogDir = filepath.Join(tmpDir, "request_logs")
	_ = os.MkdirAll(requestLogDir, 0755)

	// Reset file state
	requestLogFile = nil
	requestLogDate = ""

	entry := &RequestLogEntry{
		RequestId:    "test-req-12345",
		Timestamp:    time.Now().Unix(),
		TokenId:      1,
		TokenName:    "test-token",
		UserId:       100,
		Model:        "gpt-4",
		RelayMode:    1,
		RequestBody:  `{"messages":[{"role":"user","content":"hello"}]}`,
		ResponseBody: `{"choices":[{"message":{"content":"hi"}}]}`,
		StatusCode:   200,
		IsStream:     false,
	}

	// Write the log entry
	WriteRequestLog(entry)

	// Close and flush
	if requestLogFile != nil {
		_ = requestLogFile.Sync()
	}

	// Search for the entry
	found, err := SearchRequestLog("test-req-12345")
	if err != nil {
		t.Fatalf("SearchRequestLog failed: %v", err)
	}
	if found == nil {
		t.Fatal("SearchRequestLog returned nil")
	}
	if found.RequestId != "test-req-12345" {
		t.Errorf("expected request_id 'test-req-12345', got '%s'", found.RequestId)
	}
	if found.Model != "gpt-4" {
		t.Errorf("expected model 'gpt-4', got '%s'", found.Model)
	}
	if found.TokenName != "test-token" {
		t.Errorf("expected token_name 'test-token', got '%s'", found.TokenName)
	}
	if found.UserId != 100 {
		t.Errorf("expected user_id 100, got %d", found.UserId)
	}

	// Search for a non-existent entry
	_, err = SearchRequestLog("non-existent-id")
	if err == nil {
		t.Error("expected error for non-existent request_id, got nil")
	}
}

func TestWriteRequestLogNil(t *testing.T) {
	// Should not panic
	WriteRequestLog(nil)
}

func TestSearchRequestLogEmpty(t *testing.T) {
	_, err := SearchRequestLog("")
	if err == nil {
		t.Error("expected error for empty request_id, got nil")
	}
}

// Ensure common package is imported (used by WriteRequestLog for Marshal)
var _ = common.Version
