# k8s-config-validator

A CLI tool to validate Kubernetes YAML manifests against official Kubernetes schemas. It checks for common issues, schema violations, and configuration errors.

## Features

- Validate individual files or entire directories
- Recursive directory scanning with `-r` flag
- Multi-document YAML support
- Read from stdin using `-`
- Target specific Kubernetes versions
- Sectioned output grouped by file for easy troubleshooting

## Installation

```bash
go build -o k8s-config-validator .
```

Or download a pre-built binary from the releases page.

## Usage

### Validate a single file

```bash
./k8s-config-validator validate manifest.yaml
```

### Validate multiple files

```bash
./k8s-config-validator validate manifest1.yaml manifest2.yaml
```

### Validate all manifests in a directory

```bash
./k8s-config-validator validate ./manifests/
```

### Recursively scan directories and subdirectories

```bash
./k8s-config-validator validate -r ./kubernetes/
```

### Read from stdin

```bash
cat manifest.yaml | ./k8s-config-validator validate -
```

### Target a specific Kubernetes version

```bash
./k8s-config-validator validate --k8s-version 1.28 manifest.yaml
```

## Flags

| Flag | Description |
|------|-------------|
| `-h, --help` | Show help for the validate command |
| `--k8s-version string` | Kubernetes version to validate against (e.g., 1.27) |
| `-r, --recursive` | Recursively scan subdirectories for manifests |

## Output

Results are sectioned by file path for easy troubleshooting:

```
Using Kubernetes 1.30 schemas

Found 2 manifest(s) in example-manifests

--- example-manifests/dep.yml ---
[OK] default/nginx-deployment (Deployment)
[OK] default/nginx-service (Service)

--- example-manifests/test-invalid.yml ---
[FAIL] default/ (Pod)
  - : Pod.core "" is invalid: metadata.name: Required value: name or generateName is required
[OK] default/nginx-service (Service)
```

- `[OK]` indicates a valid resource
- `[FAIL]` indicates a resource with validation errors

The tool exits with code `1` if any validation failures are found.

## Project Structure

```
.
├── cmd/
│   ├── root.go       # CLI root command
│   └── validate.go   # Validate subcommand
├── validate/
│   ├── root.go       # Validator logic and result types
│   └── parser.go     # YAML multi-document splitting
├── example-manifests/ # Example YAML manifests
├── main.go           # Entry point
├── go.mod
└── go.sum
```

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [kubectl-validate](https://github.com/kubernetes/kubectl-validate) - Kubernetes schema validation
- [yaml.v3](https://github.com/go-yaml/yaml) - YAML parsing

## License

Apache-2.0
