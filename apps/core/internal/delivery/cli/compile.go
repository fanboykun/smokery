package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/fanboykun/smokery/apps/core/internal/cli/loader"
	"github.com/fanboykun/smokery/apps/core/internal/compiler"
	"github.com/fanboykun/smokery/apps/core/internal/spec"
)

// newCompileCmd builds: smokery compile --config X --spec Y [-o plan.yaml]
//
// Compiles a ProjectConfig + OpenAPI spec into a SmokePlan. Use this offline
// to produce a plan you can later feed to `smokery run`.
func newCompileCmd(svcs *Services) *cobra.Command {
	var (
		configPath string
		specPath   string
		outPath    string
		outFormat  string
	)
	cmd := &cobra.Command{
		Use:   "compile",
		Short: "Compile a project config + OpenAPI spec into a SmokePlan",
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" || specPath == "" {
				return fmt.Errorf("--config and --spec are required")
			}

			cfg, err := loader.ProjectConfig(configPath)
			if err != nil {
				return err
			}
			rawSpec, err := loader.OpenAPISpec(specPath)
			if err != nil {
				return err
			}
			analysis, err := spec.Parse(rawSpec)
			if err != nil {
				return fmt.Errorf("parse spec: %w", err)
			}

			out := compiler.Compile(compiler.Input{
				Config:     *cfg,
				Operations: analysis.Operations,
			})

			// Print errors and warnings to stderr
			for _, e := range out.Errors {
				fmt.Fprintf(cmd.ErrOrStderr(), "ERROR  %s: %s\n", e.Path, e.Message)
			}
			for _, w := range out.Warnings {
				fmt.Fprintf(cmd.ErrOrStderr(), "WARN   %s: %s\n", w.Path, w.Message)
			}
			if len(out.Errors) > 0 {
				return fmt.Errorf("compile failed: %d error(s)", len(out.Errors))
			}
			if out.Plan == nil {
				return fmt.Errorf("compile produced no plan")
			}

			// Marshal plan to chosen format
			var data []byte
			switch outFormat {
			case "json":
				data, err = json.MarshalIndent(out.Plan, "", "  ")
			case "yaml":
				data, err = yaml.Marshal(out.Plan)
			default:
				return fmt.Errorf("unsupported format %q (use json or yaml)", outFormat)
			}
			if err != nil {
				return err
			}

			if outPath == "" {
				_, err = cmd.OutOrStdout().Write(data)
				return err
			}
			return os.WriteFile(outPath, data, 0o644)
		},
	}
	cmd.Flags().StringVar(&configPath, "config", "", "Path to project config YAML/JSON")
	cmd.Flags().StringVar(&specPath, "spec", "", "Path to OpenAPI spec file")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "Output file (default stdout)")
	cmd.Flags().StringVar(&outFormat, "format", "yaml", "Output format: yaml or json")
	return cmd
}
