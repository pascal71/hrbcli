package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// requireArgs returns a cobra.PositionalArgs validator that ensures the
// command receives exactly n arguments. The msg parameter should describe
// the expected argument(s) for a friendly error message.
func requireArgs(n int, msg string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("%s", msg)
		}
		return nil
	}
}
