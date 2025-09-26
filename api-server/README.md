# PostgreSQL AGE Operator API Server

A REST API server for managing PostgreSQL AGE Operator customer instances programmatically. This allows you to create, manage, and delete customer instances through HTTP requests instead of manual kubectl commands.

## 🚀 Features

- **REST API**: Full CRUD operations for customer instances
- **Authentication**: API key-based authentication
- **Kubernetes Native**: Direct integration with Kubernetes API
- **Real-time Status**: Get live status of customer instances
- **Pagination**: Support for paginated listing
- **Rate Limiting**: Built-in rate limiting for API protection
- **Health Checks**: Health and readiness endpoints
- **OpenAPI Documentation**: Built-in API documentation

## 📋 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/docs` | API documentation |
| POST | `/api/v1/customers` | Create customer instance |
| GET | `/api/v1/customers` | List all customer instances |
| GET | `/api/v1/customers/:name` | Get specific customer instance |
| PUT | `/api/v1/customers/:name` | Update customer instance |
| DELETE | `/api/v1/customers/:name` | Delete customer instance |
| GET | `/api/v1/customers/:name/status` | Get customer status |

## 🏗️ Quick Start

### 1. Build and Deploy API Server

```bash
# Build the API server image
cd api-server
docker build -t localhost/postgres-age-api:latest .

# For Kind clusters: Load the image
kind load docker-image localhost/postgres-age-api:latest

# Deploy the API server
kubectl apply -f deployments/kubernetes.yaml
```

### 2. Get API Key

```bash
# Get the API key from the secret
kubectl get secret postgres-age-api-config -n postgres-age-api -o jsonpath='{.data.api-key}' | base64 -d
```

### 3. Access the API

```bash
# Port forward to access locally
kubectl port-forward -n postgres-age-api service/postgres-age-api 8080:80

# Test the API
curl http://localhost:8080/health
```

## 💻 API Usage Examples

### Create a Customer Instance

```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "acme-corp",
    "displayName": "ACME Corporation",
    "imageTag": "latest",
    "config": {
      "resources": {
        "requests": {"cpu": "200m", "memory": "256Mi"},
        "limits": {"cpu": "1000m", "memory": "1Gi"}
      },
      "highAvailability": true,
      "backupEnabled": true,
      "monitoringEnabled": true
    },
    "labels": {
      "environment": "production",
      "team": "platform"
    }
  }'
```

### List Customer Instances

```bash
curl -X GET "http://localhost:8080/api/v1/customers?page=1&pageSize=10" \
  -H "Authorization: Bearer your-api-key"
```

### Get Specific Customer

```bash
curl -X GET http://localhost:8080/api/v1/customers/acme-corp \
  -H "Authorization: Bearer your-api-key"
```

### Update Customer Instance

```bash
curl -X PUT http://localhost:8080/api/v1/customers/acme-corp \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "imageTag": "v1.2.0",
    "config": {
      "resources": {
        "limits": {"cpu": "2000m", "memory": "2Gi"}
      }
    }
  }'
```

### Delete Customer Instance

```bash
curl -X DELETE http://localhost:8080/api/v1/customers/acme-corp \
  -H "Authorization: Bearer your-api-key"
```

### Get Customer Status

```bash
curl -X GET http://localhost:8080/api/v1/customers/acme-corp/status \
  -H "Authorization: Bearer your-api-key"
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port to run the server on | `8080` |
| `API_KEY` | API key for authentication | `""` (disabled) |
| `DEBUG` | Enable debug mode | `false` |

### Command Line Flags

```bash
./api-server -help
  -api-key string
        API key for authentication (optional)
  -debug
        Enable debug mode
  -port string
        Port to run the server on (default "8080")
```

## 🏗️ Development

### Local Development

```bash
# Install dependencies
cd api-server
go mod download

# Run locally (no auth)
go run cmd/main.go -debug

