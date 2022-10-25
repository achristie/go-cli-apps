package main

import (
	"os"
	"testing"
)

func TestFilterOut(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		ext      string
		minSize  int64
		expected bool
	}{
		{"FilterNoExtension", "testdata/dir.log", "", 0, false},
		{"FilterExtensionMatch", "testdata/dir.log", ".log", 0, false},
		{"FilterExtensionNoMatch", "testdata/dir.log", ".sh", 0, true},
		{"FilterExtensionSizeMatch", "testdata/dir.log", ".log", 10, false},
		{"FilterExtensionSizeNoMatch", "testdata/dir.log", ".log", 20, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := os.Stat(tt.file)
			if err != nil {
				t.Fatal(err)
			}

			f := filterOut(tt.file, tt.ext, tt.minSize, info)

			if f != tt.expected {
				t.Errorf("Expected '%t', got '%t' instead\n", tt.expected, f)
			}
		})
	}
}
