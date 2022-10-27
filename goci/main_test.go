package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	var tests = []struct {
		name   string
		proj   string
		out    string
		expErr error
	}{
		{name: "success", proj: "./testdata/tool/", out: "Go Build: SUCCESS\n", expErr: nil},
		{name: "fail", proj: "./testdata/toolErr/", out: "", expErr: &stepErr{step: "go build"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			err := run(tt.proj, &out)

			if tt.expErr != nil {
				if err == nil {
					t.Errorf("Expected error %q. Got 'nil' instead", tt.expErr)
					return
				}
				if !errors.Is(err, tt.expErr) {
					t.Errorf("expected error %q, got %q", tt.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}

			if out.String() != tt.out {
				t.Errorf("Expected output: %q. Got %q", tt.out, out.String())
			}
		})
	}
}
