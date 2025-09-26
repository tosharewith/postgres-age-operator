#!/bin/bash

# PostgreSQL AGE Operator - Multi-Instance Installation Script
# This script allows easy deployment of the operator per client/environment

set -euo pipefail

# Default values
CLIENT_NAME=""
IMAGE_TAG="latest"
DRY_RUN=false
UNINSTALL=false
LIST_CLIENTS=false

# Help function
show_help() {
    cat << EOF
PostgreSQL AGE Operator Multi-Instance Installer

Usage: $0 --client CLIENT_NAME [OPTIONS]

REQUIRED:
  --client CLIENT_NAME      Name of the client/instance (e.g., acme-corp, dev, staging)

OPTIONS:
  --image-tag TAG          Docker image tag (default: latest)
  --dry-run               Show what would be deployed without applying
  --uninstall             Remove the client instance
  --list                  List all deployed client instances
  --help                  Show this help message

EXAMPLES:
  # Deploy for client "acme-corp"
  $0 --client acme-corp

  # Deploy development instance
  $0 --client dev --image-tag dev-latest

  # Deploy staging with dry-run
  $0 --client staging --image-tag v1.2.3 --dry-run

  # Remove client instance
  $0 --client acme-corp --uninstall

  # List all deployed instances
  $0 --list

WHAT THIS CREATES:
  - Namespace: postgres-operator-CLIENT_NAME
  - ServiceAccount: CLIENT_NAME-pgo-age
  - ClusterRole: postgres-operator-CLIENT_NAME
  - ClusterRoleBinding: postgres-operator-CLIENT_NAME
  - Deployment: CLIENT_NAME-pgo-age
  - All resources labeled with client name for easy management

EOF
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --client)
            CLIENT_NAME="$2"
            shift 2
            ;;
        --image-tag)
            IMAGE_TAG="$2"
            shift 2
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --uninstall)
            UNINSTALL=true
            shift
            ;;
        --list)
            LIST_CLIENTS=true
            shift
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Function to list existing client instances
list_client_instances() {
    echo "=== Deployed PostgreSQL AGE Operator Instances ==="
    echo

    # List namespaces
    echo "Client Namespaces:"
    kubectl get namespaces -l app.kubernetes.io/client --no-headers 2>/dev/null || echo "No client instances found"
    echo

    # List deployments
    echo "Client Deployments:"
    kubectl get deployments -A -l app.kubernetes.io/client --no-headers 2>/dev/null || echo "No client deployments found"
    echo

    # List cluster roles
    echo "Client ClusterRoles:"
    kubectl get clusterroles | grep "postgres-operator-" || echo "No client cluster roles found"
    echo
}

