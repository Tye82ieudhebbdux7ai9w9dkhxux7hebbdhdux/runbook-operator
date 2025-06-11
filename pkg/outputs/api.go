package outputs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	runbookv1alpha1 "github.com/guibes/runbook-operator/api/v1alpha1"
)

type APIOutput struct {
	BaseURL string
	ApiKey  string
}

type RunbookAPI struct {
	ID          string                 `json:"id"`
	AlertName   string                 `json:"alert_name"`
	Severity    string                 `json:"severity"`
	Team        string                 `json:"team"`
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata"`
	GeneratedAt time.Time              `json:"generated_at"`
}

func (a *APIOutput) Generate(runbook *runbookv1alpha1.Runbook, content string) error {
	apiData := RunbookAPI{
		ID:        fmt.Sprintf("%s-%s", runbook.Namespace, runbook.Name),
		AlertName: runbook.Spec.AlertName,
		Severity:  runbook.Spec.Severity,
		Team:      runbook.Spec.Team,
		Content:   content,
		Metadata: map[string]interface{}{
			"namespace": runbook.Namespace,
			"outputs":   runbook.Spec.Outputs,
		},
		GeneratedAt: time.Now(),
	}

	jsonData, err := json.Marshal(apiData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/runbooks", a.BaseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if a.ApiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.ApiKey))
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	return nil
}
