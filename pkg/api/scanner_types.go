package api

// VulnerabilityReport represents a vulnerability scanning report
// Only fields needed by the CLI are included.
type VulnerabilityReport struct {
	Severity        string               `json:"severity"`
	Summary         VulnerabilitySummary `json:"summary"`
	Vulnerabilities []VulnerabilityItem  `json:"vulnerabilities"`
}

// VulnerabilitySummary contains summary information about the
// vulnerabilities found in a scan.
type VulnerabilitySummary struct {

	Total   int            `json:"total"`
	Fixable int            `json:"fixable"`
	Summary map[string]int `json:"summary"`

}

// VulnerabilityItem represents a single vulnerability entry
// returned by Harbor scanners.
type VulnerabilityItem struct {
	CVEID        string `json:"id"`
	Package      string `json:"package"`
	Version      string `json:"version"`
	FixedVersion string `json:"fix_version"`
	Severity     string `json:"severity"`
}

// NativeReportSummary represents the summary of a scan report
// attached to an artifact. Only fields relevant for displaying
// running scans are included.
type NativeReportSummary struct {
	ReportID    string               `json:"report_id"`
	ScanStatus  string               `json:"scan_status"`
	Severity    string               `json:"severity"`
	CompletePct int                  `json:"complete_percent"`
	Summary     VulnerabilitySummary `json:"summary"`
}
