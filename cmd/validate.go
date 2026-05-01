package cmd

import (
	"fmt"

	"github.com/minalan79/k8s-config-validator/validate"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Kubernetes manifests",
	Long:  `Validate Kubernetes manifests against schemas and best practices`,
	Run: func(cmd *cobra.Command, args []string) {
		sample := []byte(`
apiVersion: apps/v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: nginx
    image: nginx:latest
---
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: nginx
    image: nginx:latest
`)
		manifests, err := validate.Validate(sample)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Parsed %d manifests: %#v\n", len(manifests), manifests)
		fmt.Println(manifests[0].APIVersion) // Output: apps/v1
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
