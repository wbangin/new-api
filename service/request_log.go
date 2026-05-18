package service

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
)

// RequestLogEntry represents a single request/response log entry in JSONL format
type RequestLogEntry struct {
	RequestId    string `json:"request_id"`
	Timestamp    int64  `json:"timestamp"`
	TokenId      int    `json:"token_id"`
	TokenName    string `json:"token_name"`
	UserId       int    `json:"user_id"`
	Model        string `json:"model"`
	RelayMode    int    `json:"relay_mode"`
	RequestBody  string `json:"request_body"`
	ResponseBody string `json:"response_body"`
	StatusCode   int    `json:"status_code"`
	IsStream     bool   `json:"is_stream"`
}

var (
	requestLogMu      sync.Mutex
	requestLogFile    *os.File
	requestLogDate    string
	requestLogDir     string
	requestLogDirOnce sync.Once
)

// getRequestLogDir returns the directory for request log files
func getRequestLogDir() string {
	requestLogDirOnce.Do(func() {
		if *common.LogDir != "" {
			requestLogDir = filepath.Join(*common.LogDir, "request_logs")
		} else {
			requestLogDir = filepath.Join(".", "logs", "request_logs")
		}
		if _, err := os.Stat(requestLogDir); os.IsNotExist(err) {
			_ = os.MkdirAll(requestLogDir, 0755)
		}
	})
	return requestLogDir
}

// getRequestLogFile returns the current log file, rotating by date
func getRequestLogFile() (*os.File, error) {
	today := time.Now().Format("2006-01-02")

	requestLogMu.Lock()
	defer requestLogMu.Unlock()

	if requestLogFile != nil && requestLogDate == today {
		return requestLogFile, nil
	}

	// Close old file if date changed
	if requestLogFile != nil {
		_ = requestLogFile.Close()
		requestLogFile = nil
	}

	dir := getRequestLogDir()
	logPath := filepath.Join(dir, fmt.Sprintf("%s.jsonl", today))
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open request log file %s: %w", logPath, err)
	}

	requestLogFile = f
	requestLogDate = today
	return f, nil
}

// WriteRequestLog writes a request/response log entry to the JSONL file
func WriteRequestLog(entry *RequestLogEntry) {
	if entry == nil {
		return
	}

	data, err := common.Marshal(entry)
	if err != nil {
		common.SysError(fmt.Sprintf("failed to marshal request log entry: %v", err))
		return
	}

	f, err := getRequestLogFile()
	if err != nil {
		common.SysError(fmt.Sprintf("failed to get request log file: %v", err))
		return
	}

	requestLogMu.Lock()
	defer requestLogMu.Unlock()

	_, err = f.Write(append(data, '\n'))
	if err != nil {
		common.SysError(fmt.Sprintf("failed to write request log: %v", err))
	}
}

// SearchRequestLog searches for a request log entry by request_id
// It searches today's file first, then recent files (up to 7 days back)
func SearchRequestLog(requestId string) (*RequestLogEntry, error) {
	if requestId == "" {
		return nil, fmt.Errorf("request_id is empty")
	}

	dir := getRequestLogDir()

	// Search from today backwards up to 7 days
	for i := 0; i < 7; i++ {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		logPath := filepath.Join(dir, fmt.Sprintf("%s.jsonl", date))

		entry, err := searchInFile(logPath, requestId)
		if err != nil {
			continue // File doesn't exist or can't be read
		}
		if entry != nil {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("request log not found for request_id: %s", requestId)
}

// searchInFile searches for a request_id in a single JSONL file
func searchInFile(filePath string, requestId string) (*RequestLogEntry, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// Increase buffer size for large log entries (up to 10MB per line)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	searchKey := fmt.Sprintf(`"request_id":"%s"`, requestId)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		// Quick string check before doing full unmarshal
		if !strings.Contains(line, searchKey) {
			continue
		}

		var entry RequestLogEntry
		if err := common.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		if entry.RequestId == requestId {
			return &entry, nil
		}
	}

	return nil, nil
}
