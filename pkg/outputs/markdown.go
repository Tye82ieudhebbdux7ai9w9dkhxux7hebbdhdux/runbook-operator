package outputs

import (
	"fmt"
	"os"
	"path/filepath"

	runbookv1alpha1 "github.com/guibes/runbook-operator/api/v1alpha1"
)

type MarkdownOutput struct {
	BasePath string
}

func (m *MarkdownOutput) Generate(runbook *runbookv1alpha1.Runbook, content string) error {
	filename := fmt.Sprintf("%s.md", runbook.Spec.AlertName)
	fullPath := filepath.Join(m.BasePath, filename)

	// Ensure directory exists
	if err := os.MkdirAll(m.BasePath, 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, []byte(content), 0644)
}
