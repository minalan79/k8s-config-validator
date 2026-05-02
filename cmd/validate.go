package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/minalan79/k8s-config-validator/validate"
	"github.com/spf13/cobra"
)

var k8sVersion string
var recursive bool

var validateCmd = &cobra.Command{
	Use:   "validate [file|dir...]",
	Short: "Validate Kubernetes manifests",
	Long:  `Validate Kubernetes YAML manifests against official Kubernetes schemas. Use "-" to read from stdin. Directories will be scanned for YAML files.`,
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
			var files []string

			if arg == "-" {
				files = append(files, "-")
			} else {
				info, err := os.Stat(arg)
				if err != nil {
					return fmt.Errorf("failed to access %s: %w", arg, err)
				}

				if info.IsDir() {
					files, err = findManifests(arg, recursive)
					if err != nil {
						return fmt.Errorf("failed to scan directory %s: %w", arg, err)
					}
					if len(files) == 0 {
						fmt.Printf("No YAML files found in %s\n", arg)
						continue
					}
					fmt.Printf("Found %d manifest(s) in %s\n\n", len(files), arg)
				} else {
					files = append(files, arg)
				}
			}

		for _, file := range files {
			var data []byte

			if file == "-" {
				data, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("failed to read stdin: %w", err)
				}
			} else {
				data, err = os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", file, err)
				}
			}

			results, err := v.ValidateDocument(data)
			if err != nil {
				return fmt.Errorf("failed to validate %s: %w", file, err)
			}

			if len(files) > 1 || len(args) > 0 {
				fmt.Printf("--- %s ---\n", file)
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
			fmt.Println()
		}
		}

		if exitCode != 0 {
			os.Exit(exitCode)
		}

		return nil
	},
}

func findManifests(dir string, recurse bool) ([]string, error) {
	var manifests []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if !recurse && path != dir {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".yaml" || ext == ".yml" {
			manifests = append(manifests, path)
		}

		return nil
	})

	return manifests, err
}

func init() {
	validateCmd.Flags().StringVar(&k8sVersion, "k8s-version", "", "Kubernetes version to validate against (e.g., 1.27)")
	validateCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively scan subdirectories for manifests")
	rootCmd.AddCommand(validateCmd)
}
