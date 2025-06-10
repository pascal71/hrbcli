package cmd

import "testing"

func TestParseArtifactRef(t *testing.T) {
	tests := []struct {
		input   string
		project string
		repo    string
		ref     string
		wantErr bool
	}{
		{"proj/repo:tag", "proj", "repo", "tag", false},
		{"proj/repo@sha256:abcd", "proj", "repo", "sha256:abcd", false},
		{"proj/repo", "proj", "repo", "latest", false},
		{"proj/repo:", "", "", "", true},
		{"proj/:tag", "", "", "", true},
		{"proj", "", "", "", true},
	}

	for _, tt := range tests {
		project, repo, ref, err := parseArtifactRef(tt.input)
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
		if project != tt.project || repo != tt.repo || ref != tt.ref {
			t.Errorf("parseArtifactRef(%q) = (%q,%q,%q), want (%q,%q,%q)", tt.input, project, repo, ref, tt.project, tt.repo, tt.ref)
		}
	}
}
