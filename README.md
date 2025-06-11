
# RunbookOperator

A cloud-native Kubernetes operator that automatically generates, manages, and distributes runbook documentation from PrometheusRule configurations.

## The Problem

Manual runbook creation and maintenance is time-consuming, error-prone, and often results in outdated documentation when you need it most. Traditional approaches suffer from:

- Runbooks that become stale and disconnected from actual alert definitions
- Inconsistent documentation formats across teams
- Manual effort required to keep runbooks synchronized with monitoring rules
- Lack of standardized incident response procedures
- No automated validation of runbook procedures

## The Solution

RunbookOperator bridges the gap between alerting and incident response by automatically generating comprehensive runbooks directly from your PrometheusRule annotations. It provides:

- **Automated Generation**: Runbooks are created and updated automatically when PrometheusRule resources change
- **Multiple Output Formats**: Generate Markdown, HTML, ConfigMaps, and integrate with external APIs
- **Validation**: Built-in validation ensures runbooks are complete and accurate
- **GitOps Integration**: Seamlessly integrates with existing GitOps workflows
- **Template System**: Customizable templates for different alert types and teams

## Quick Example

Define your alert with runbook content in your PrometheusRule:
```yaml
    apiVersion: monitoring.coreos.com/v1
    kind: PrometheusRule
    metadata:
      name: app-alerts
    spec:
      groups:
      - name: application
        rules:
        - alert: HighMemoryUsage
          expr: memory_usage > 0.8
          annotations:
            runbook: |
              ## Impact
              Application performance degrades when memory exceeds 80%
              
              ## Investigation
              1. Check pod memory: kubectl top pods -l app=myapp
              2. Review logs: kubectl logs -l app=myapp --tail=100
              
              ## Remediation
              1. Restart pod: kubectl delete pod -l app=myapp
              2. Scale horizontally: kubectl scale deployment myapp --replicas=3
```
The operator automatically creates a Runbook resource and generates documentation:
```yaml
    apiVersion: runbook.runbook.io/v1alpha1
    kind: Runbook
    metadata:
      name: high-memory-usage
    spec:
      alertName: "HighMemoryUsage"
      severity: "warning"
      team: "platform"
      outputs:
        - format: "markdown"
          destination: "/tmp/runbooks"
        - format: "html"
          destination: "/tmp/runbooks/html"
```
## Prerequisites

- Kubernetes v1.11.3+ cluster
- Go 1.21+ (for development)
- Docker 17.03+ (for building images)
- kubectl v1.11.3+

## Installation

### Quick Start with Kind

Create development cluster:
```bash
    kind create cluster --name runbook-dev
```

Install the operator:
```bash
    kubectl apply -f https://raw.githubusercontent.com/guibes/runbook-operator/main/config/default
```
Apply a sample runbook:
```bash
    kubectl apply -f https://raw.githubusercontent.com/guibes/runbook-operator/main/config/samples/runbook_v1alpha1_runbook.yaml
```
### Production Installation

1. **Install CRDs:**
```bash
       kubectl apply -f config/crd/bases/
```
2. **Deploy the operator:**
```bash
       make deploy IMG=your-registry/runbook-operator:latest
```
3. **Create storage for runbook outputs:**
```bash
       kubectl apply -f config/storage/
```
## Usage

### Basic Runbook Creation

Create a runbook resource:
```yaml
    apiVersion: runbook.runbook.io/v1alpha1
    kind: Runbook
    metadata:
      name: database-alert
    spec:
      alertName: "DatabaseDown"
      severity: "critical"
      team: "database"
      content:
        impact: "Database is unavailable, affecting all services"
        investigation:
          - description: "Check database pod status"
            command: "kubectl get pods -l app=database"
            expected: "All pods should be Running"
        remediation:
          - description: "Restart database service"
            command: "kubectl rollout restart deployment/database"
            risk: "medium"
            automated: false
        prevention: "Implement database health monitoring"
      outputs:
        - format: "markdown"
          destination: "/tmp/runbooks"
        - format: "html"
          destination: "/tmp/runbooks/html"
```
### Output Formats

