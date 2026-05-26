// Package cli is the CLI delivery layer. It wires cobra commands to the
// shared app services. This is the only place allowed to import cobra.
package cli

import (
	"github.com/spf13/cobra"

	"github.com/fanboykun/smokery/apps/core/internal/app"
	"github.com/fanboykun/smokery/apps/core/internal/runner"
)

// Services bundles the dependencies a CLI command needs.
type Services struct {
	Project   *app.ProjectService
	Spec      *app.SpecService
	Operation *app.OperationService
	Run       *app.RunService
	Report    *app.ReportService
	Runner    *runner.Runner // direct runner access for one-shot CLI execution
}

// NewRootCmd constructs the top-level `smokery` command with all subcommands.
func NewRootCmd(svcs *Services) *cobra.Command {
	root := &cobra.Command{
		Use:   "smokery",
		Short: "OpenAPI smoke testing CLI",
		Long: "Smokery runs compiled smoke plans against an API. " +
			"It uses the same compiler and runner as the Smokery web platform.",
	}

	root.AddCommand(newRunCmd(svcs))
	root.AddCommand(newCompileCmd(svcs))
	root.AddCommand(newImportSpecCmd(svcs))
	root.AddCommand(newReportCmd(svcs))

	return root
}