# Function to validate client name
validate_client_name() {
    if [[ -z "$CLIENT_NAME" ]]; then
        echo "Error: --client CLIENT_NAME is required"
        show_help
        exit 1
    fi

    # Check if client name is valid for Kubernetes
    if [[ ! "$CLIENT_NAME" =~ ^[a-z0-9]([-a-z0-9]*[a-z0-9])?$ ]]; then
        echo "Error: Client name '$CLIENT_NAME' is not valid for Kubernetes resources"
        echo "Must be lowercase alphanumeric with optional hyphens (RFC 1123)"
        exit 1
    fi

    # Check length
    if [[ ${#CLIENT_NAME} -gt 20 ]]; then
        echo "Error: Client name '$CLIENT_NAME' is too long (max 20 characters)"
        echo "This ensures generated resource names stay within Kubernetes limits"
        exit 1
    fi
}

# Function to check if instance already exists
check_existing_instance() {
    local namespace="postgres-operator-$CLIENT_NAME"
    if kubectl get namespace "$namespace" >/dev/null 2>&1; then
        return 0  # exists
    else
        return 1  # doesn't exist
    fi
}

# Function to uninstall client instance
uninstall_client() {
    local namespace="postgres-operator-$CLIENT_NAME"

    if ! check_existing_instance; then
        echo "Client instance '$CLIENT_NAME' not found"
        exit 1
    fi

    echo "🗑️  Uninstalling PostgreSQL AGE Operator for client: $CLIENT_NAME"
    echo

    # Delete namespace (this will delete most resources)
    echo "Deleting namespace: $namespace"
    kubectl delete namespace "$namespace" --ignore-not-found=true

    # Delete cluster-level resources
    echo "Deleting ClusterRole: postgres-operator-$CLIENT_NAME"
    kubectl delete clusterrole "postgres-operator-$CLIENT_NAME" --ignore-not-found=true

    echo "Deleting ClusterRoleBinding: postgres-operator-$CLIENT_NAME"
    kubectl delete clusterrolebinding "postgres-operator-$CLIENT_NAME" --ignore-not-found=true

    echo "✅ Client instance '$CLIENT_NAME' uninstalled successfully"
}

# Function to install client instance
install_client() {
    local namespace="postgres-operator-$CLIENT_NAME"
    local temp_dir=$(mktemp -d)

    if check_existing_instance; then
        echo "⚠️  Client instance '$CLIENT_NAME' already exists"
        echo "Use --uninstall to remove it first, or choose a different client name"
        exit 1
    fi

    echo "🚀 Installing PostgreSQL AGE Operator for client: $CLIENT_NAME"
    echo "   Namespace: $namespace"
    echo "   Image Tag: $IMAGE_TAG"
    echo

    # Copy template files and replace variables
    cp -r config/multi-instance/* "$temp_dir/"

    # Replace variables in all files
    find "$temp_dir" -type f -name "*.yaml" -exec sed -i "s/CLIENT_NAME/$CLIENT_NAME/g" {} \;
    find "$temp_dir" -type f -name "*.yaml" -exec sed -i "s/IMAGE_TAG/$IMAGE_TAG/g" {} \;

    # Create namespace first
    echo "Creating namespace: $namespace"
    kubectl create namespace "$namespace" --dry-run=client -o yaml | kubectl apply -f -

    # Apply the configuration
    if [[ "$DRY_RUN" == "true" ]]; then
        echo "=== DRY RUN - Would apply the following resources ==="
        $KUSTOMIZE_CMD build "$temp_dir"
    else
        echo "Applying configuration..."
        $KUSTOMIZE_CMD build "$temp_dir" | kubectl apply --server-side -f -

        echo
        echo "✅ PostgreSQL AGE Operator installed successfully for client: $CLIENT_NAME"
        echo
        echo "🔍 Verify installation:"
        echo "   kubectl get pods -n $namespace"
        echo "   kubectl get postgresclusters -n $namespace"
        echo
        echo "📚 Next steps:"
        echo "   1. Create an AGE cluster: kubectl apply -f examples/age-cluster/"
        echo "   2. Modify the namespace in the example to: $namespace"
        echo "   3. Monitor deployment: kubectl logs -n $namespace -l app.kubernetes.io/instance=$CLIENT_NAME-age"
    fi

    # Cleanup
    rm -rf "$temp_dir"
}

# Main logic
main() {
    # Check if kubectl is available
    if ! command -v kubectl >/dev/null 2>&1; then
        echo "Error: kubectl is not installed or not in PATH"
        exit 1
    fi

    # Check if kustomize is available (try kubectl kustomize first)
    if command -v kustomize >/dev/null 2>&1; then
        KUSTOMIZE_CMD="kustomize"
    elif kubectl kustomize --help >/dev/null 2>&1; then
        KUSTOMIZE_CMD="kubectl kustomize"
    else
        echo "Error: Neither 'kustomize' nor 'kubectl kustomize' is available"
        echo "Install kustomize from: https://kubectl.docs.kubernetes.io/installation/kustomize/"
        echo "Or use kubectl 1.14+ which includes kustomize"
        exit 1
    fi

    # Handle list command
    if [[ "$LIST_CLIENTS" == "true" ]]; then
        list_client_instances
        exit 0
    fi

    # Validate required parameters
    validate_client_name

    # Handle uninstall
    if [[ "$UNINSTALL" == "true" ]]; then
        uninstall_client
        exit 0
    fi

    # Install client instance
    install_client
}

# Run main function
main