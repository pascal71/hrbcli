package api

// VulnerabilityReport represents a vulnerability scanning report
// Only fields needed by the CLI are included.
type VulnerabilityReport struct {
	Vulnerabilities []VulnerabilityItem `json:"vulnerabilities"`
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
	ReportID    string `json:"report_id"`
	ScanStatus  string `json:"scan_status"`
	Severity    string `json:"severity"`
	CompletePct int    `json:"complete_percent"`
}
