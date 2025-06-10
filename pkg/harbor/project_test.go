package harbor

import "testing"

func TestParseStorageLimit(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		wantErr  bool
	}{
		{"", -1, false},
		{"-1", -1, false},
		{"10G", 10 * 1024 * 1024 * 1024, false},
		{"2T", 2 * 1024 * 1024 * 1024 * 1024, false},
		{"5M", 5 * 1024 * 1024, false},
		{"123", 123, false},
		{"bad", 0, true},
	}

	for _, tt := range tests {
		got, err := ParseStorageLimit(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("expected error for %q", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error for %q: %v", tt.input, err)
			continue
		}
		if got != tt.expected {
			t.Errorf("ParseStorageLimit(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}
