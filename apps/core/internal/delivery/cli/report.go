package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/report"
)

// newReportCmd builds: smokery report <result-file> --view <debug|ci|mermaid>
//
// Loads a saved RunResult JSON and renders the chosen view.
func newReportCmd(svcs *Services) *cobra.Command {
	var view string
	cmd := &cobra.Command{
		Use:   "report [result-file]",
		Short: "Render a report view from a saved RunResult",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			var result model.RunResult
			if err := json.Unmarshal(data, &result); err != nil {
				return fmt.Errorf("parse result: %w", err)
			}

			switch view {
			case "debug":
				v := report.GenerateDebugView(&result)
				out, _ := json.MarshalIndent(v, "", "  ")
				fmt.Fprintln(cmd.OutOrStdout(), string(out))
			case "ci":
				v := report.GenerateCISummary(&result)
				out, _ := json.MarshalIndent(v, "", "  ")
				fmt.Fprintln(cmd.OutOrStdout(), string(out))
			case "mermaid":
				fmt.Fprintln(cmd.OutOrStdout(), report.GenerateMermaidDiagram(&result))
			default:
				return fmt.Errorf("unknown view %q (use debug|ci|mermaid)", view)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&view, "view", "ci", "Report view: debug | ci | mermaid")
	return cmd
}
