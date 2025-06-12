---
layout: default
title: Home
---

# RunbookOperator

A cloud-native Kubernetes operator that automatically generates, manages, and distributes runbook documentation from PrometheusRule configurations.

## Why RunbookOperator?

Manual runbook creation and maintenance is time-consuming and error-prone. RunbookOperator solves this by:

- **Automating runbook generation** from PrometheusRule annotations
- **Supporting multiple output formats** (Markdown, HTML, ConfigMaps, APIs)
- **Ensuring consistency** across teams and alert types
- **Integrating with GitOps** workflows seamlessly
- **Validating runbook content** automatically

## Quick Start

    # Install the operator
    kubectl apply -f https://raw.githubusercontent.com/guibes/runbook-operator/main/config/default

    # Create a runbook
    kubectl apply -f - <<EOF
    apiVersion: runbook.runbook.io/v1alpha1
    kind: Runbook
    metadata:
      name: example-runbook
    spec:
      alertName: "HighMemoryUsage"
      severity: "warning"
      team: "platform"
      outputs:
        - format: "markdown"
          destination: "/tmp/runbooks"
    EOF

## Features

### Automated Generation
Runbooks are automatically created and updated when PrometheusRule resources change, ensuring your documentation is always current.

### Multiple Output Formats
- **Markdown**: For documentation repositories
- **HTML**: For web-based dashboards  
- **ConfigMaps**: For Kubernetes-native storage
- **API**: For integration with external systems

### Template System
Customize runbook formats for different alert types and teams using Go templates.

### GitOps Integration
Seamlessly integrate with existing GitOps workflows by automatically committing generated runbooks to Git repositories.

### Dark Mode Support
Modern dark mode with automatic system preference detection and manual toggle.

## Documentation

<div class="nav-toc">
<h4>Documentation Sections</h4>
<ul>
  <li><a href="{{ '/getting-started/' | relative_url }}">Getting Started</a> - Installation and basic usage</li>
  <li><a href="{{ '/api-reference/' | relative_url }}">API Reference</a> - Complete API documentation</li>
  <li><a href="{{ '/examples/' | relative_url }}">Examples</a> - Practical usage examples</li>
  <li><a href="{{ '/developer-guide/' | relative_url }}">Developer Guide</a> - Development and contribution guide</li>
  <li><a href="{{ '/configuration/' | relative_url }}">Configuration</a> - Configuration options and reference</li>
</ul>
</div>

## Community

- **GitHub**: [guibes/runbook-operator](https://github.com/guibes/runbook-operator)
- **Issues**: [Report bugs or request features](https://github.com/guibes/runbook-operator/issues)
- **Discussions**: [Community discussions](https://github.com/guibes/runbook-operator/discussions)

## License

RunbookOperator is licensed under the Apache 2.0 License. See [LICENSE](https://github.com/guibes/runbook-operator/blob/main/LICENSE) for details.
