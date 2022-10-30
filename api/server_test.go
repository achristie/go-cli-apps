package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	todo "github.com/achristie/go-cli-apps/ch1"
)

func setupAPI(t *testing.T) (string, func()) {
	t.Helper()

	tempTodoFile, err := os.CreateTemp("", "todotest")
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(newMux(tempTodoFile.Name()))

	for i := 1; i < 3; i++ {
		var body bytes.Buffer
		taskName := fmt.Sprintf("Task number %d", i)
		item := struct {
			Task string `json:"task"`
		}{
			Task: taskName,
		}

		if err := json.NewEncoder(&body).Encode(item); err != nil {
			t.Fatal(err)
		}

		r, err := http.Post(ts.URL+"/todo", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}

		if r.StatusCode != http.StatusCreated {
			t.Fatalf("Failed to add initial items: Status: %d", r.StatusCode)
		}
	}

	return ts.URL, func() {
		ts.Close()
		os.Remove(tempTodoFile.Name())
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		expCode    int
		expItems   int
		expContent string
	}{
		{name: "GetRoot", path: "/", expCode: http.StatusOK, expContent: "API HERE"},
		{name: "GetAll", path: "/todo", expCode: http.StatusOK, expItems: 2, expContent: "Task number 1"},
		{name: "GetOne", path: "/todo/1", expCode: http.StatusOK, expItems: 1, expContent: "Task number 1"},
		{name: "NotFound", path: "/todo/5000", expCode: http.StatusNotFound},
	}

	url, cleanup := setupAPI(t)
	defer cleanup()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				resp struct {
					Results      todo.List `json:"results"`
					Date         int64     `json:"date"`
					TotalResults int       `json:"total_results"`
				}
				body []byte
				err  error
			)

			r, err := http.Get(url + tt.path)
			if err != nil {
				t.Error(err)
			}
			defer r.Body.Close()

			if r.StatusCode != tt.expCode {
				t.Fatalf("Expected %q, got %q\n", http.StatusText(tt.expCode), http.StatusText(r.StatusCode))
			}

			switch {
			case r.Header.Get("Content-Type") == "application/json":
				if err = json.NewDecoder(r.Body).Decode(&resp); err != nil {
					t.Error(err)
				}
				if resp.TotalResults != tt.expItems {
					t.Errorf("Expected %d items, got %d\n", tt.expItems, resp.TotalResults)
				}

				if resp.Results[0].Task != tt.expContent {
					t.Errorf("Expected %q, got %q", tt.expContent, resp.Results[0].Task)
				}
			case strings.Contains(r.Header.Get("Content-Type"), "text/plain"):
				if body, err = io.ReadAll(r.Body); err != nil {
					t.Error(err)
				}

				if !strings.Contains(string(body), tt.expContent) {
					t.Errorf("Expected %q, got %q\n", tt.expContent, string(body))
				}
			default:
				t.Fatalf("Unsupported Content-Type: %q", r.Header.Get("Content-Type"))
			}

		})
	}
}

func TestAdd(t *testing.T) {
	url, cleanup := setupAPI(t)
	defer cleanup()

	taskName := "Task number 3"
	t.Run("Add", func(t *testing.T) {
		var body bytes.Buffer

		item := struct {
			Task string `json:"task"`
		}{
			Task: taskName,
		}

		if err := json.NewEncoder(&body).Encode(item); err != nil {
			t.Fatal(err)
		}

		r, err := http.Post(url+"/todo", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}

		if r.StatusCode != http.StatusCreated {
			t.Errorf("Expected %q, got %q", http.StatusText(http.StatusCreated), http.StatusText(r.StatusCode))
		}
	})

	t.Run("CheckAdd", func(t *testing.T) {
		r, err := http.Get(url + "/todo/3")
		if err != nil {
			t.Error(err)
		}

		if r.StatusCode != http.StatusOK {
			t.Errorf("Expected %q, got %q", http.StatusText(http.StatusOK), http.StatusText(r.StatusCode))
		}

		var resp todoResponse
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			t.Fatal(err)
		}
		r.Body.Close()

		if resp.Results[0].Task != taskName {
			t.Errorf("Expected %q, got %q\n", taskName, resp.Results[0].Task)
		}
	})
}