# Run with API key
go run cmd/main.go -debug -api-key="dev-api-key"
```

### Building

```bash
# Build binary
go build -o api-server cmd/main.go

# Build Docker image
docker build -t postgres-age-api:latest .
```

### Testing

```bash
# Run tests
go test ./...

# Test API endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/docs
```

## 🔐 Security

### Authentication

The API uses Bearer token authentication:

```bash
curl -H "Authorization: Bearer your-api-key" ...
```

### Rate Limiting

- 100 requests per minute per IP address
- Configurable in the middleware

### RBAC

The API server requires the following Kubernetes permissions:
- Full access to namespaces, service accounts, RBAC resources
- Full access to deployments and replica sets
- Read access to pods, services, configmaps, secrets
- Read access to PostgreSQL custom resources

## 📊 Monitoring

### Health Checks

```bash
# Health check
curl http://localhost:8080/health

# Response
{
  "status": "healthy",
  "timestamp": "2023-10-01T12:00:00Z",
  "version": "1.0.0"
}
```

### Metrics

The API server includes:
- Request/response logging
- Request ID tracing
- Error tracking
- Performance metrics

### Kubernetes Probes

- **Liveness Probe**: `/health` endpoint
- **Readiness Probe**: `/health` endpoint

## 🚀 Integration Examples

### SaaS Platform Integration

```python
# Python example
import requests

class PostgresAGEManager:
    def __init__(self, api_url, api_key):
        self.api_url = api_url
        self.headers = {
            'Authorization': f'Bearer {api_key}',
            'Content-Type': 'application/json'
        }

    def create_customer(self, customer_name, config=None):
        payload = {
            'name': customer_name,
            'config': config or {}
        }
        response = requests.post(
            f'{self.api_url}/api/v1/customers',
            json=payload,
            headers=self.headers
        )
        return response.json()

    def delete_customer(self, customer_name):
        response = requests.delete(
            f'{self.api_url}/api/v1/customers/{customer_name}',
            headers=self.headers
        )
        return response.json()

# Usage
manager = PostgresAGEManager('http://localhost:8080', 'your-api-key')
result = manager.create_customer('new-customer')
```

### Terraform Integration

```hcl
# terraform/customer.tf
resource "null_resource" "postgres_age_customer" {
  provisioner "local-exec" {
    command = <<EOF
curl -X POST ${var.api_url}/api/v1/customers \
  -H "Authorization: Bearer ${var.api_key}" \
  -H "Content-Type: application/json" \
  -d '${jsonencode({
    name = var.customer_name
    config = var.customer_config
  })}'
EOF
  }
}
```

### CI/CD Pipeline Integration

```yaml
# .github/workflows/deploy-customer.yml
name: Deploy Customer Instance
on:
  workflow_dispatch:
    inputs:
      customer_name:
        description: 'Customer name'
        required: true
      environment:
        description: 'Environment'
        required: true
        default: 'staging'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - name: Create Customer Instance
      run: |
        curl -X POST ${{ vars.API_URL }}/api/v1/customers \
          -H "Authorization: Bearer ${{ secrets.API_KEY }}" \
          -H "Content-Type: application/json" \
          -d '{
            "name": "${{ github.event.inputs.customer_name }}",
            "imageTag": "${{ github.event.inputs.environment == \"production\" && \"latest\" || \"dev-latest\" }}",
            "labels": {
              "environment": "${{ github.event.inputs.environment }}"
            }
          }'
```

## 🎯 Use Cases

1. **Multi-Tenant SaaS**: Automatically provision PostgreSQL AGE instances for new customers
2. **DevOps Automation**: Integrate with CI/CD pipelines for environment provisioning
3. **Self-Service Portals**: Allow teams to provision their own database instances
4. **Infrastructure as Code**: Manage database instances through Terraform/Pulumi
5. **Monitoring Dashboards**: Build custom dashboards showing instance status
6. **Automated Scaling**: Programmatically create instances based on demand

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](../LICENSE.md) file for details.