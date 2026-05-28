package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/fanboykun/smokery/apps/core/internal/cli/loader"
	"github.com/fanboykun/smokery/apps/core/internal/report"
)

// newRunCmd builds: smokery run <plan-file>
//
// Loads a precompiled SmokePlan from YAML/JSON and executes it via the runner.
// Outputs a CI-friendly summary by default; use --json for the full RunResult.
func newRunCmd(svcs *Services) *cobra.Command {
	var (
		jsonOutput bool
		outputDir  string
	)
	cmd := &cobra.Command{
		Use:   "run [plan-file]",
		Short: "Execute a compiled SmokePlan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			plan, err := loader.SmokePlan(args[0])
			if err != nil {
				return err
			}
			result := svcs.Runner.Execute(context.Background(), plan)

			// Persist artifacts if output dir specified
			if outputDir != "" {
				_ = os.MkdirAll(outputDir, 0755)
				jsonData := report.GenerateJSONArtifact(result)
				_ = os.WriteFile(filepath.Join(outputDir, "result.json"), jsonData, 0644)
				if htmlData, err := report.GenerateHTMLReport(result); err == nil {
					_ = os.WriteFile(filepath.Join(outputDir, "report.html"), htmlData, 0644)
				}
			}

			if jsonOutput {
				b, _ := json.MarshalIndent(result, "", "  ")
				fmt.Fprintln(cmd.OutOrStdout(), string(b))
				if result.Status != "passed" {
					return fmt.Errorf("run failed")
				}
				return nil
			}

			summary := report.GenerateCISummary(result)
			fmt.Fprintf(cmd.OutOrStdout(), "Status:   %s\n", summary.Status)
			fmt.Fprintf(cmd.OutOrStdout(), "Total:    %d\n", summary.Total)
			fmt.Fprintf(cmd.OutOrStdout(), "Passed:   %d\n", summary.Passed)
			fmt.Fprintf(cmd.OutOrStdout(), "Failed:   %d\n", summary.Failed)
			fmt.Fprintf(cmd.OutOrStdout(), "Duration: %dms\n", summary.Duration)
			if len(summary.Failures) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Failures:")
				for _, f := range summary.Failures {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", f)
				}
			}
			if result.Status != "passed" {
				return fmt.Errorf("run failed")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output the full RunResult as JSON")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Directory to save report artifacts (JSON + HTML)")
	return cmd
}
