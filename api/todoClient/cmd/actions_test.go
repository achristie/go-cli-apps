//go:build !integration
// +build !integration

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestListAction(t *testing.T) {
	tests := []struct {
		name     string
		expError error
		expOut   string
		resp     struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{name: "Results", expError: nil, expOut: "-  1  Task 1\n-  2  Task 2\n", resp: testResponses["resultsMany"]},
		{name: "NoResults", expError: ErrNotFound, resp: testResponses["noResults"]},
		{name: "InvalidURL", expError: ErrConnection, resp: testResponses["noResults"], closeServer: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.resp.Status)
				fmt.Fprintln(w, tt.resp.Body)
			})
			defer cleanup()

			if tt.closeServer {
				cleanup()
			}

			var out bytes.Buffer
			err := listAction(&out, url)

			if tt.expError != nil {
				if err == nil {
					t.Fatalf("Expected error %q, got nothing", tt.expError)
				}

				if !errors.Is(err, tt.expError) {
					t.Errorf("Expected error %q, got %q", tt.expError, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %q", err)
			}

			if tt.expOut != out.String() {
				t.Errorf("expected output %q, got %q", tt.expOut, out.String())
			}
		})
	}
}

func TestViewAction(t *testing.T) {
	tests := []struct {
		name     string
		expError error
		expOut   string
		resp     struct {
			Status int
			Body   string
		}
		id string
	}{
		{name: "ResultsOne", expError: nil, expOut: "Task:         Task 1\nCreated at:   Oct/28 @ 08:23\nCompleted:    No\n", resp: testResponses["resultsOne"], id: "1"},
		{name: "NotFound", expError: ErrNotFound, resp: testResponses["notFound"], id: "1"},
		{name: "InvalidID", expError: ErrNotNumber, resp: testResponses["noResults"], id: "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, cleanup := mockServer(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.resp.Status)
					fmt.Fprintln(w, tt.resp.Body)
				})
			defer cleanup()

			var out bytes.Buffer

			err := viewAction(&out, url, tt.id)

			if tt.expError != nil {
				if err == nil {
					t.Fatalf("expected error %q, got no error", tt.expError)
				}

				if !errors.Is(err, tt.expError) {
					t.Errorf("Expected error %q, got %q", tt.expError, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %q", err)
			}

			if tt.expOut != out.String() {
				t.Errorf("expected output %q, got %q", tt.expOut, out.String())
			}
		})
	}
}

func TestAddAction(t *testing.T) {
	expURLPath := "/todo"
	expMethod := http.MethodPost
	expBody := "{\"task\":\"Task 1\"}\n"
	expContentType := "application/json"
	expOut := "Added task \"Task 1\" to the list\n"
	args := []string{"Task", "1"}

	url, cleanup := mockServer(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != expURLPath {
				t.Errorf("Expected path %q got %q\n", expURLPath, r.URL.Path)
			}

			if r.Method != expMethod {
				t.Errorf("Expected method %q got %q\n", expMethod, r.Method)
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			r.Body.Close()

			if string(body) != expBody {
				t.Errorf("Expected body %q, got %q", expBody, string(body))
			}

			contentType := r.Header.Get("Content-Type")
			if contentType != expContentType {
				t.Errorf("Expected Content-Type %q, got %q", expContentType, contentType)
			}

			w.WriteHeader(testResponses["created"].Status)
			fmt.Fprintln(w, testResponses["created"].Body)
		})
	defer cleanup()

	var out bytes.Buffer
	if err := addAction(&out, url, args); err != nil {
		t.Fatalf("expected no error, got %q", err)
	}

	if expOut != out.String() {
		t.Errorf("expected output %q, got %q", expOut, out.String())
	}
}
