package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/minalan79/k8s-config-validator/validate"
	"github.com/spf13/cobra"
)

var k8sVersion string

var validateCmd = &cobra.Command{
	Use:   "validate [file...]",
	Short: "Validate Kubernetes manifests",
	Long:  `Validate Kubernetes YAML manifests against official Kubernetes schemas. Use "-" to read from stdin.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no files specified, provide file paths or use '-' for stdin")
		}

		v, err := validate.New(k8sVersion)
		if err != nil {
			return fmt.Errorf("failed to initialize validator: %w", err)
		}

		fmt.Printf("Using Kubernetes %s schemas\n\n", v.K8sVersion())

	var exitCode int

		for _, arg := range args {
			var data []byte

			if arg == "-" {
				data, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("failed to read stdin: %w", err)
				}
			} else {
				data, err = os.ReadFile(arg)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", arg, err)
				}
			}

			results, err := v.ValidateDocument(data)
			if err != nil {
				return fmt.Errorf("failed to validate %s: %w", arg, err)
			}

			for _, r := range results {
				resourceID := r.Name
				if r.Namespace != "" {
					resourceID = r.Namespace + "/" + r.Name
				} else {
					resourceID = "default/" + r.Name
				}
				resourceID += fmt.Sprintf(" (%s)", r.Kind)

				if len(r.Errors) > 0 {
					exitCode = 1
					fmt.Printf("[FAIL] %s\n", resourceID)
					for _, e := range r.Errors {
						fmt.Printf("  - %s: %s\n", e.Field, e.Message)
					}
				} else {
					fmt.Printf("[OK] %s\n", resourceID)
				}
			}
		}

		if exitCode != 0 {
			os.Exit(exitCode)
		}

		return nil
	},
}

func init() {
	validateCmd.Flags().StringVar(&k8sVersion, "k8s-version", "", "Kubernetes version to validate against (e.g., 1.27)")
	rootCmd.AddCommand(validateCmd)
}
