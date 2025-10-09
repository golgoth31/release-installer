package progressbar

import (
	"io"
	"strings"
	"testing"
	"time"
)

func TestSimpleProgress(t *testing.T) {
	// Test the SimpleProgress function
	SimpleProgress(50, 100, "Test")
	time.Sleep(100 * time.Millisecond)
	SimpleProgress(100, 100, "Test")
	time.Sleep(100 * time.Millisecond)
}

func TestCreateSimpleBar(t *testing.T) {
	tests := []struct {
		percent  float64
		width    int
		expected string
	}{
		{0.0, 10, "[░░░░░░░░░░]"},
		{0.5, 10, "[█████░░░░░]"},
		{1.0, 10, "[██████████]"},
		{0.25, 4, "[█░░░]"},
	}

	for _, test := range tests {
		result := createSimpleBar(test.percent, test.width)
		if result != test.expected {
			t.Errorf("createSimpleBar(%.2f, %d) = %s, expected %s",
				test.percent, test.width, result, test.expected)
		}
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, test := range tests {
		result := formatBytes(test.bytes)
		if result != test.expected {
			t.Errorf("formatBytes(%d) = %s, expected %s",
				test.bytes, result, test.expected)
		}
	}
}

func TestProgressReader(t *testing.T) {
	// Create a test reader
	testData := "Hello, World!"
	reader := strings.NewReader(testData)

	// Create a progress bar
	pb := New()

	// Track progress
	progressReader := pb.TrackProgress("test", 0, int64(len(testData)), io.NopCloser(reader))

	// Read all data
	buf := make([]byte, len(testData))
	n, err := progressReader.Read(buf)

	if err != nil && err != io.EOF {
		t.Errorf("Unexpected error: %v", err)
	}

	if n != len(testData) {
		t.Errorf("Expected to read %d bytes, got %d", len(testData), n)
	}

	// Close the progress reader
	err = progressReader.Close()
	if err != nil {
		t.Errorf("Unexpected error on close: %v", err)
	}
}
