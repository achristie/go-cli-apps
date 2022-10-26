package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()

	tempDir, err := ioutil.TempDir("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for k, n := range files {
		for j := 1; j <= n; j++ {
			fname := fmt.Sprintf("file%d%s", j, k)
			fpath := filepath.Join(tempDir, fname)

			if err := ioutil.WriteFile(fpath, []byte("dummy data"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	return tempDir, func() { os.RemoveAll(tempDir) }
}

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{name: "NoFilter", root: "testdata", cfg: config{ext: "", size: 0, list: true}, expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{name: "FilterExtensionMatch", root: "testdata", cfg: config{ext: ".log", size: 0, list: true}, expected: "testdata/dir.log\n"},
		{name: "FilterExtensionNoMatch", root: "testdata", cfg: config{ext: ".log", size: 10, list: true}, expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeMatch", root: "testdata", cfg: config{ext: ".log", size: 20, list: true}, expected: ""},
		{name: "FilterExtensionNoMatch", root: "testdata", cfg: config{ext: ".gz", size: 0, list: true}, expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer
			if err := run(tt.root, &buffer, tt.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if tt.expected != res {
				t.Errorf("Expected %q, got %q instead\n", tt.expected, res)
			}
		})
	}
}

func TestDelExtension(t *testing.T) {
	tests := []struct {
		name        string
		cfg         config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		expected    string
	}{
		{name: "DeleteExtensionNoMatch", cfg: config{ext: ".log", del: true}, extNoDelete: ".gz", nDelete: 0, nNoDelete: 10, expected: ""},
		{name: "DeleteExtensionMatch", cfg: config{ext: ".log", del: true}, extNoDelete: "", nDelete: 10, nNoDelete: 0, expected: ""},
		{name: "DeleteExtensionMixed", cfg: config{ext: ".log", del: true}, extNoDelete: ".gz", nDelete: 5, nNoDelete: 5, expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)

			tt.cfg.wLog = &logBuffer

			tempDir, cleanup := createTempDir(t, map[string]int{
				tt.cfg.ext:     tt.nDelete,
				tt.extNoDelete: tt.nNoDelete,
			})
			defer cleanup()

			if err := run(tempDir, &buffer, tt.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if tt.expected != res {
				t.Errorf("Expected %q, got %q instead\n", tt.expected, res)
			}

			filesLeft, err := ioutil.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}

			if len(filesLeft) != tt.nNoDelete {
				t.Errorf("Expected %d fiels left, got %d instead\n", tt.nNoDelete, len(filesLeft))
			}

			expLogLines := tt.nDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != expLogLines {
				t.Errorf("Expected %d log lines, got %d instead\n", expLogLines, len(lines))
			}
		})
	}
}

func TestRunArchive(t *testing.T) {
	tests := []struct {
		name         string
		cfg          config
		extNoArchive string
		nArchive     int
		nNoArchive   int
	}{
		{name: "ArchiveExtensionNoMatch", cfg: config{ext: ".log"}, extNoArchive: ".gz", nArchive: 0, nNoArchive: 10},
		{name: "ArchiveExtensionMatch", cfg: config{ext: ".gz"}, extNoArchive: "", nArchive: 10, nNoArchive: 0},
		{name: "ArchiveExtensionMixed", cfg: config{ext: ".log"}, extNoArchive: ".gz", nArchive: 5, nNoArchive: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer

			tempDir, cleanup := createTempDir(t, map[string]int{
				tt.cfg.ext:      tt.nArchive,
				tt.extNoArchive: tt.nNoArchive,
			})
			defer cleanup()

			archiveDir, cleanupArchive := createTempDir(t, nil)
			defer cleanupArchive()

			tt.cfg.archive = archiveDir
			if err := run(tempDir, &buffer, tt.cfg); err != nil {
				t.Fatal(err)
			}

			pattern := filepath.Join(tempDir, fmt.Sprintf("*%s", tt.cfg.ext))
			expFiles, err := filepath.Glob(pattern)
			if err != nil {
				t.Fatal(err)
			}

			expOut := strings.Join(expFiles, "\n")

			res := strings.TrimSpace(buffer.String())

			if expOut != res {
				t.Errorf("Expected %q, got %q instead\n", expOut, res)
			}

			filesArchived, err := os.ReadDir(archiveDir)
			if err != nil {
				t.Fatal(err)
			}

			if len(filesArchived) != tt.nArchive {
				t.Errorf("Expected %d files archived, got %d instead\n", tt.nArchive, len(filesArchived))
			}
		})
	}
}
