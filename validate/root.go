package validate

import (
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
	"sigs.k8s.io/kubectl-validate/pkg/validator"
)

type Result struct {
	GVK        schema.GroupVersionKind
	Name       string
	Namespace  string
	Kind       string
	APIVersion string
	Errors     []FieldError
}

type FieldError struct {
	Type    string
	Field   string
	Message string
}

func New(k8sVersion string) (*Validator, error) {
	if k8sVersion == "" {
		k8sVersion = openapiclient.HardcodedBuiltinVersions[len(openapiclient.HardcodedBuiltinVersions)-1]
	}

	client := openapiclient.NewHardcodedBuiltins(k8sVersion)
	v, err := validator.New(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create validator: %w", err)
	}

	return &Validator{validator: v, k8sVersion: k8sVersion}, nil
}

type Validator struct {
	validator  *validator.Validator
	k8sVersion string
}

func (v *Validator) K8sVersion() string {
	return v.k8sVersion
}

func (v *Validator) ValidateDocument(data []byte) ([]Result, error) {
	docs, err := SplitDocuments(data)
	if err != nil {
		return nil, fmt.Errorf("failed to split documents: %w", err)
	}

	var results []Result
	for _, doc := range docs {
		gvk, obj, err := v.validator.Parse(doc)
		if err != nil {
			results = append(results, Result{
				GVK:  gvk,
				Errors: []FieldError{{
					Type:    "parse",
					Field:   "",
					Message: err.Error(),
				}},
			})
			continue
		}

		validationErr := v.validator.Validate(obj)
		result := Result{
			GVK:        gvk,
			Name:       obj.GetName(),
			Namespace:  obj.GetNamespace(),
			Kind:       obj.GetKind(),
			APIVersion: obj.GetAPIVersion(),
		}

		if validationErr != nil {
			result.Errors = parseValidationErrors(validationErr, obj)
		}

		results = append(results, result)
	}

	return results, nil
}

func parseValidationErrors(err error, obj *unstructured.Unstructured) []FieldError {
	var fieldErrors []FieldError
	errStr := err.Error()

	if isMultiError(errStr) {
		for _, line := range splitErrors(errStr) {
			fieldErrors = append(fieldErrors, parseSingleError(line, obj))
		}
	} else {
		fieldErrors = append(fieldErrors, parseSingleError(errStr, obj))
	}

	return fieldErrors
}

func isMultiError(errStr string) bool {
	return len(errStr) > 0 && errStr[0] == '['
}

func splitErrors(errStr string) []string {
	errStr = errStr[1 : len(errStr)-1]
	var parts []string
	depth := 0
	start := 0
	for i := 0; i < len(errStr); i++ {
		switch errStr[i] {
		case '[':
			depth++
		case ']':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, errStr[start:i])
				start = i + 2
			}
		}
	}
	if start < len(errStr) {
		parts = append(parts, errStr[start:])
	}
	return parts
}

func parseSingleError(errStr string, obj *unstructured.Unstructured) FieldError {
	return FieldError{
		Type:    "invalid",
		Field:   extractField(errStr),
		Message: errStr,
	}
}

func extractField(errStr string) string {
	if idx := extractFieldWithRegex(errStr, "spec."); idx >= 0 {
		return "spec." + errStr[idx:extractFieldEnd(errStr[idx:])]
	}
	return ""
}

func extractFieldWithRegex(errStr, prefix string) int {
	for i := 0; i <= len(errStr)-len(prefix); i++ {
		if errStr[i:i+len(prefix)] == prefix {
			return i
		}
	}
	return -1
}

func extractFieldEnd(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ':' || s[i] == ' ' {
			return i
		}
	}
	return len(s)
}
