package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fanboykun/smokery/apps/core/internal/cli/loader"
	"github.com/fanboykun/smokery/apps/core/internal/spec"
)

// newImportSpecCmd builds: smokery import-spec <spec-file>
//
// Parses an OpenAPI spec and prints the classified operation list. Useful for
// dry-running spec ingestion or seeding a project config.
func newImportSpecCmd(svcs *Services) *cobra.Command {
	var (
		jsonOutput bool
		outPath    string
	)
	cmd := &cobra.Command{
		Use:   "import-spec [spec-file]",
		Short: "Parse an OpenAPI spec and list operations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := loader.OpenAPISpec(args[0])
			if err != nil {
				return err
			}
			analysis, err := spec.Parse(raw)
			if err != nil {
				return err
			}

			if jsonOutput || outPath != "" {
				data, err := json.MarshalIndent(analysis, "", "  ")
				if err != nil {
					return err
				}
				if outPath != "" {
					return os.WriteFile(outPath, data, 0o644)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s v%s\n", analysis.Title, analysis.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "%d operation(s):\n", len(analysis.Operations))
			for _, op := range analysis.Operations {
				marker := " "
				if op.IsDestructive {
					marker = "!"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  %s %-7s %-12s %s   (%s)\n",
					marker, op.Method, op.Classification, op.Path, op.OperationID)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "Write JSON to file")
	return cmd
}
