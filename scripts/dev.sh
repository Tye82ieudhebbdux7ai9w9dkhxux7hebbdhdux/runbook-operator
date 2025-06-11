#!/bin/bash
set -e

echo "🚀 Starting RunbookOperator development environment..."

# Check if kind cluster exists
if ! kind get clusters | grep -q "runbook-dev"; then
    echo "📦 Creating kind cluster..."
    kind create cluster --name runbook-dev --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF
fi

# Switch to the dev cluster
kubectl config use-context kind-runbook-dev

# Install CRDs
echo "📋 Installing CRDs..."
make install

# Start the operator with hot reload
echo "🔥 Starting hot reload development..."
air