The operator supports multiple output formats:

- **Markdown**: Standard markdown files for documentation
- **HTML**: Web-ready HTML with styling
- **API**: Send to external systems via REST APIs

### Template System

Create custom templates for different alert types:
```yaml
    apiVersion: runbook.runbook.io/v1alpha1
    kind: RunbookTemplate
    metadata:
      name: database-template
    spec:
      name: "database-incidents"
      description: "Template for database-related incidents"
      template: |
        # Database Alert: {{ .Spec.AlertName }}
        
        **Severity**: {{ .Spec.Severity }}
        **Team**: {{ .Spec.Team }}
        
        ## Database-Specific Checks
        - Connection pool status
        - Replication lag
        - Query performance
```
## Development

### Local Development Setup

1. **Clone and setup:**
```bash
       git clone https://github.com/guibes/runbook-operator.git
       cd runbook-operator
```
2. **Install dependencies:**
```bash
       # Install Air for hot reload
       go install github.com/cosmtrek/air@latest
       
       # Setup development environment
       ./scripts/setup-dev.sh
```
3. **Start development with hot reload:**

   Terminal 1 - Start the operator:
```bash
       air
```
   Terminal 2 - Test with samples:
```bash
       kubectl apply -f config/samples/runbook_v1alpha1_runbook.yaml
       kubectl get runbooks
```
### Project Structure
```
    runbook-operator/
    ├── api/v1alpha1/           # Custom Resource Definitions
    ├── internal/controller/    # Controller implementations
    ├── pkg/
    │   ├── generator/         # Runbook content generation
    │   ├── outputs/           # Output format handlers
    │   ├── validator/         # Runbook validation
    │   └── templates/         # Template management
    ├── config/
    │   ├── crd/              # CRD manifests
    │   ├── samples/          # Example resources
    │   └── default/          # Default deployment config
    └── scripts/              # Development scripts
```
### Running Tests
```bash
    # Run unit tests
    make test
    
    # Run integration tests
    make test-integration
    
    # Run with coverage
    make test-coverage
```
### Building and Deployment
```bash
    # Build locally
    make build
    
    # Build and push Docker image
    make docker-build docker-push IMG=your-registry/runbook-operator:tag
    
    # Deploy to cluster
    make deploy IMG=your-registry/runbook-operator:tag
```
## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `RECONCILE_INTERVAL` | How often to reconcile resources | `30s` |
| `OUTPUT_BASE_PATH` | Base path for file outputs | `/tmp/runbooks` |
| `ENABLE_VALIDATION` | Enable runbook validation | `true` |

### ConfigMap Configuration
```yaml
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: runbook-operator-config
    data:
      config.yaml: |
        generator:
          defaultTemplate: "standard-runbook"
          outputFormats: ["markdown", "html"]
        validation:
          enabled: true
          checkLinks: true
        outputs:
          markdown:
            enabled: true
            basePath: "/tmp/runbooks"
          html:
            enabled: true
            basePath: "/tmp/runbooks/html"
```
## Roadmap

- [ ] Integration with Confluence and Notion
- [ ] Advanced template engine with conditionals
- [ ] Runbook execution automation
- [ ] Metrics and monitoring dashboard
- [ ] Multi-language support
- [ ] Plugin architecture for custom outputs

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Process

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Ensure all tests pass: `make test`
5. Submit a pull request

### Code Standards

- Follow Go conventions and best practices
- Add tests for new functionality
- Update documentation for API changes
- Use conventional commit messages

## Uninstallation
```bash
    # Delete sample resources
    kubectl delete -k config/samples/
    
    # Remove the operator
    make undeploy
    
    # Remove CRDs
    make uninstall
```
## License

Copyright 2025 Geovane Guibes.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

