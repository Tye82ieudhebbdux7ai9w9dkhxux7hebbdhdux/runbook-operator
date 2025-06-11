package outputs

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	runbookv1alpha1 "github.com/guibes/runbook-operator/api/v1alpha1"
)

type HTMLOutput struct {
	BasePath string
}

const htmlTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>{{.Spec.AlertName}} Runbook</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .header { background: #f5f5f5; padding: 15px; border-radius: 5px; margin-bottom: 20px; }
        .severity-{{.Spec.Severity}} { border-left: 4px solid #ff6b6b; }
        .severity-warning { border-left: 4px solid #feca57; }
        .severity-info { border-left: 4px solid #48cae4; }
        .step { background: #f8f9fa; padding: 10px; margin: 10px 0; border-radius: 3px; }
        code { background: #e9ecef; padding: 2px 4px; border-radius: 3px; }
        pre { background: #2d3748; color: #e2e8f0; padding: 15px; border-radius: 5px; overflow-x: auto; }
    </style>
</head>
<body>
    <div class="header severity-{{.Spec.Severity}}">
        <h1>🚨 {{.Spec.AlertName}}</h1>
        <p><strong>Severity:</strong> {{.Spec.Severity}} | <strong>Team:</strong> {{.Spec.Team}}</p>
        <p><em>Generated: {{.GeneratedAt}}</em></p>
    </div>
    
    <h2>💥 Impact</h2>
    <p>{{.Spec.Content.Impact}}</p>
    
    <h2>🔍 Investigation Steps</h2>
    {{range $i, $step := .Spec.Content.Investigation}}
    <div class="step">
        <h3>Step {{add $i 1}}: {{.Description}}</h3>
        {{if .Command}}<pre>{{.Command}}</pre>{{end}}
        {{if .Expected}}<p><strong>Expected:</strong> {{.Expected}}</p>{{end}}
    </div>
    {{end}}
    
    <h2>🛠️ Remediation</h2>
    {{range $i, $step := .Spec.Content.Remediation}}
    <div class="step">
        <h3>{{add $i 1}}. {{.Description}} {{if .Risk}}(Risk: {{.Risk}}){{end}}</h3>
        {{if .Command}}<pre>{{.Command}}</pre>{{end}}
    </div>
    {{end}}
    
    <h2>🛡️ Prevention</h2>
    <p>{{.Spec.Content.Prevention}}</p>
</body>
</html>`

func (h *HTMLOutput) Generate(runbook *runbookv1alpha1.Runbook) error {
	tmpl := template.Must(template.New("runbook").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}).Parse(htmlTemplate))

	filename := fmt.Sprintf("%s.html", runbook.Spec.AlertName)
	fullPath := filepath.Join(h.BasePath, filename)

	if err := os.MkdirAll(h.BasePath, 0755); err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := struct {
		*runbookv1alpha1.Runbook
		GeneratedAt string
	}{
		Runbook:     runbook,
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	return tmpl.Execute(file, data)
}
