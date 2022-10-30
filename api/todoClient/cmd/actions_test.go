package cmd

import (
	"bytes"
	"errors"
	"fmt"
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
