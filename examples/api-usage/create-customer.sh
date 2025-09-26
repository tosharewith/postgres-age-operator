#!/bin/bash

# PostgreSQL AGE Operator API - Create Customer Example
# This script demonstrates how to create a new customer instance via API

set -euo pipefail

# Configuration
API_URL="${API_URL:-http://localhost:8080}"
API_KEY="${API_KEY:-}"

# Customer configuration
CUSTOMER_NAME="${1:-demo-customer}"
DISPLAY_NAME="${2:-Demo Customer}"
IMAGE_TAG="${3:-latest}"
ENVIRONMENT="${4:-development}"

# Help function
show_help() {
    cat << EOF
Create Customer Instance via API

Usage: $0 [CUSTOMER_NAME] [DISPLAY_NAME] [IMAGE_TAG] [ENVIRONMENT]

Arguments:
  CUSTOMER_NAME    Name of the customer (default: demo-customer)
  DISPLAY_NAME     Display name (default: Demo Customer)
  IMAGE_TAG        Docker image tag (default: latest)
  ENVIRONMENT      Environment type (default: development)

Environment Variables:
  API_URL          API server URL (default: http://localhost:8080)
  API_KEY          API authentication key (required)

Examples:
  # Create development customer
  API_KEY="your-key" $0 dev-client "Development Client" dev-latest development

  # Create production customer
  API_KEY="your-key" $0 acme-corp "ACME Corporation" v1.2.3 production

EOF
}

# Check for help flag
if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
    show_help
    exit 0
fi

# Validate API key
if [[ -z "$API_KEY" ]]; then
    echo "Error: API_KEY environment variable is required"
    echo "Set it with: export API_KEY='your-api-key'"
    exit 1
fi

# Validate customer name
if [[ ! "$CUSTOMER_NAME" =~ ^[a-z0-9]([-a-z0-9]*[a-z0-9])?$ ]]; then
    echo "Error: Customer name '$CUSTOMER_NAME' is not valid"
    echo "Must be lowercase alphanumeric with optional hyphens"
    exit 1
fi

if [[ ${#CUSTOMER_NAME} -gt 20 ]]; then
    echo "Error: Customer name '$CUSTOMER_NAME' is too long (max 20 characters)"
    exit 1
fi

# Set configuration based on environment
if [[ "$ENVIRONMENT" == "production" ]]; then
    CPU_REQUEST="500m"
    CPU_LIMIT="2000m"
    MEMORY_REQUEST="1Gi"
    MEMORY_LIMIT="4Gi"
    HIGH_AVAILABILITY="true"
    BACKUP_ENABLED="true"
    MONITORING_ENABLED="true"
elif [[ "$ENVIRONMENT" == "staging" ]]; then
    CPU_REQUEST="250m"
    CPU_LIMIT="1000m"
    MEMORY_REQUEST="512Mi"
    MEMORY_LIMIT="2Gi"
    HIGH_AVAILABILITY="true"
    BACKUP_ENABLED="true"
    MONITORING_ENABLED="true"
else  # development
    CPU_REQUEST="100m"
    CPU_LIMIT="500m"
    MEMORY_REQUEST="256Mi"
    MEMORY_LIMIT="1Gi"
    HIGH_AVAILABILITY="false"
    BACKUP_ENABLED="false"
    MONITORING_ENABLED="false"
fi

# Create JSON payload
PAYLOAD=$(cat <<EOF
{
  "name": "$CUSTOMER_NAME",
  "displayName": "$DISPLAY_NAME",
  "imageTag": "$IMAGE_TAG",
  "config": {
    "resources": {
      "requests": {
        "cpu": "$CPU_REQUEST",
        "memory": "$MEMORY_REQUEST"
      },
      "limits": {
        "cpu": "$CPU_LIMIT",
        "memory": "$MEMORY_LIMIT"
      }
    },
    "highAvailability": $HIGH_AVAILABILITY,
    "backupEnabled": $BACKUP_ENABLED,
    "monitoringEnabled": $MONITORING_ENABLED
  },
  "labels": {
    "environment": "$ENVIRONMENT",
    "managed-by": "api-script",
    "created-at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  }
}
EOF
)

echo "🚀 Creating customer instance: $CUSTOMER_NAME"
echo "   Display Name: $DISPLAY_NAME"
echo "   Environment: $ENVIRONMENT"
echo "   Image Tag: $IMAGE_TAG"
echo "   API URL: $API_URL"
echo

# Make API request
echo "📡 Making API request..."
RESPONSE=$(curl -s -w "\n%{http_code}" \
  -X POST "$API_URL/api/v1/customers" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d "$PAYLOAD")

# Parse response
HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
BODY=$(echo "$RESPONSE" | head -n -1)

# Handle response
if [[ "$HTTP_CODE" == "201" ]]; then
    echo "✅ Customer instance created successfully!"
    echo
    echo "📋 Customer Details:"
    echo "$BODY" | jq -r '
        .data |
        "   ID: \(.id)",
        "   Name: \(.name)",
        "   Namespace: \(.namespace)",
        "   Status: \(.status.phase)",
        "   Created: \(.createdAt)"
    '
    echo
    echo "🔍 Next steps:"
    echo "   1. Check status: curl -H \"Authorization: Bearer \$API_KEY\" $API_URL/api/v1/customers/$CUSTOMER_NAME/status"
    echo "   2. Monitor deployment: kubectl get pods -n postgres-operator-$CUSTOMER_NAME -w"
    echo "   3. Create PostgreSQL cluster in the customer namespace"
    echo
elif [[ "$HTTP_CODE" == "400" ]]; then
    echo "❌ Bad Request - Validation failed"
    echo "$BODY" | jq -r '.error.message'
    if echo "$BODY" | jq -e '.error.details' > /dev/null; then
        echo
        echo "Validation errors:"
        echo "$BODY" | jq -r '.error.details[] | "   \(.field): \(.message)"'
    fi
    exit 1
elif [[ "$HTTP_CODE" == "401" ]]; then
    echo "❌ Unauthorized - Check your API key"
    exit 1
elif [[ "$HTTP_CODE" == "409" ]]; then
    echo "❌ Conflict - Customer instance already exists"
    echo "$BODY" | jq -r '.error.message'
    exit 1
elif [[ "$HTTP_CODE" == "500" ]]; then
    echo "❌ Internal Server Error"
    echo "$BODY" | jq -r '.error.message'
    exit 1
else
    echo "❌ Unexpected response (HTTP $HTTP_CODE)"
    echo "$BODY"
    exit 1
fi