package logger

import (
	"context"
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name       string
		level      string
		outputFile string
		wantErr    bool
	}{
		{"valid log level", "info", "", false},
		{"invalid log level", " invalid", "", true},
		{"output to stdout", "info", "stdout", false},
		{"output to file", "info", "stdout", false},
		{"output to file with existing file", "info", "./logs/rpeviewer.log", false},
		{"output to file with non-existent directory", "info", "non-existent-dir/test.log", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := New(ctx, tt.level, tt.outputFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.outputFile != "" && tt.outputFile != "stdout" {
				defer os.Remove(tt.outputFile)
			}
		})
	}
}

func TestNewLoggerWithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := New(ctx, "info", "test.log")
	if err != nil {
		t.Fatal(err)
	}

	cancel()
	_, err = logger.output.(*os.File).Write([]byte("test"))
	if err != nil {
		t.Errorf("expected file to be closed, but got error %v", err)
	}
}

func TestNewLoggerWithExistingFile(t *testing.T) {
	f, err := os.Create("test.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	defer os.Remove("test.log")

	logger, err := New(context.Background(), "info", "test.log")
	if err != nil {
		t.Fatal(err)
	}

	_, err = logger.output.(*os.File).Write([]byte("test"))
	if err != nil {
		t.Errorf("expected file to be writable, but got error %v", err)
	}
}
