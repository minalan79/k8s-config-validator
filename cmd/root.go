/*
Copyright © 2026 k8s-config-validator contributors
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "k8s-config-validator",
	Short: "Validate Kubernetes manifests against best practices",
	Long: `A CLI tool to validate Kubernetes YAML manifests against
schemas and best practices. It checks for common issues, security
misconfigurations, and policy violations.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}


