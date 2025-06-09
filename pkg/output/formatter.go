package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

var (
	format  = "table"
	noColor = false
	debug   = false
)

// SetFormat sets the output format
func SetFormat(f string) {
	format = f
}

// GetFormat returns the current output format
func GetFormat() string {
	return format
}

// SetNoColor sets whether to disable colored output
func SetNoColor(nc bool) {
	noColor = nc
	if nc {
		color.NoColor = true
	}
}

// SetDebug sets debug mode
func SetDebug(d bool) {
	debug = d
}

// JSON outputs data as JSON
func JSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// YAML outputs data as YAML
func YAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	return encoder.Encode(data)
}

// Table creates a new table writer
func Table() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	// Just use the basic table configuration
	return table
}

// Success prints a success message
func Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if noColor {
		fmt.Println("✓", msg)
	} else {
		fmt.Println(color.GreenString("✓"), msg)
	}
}

// Error prints an error message
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if noColor {
		fmt.Fprintln(os.Stderr, "✗", msg)
	} else {
		fmt.Fprintln(os.Stderr, color.RedString("✗"), msg)
	}
}

// Warning prints a warning message
func Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if noColor {
		fmt.Println("!", msg)
	} else {
		fmt.Println(color.YellowString("!"), msg)
	}
}

// Info prints an info message
func Info(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

// Debug prints a debug message if debug mode is enabled
func Debug(format string, args ...interface{}) {
	if debug {
		msg := fmt.Sprintf(format, args...)
		if noColor {
			fmt.Fprintln(os.Stderr, "[DEBUG]", msg)
		} else {
			fmt.Fprintln(os.Stderr, color.HiBlackString("[DEBUG]"), msg)
		}
	}
}

// FormatSize converts a byte count to a human readable string.
// Use this helper throughout the CLI when reporting storage quotas.
// Negative sizes return "unlimited".
func FormatSize(bytes int64) string {
	if bytes < 0 {
		return "unlimited"
	}

	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Truncate truncates a string to a maximum length
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// Bold returns a bold string (if color is enabled)
func Bold(s string) string {
	if noColor {
		return s
	}
	return color.New(color.Bold).Sprint(s)
}

// Green returns a green string (if color is enabled)
func Green(s string) string {
	if noColor {
		return s
	}
	return color.GreenString(s)
}

// Red returns a red string (if color is enabled)
func Red(s string) string {
	if noColor {
		return s
	}
	return color.RedString(s)
}

// Yellow returns a yellow string (if color is enabled)
func Yellow(s string) string {
	if noColor {
		return s
	}
	return color.YellowString(s)
}

// PrintList prints a list of items with proper formatting
func PrintList(title string, items []string) {
	if len(items) == 0 {
		Info("No %s found", strings.ToLower(title))
		return
	}

	Info("%s:", title)
	for _, item := range items {
		fmt.Printf("  • %s\n", item)
	}
}
