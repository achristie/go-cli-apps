package scan_test

import (
	"errors"
	"os"
	"testing"

	"achristie.net/cobra/scan"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, scan.ErrExists},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hl := &scan.HostsList{}

			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}

			err := hl.Add(tt.host)

			if tt.expectErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil instead\n")
				}

				if !errors.Is(err, tt.expectErr) {
					t.Errorf("expected error %q, got %q\n", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %q instead\n", err)
			}

			if len(hl.Hosts) != tt.expectLen {
				t.Errorf("expected list length %d, got %d\n", tt.expectLen, len(hl.Hosts))
			}

			if hl.Hosts[1] != tt.host {
				t.Errorf("expected host name %q at idx 1, got %q\n", tt.host, hl.Hosts[1])
			}
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"Remove", "host2", 1, nil},
		{"RemoveNotExists", "host1", 1, scan.ErrNotExists},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hl := &scan.HostsList{}

			for _, h := range []string{"host2", "host3"} {
				if err := hl.Add(h); err != nil {
					t.Fatal(err)
				}
			}

			err := hl.Remove(tt.host)

			if tt.expectErr != nil {
				if err == nil {
					t.Fatalf("expected error, got nil instead\n")
				}

				if !errors.Is(err, tt.expectErr) {
					t.Errorf("expected error %q, got %q\n", tt.expectErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %q instead\n", err)
			}

			if len(hl.Hosts) != tt.expectLen {
				t.Errorf("expected list length %d, got %d\n", tt.expectLen, len(hl.Hosts))
			}

			if hl.Hosts[0] == tt.host {
				t.Errorf("host name %q should not be in the list\n", tt.host)
			}
		})
	}
}

func TestSaveLoad(t *testing.T) {
	hl1 := scan.HostsList{}
	hl2 := scan.HostsList{}

	hostName := "host1"
	hl1.Add(hostName)

	tf, err := os.CreateTemp("", "")

	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}

	defer os.Remove(tf.Name())

	if err := hl1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}

	if err := hl2.Load(tf.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}

	if hl1.Hosts[0] != hl2.Hosts[0] {
		t.Errorf("Host %q should match %q host", hl1.Hosts[0], hl2.Hosts[0])
	}
}

func TestLoadNotExists(t *testing.T) {
	hl1 := scan.HostsList{}

	err := hl1.Load("notexists")

	if err != nil {
		t.Fatalf("expected no error; got %q instead\n", err)
	}
}
