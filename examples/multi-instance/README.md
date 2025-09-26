# Multi-Instance Examples

This directory contains examples for deploying PostgreSQL AGE clusters with the multi-instance operator setup.

## Quick Start

### 1. Deploy Operator Instance for a Client

```bash
# Deploy for ACME Corporation
./scripts/install-client.sh --client acme-corp

# Deploy for development environment
./scripts/install-client.sh --client dev --image-tag dev-latest

# Deploy for staging environment
./scripts/install-client.sh --client staging --image-tag v1.2.3
```

### 2. Deploy PostgreSQL AGE Cluster

```bash
# Deploy ACME Corp cluster
kubectl apply -f examples/multi-instance/acme-corp-cluster.yaml

# Deploy development cluster
kubectl apply -f examples/multi-instance/dev-cluster.yaml
```

## Example Scenarios

### Production Client Deployment

```bash
# 1. Install operator for client
./scripts/install-client.sh --client acme-corp

# 2. Deploy high-availability cluster
kubectl apply -f examples/multi-instance/acme-corp-cluster.yaml

# 3. Verify deployment
kubectl get pods -n postgres-operator-acme-corp
kubectl get postgresclusters -n postgres-operator-acme-corp
```

### Development Environment

```bash
# 1. Install operator for dev
./scripts/install-client.sh --client dev --image-tag dev-latest

# 2. Deploy single-instance cluster
kubectl apply -f examples/multi-instance/dev-cluster.yaml

# 3. Connect to database
kubectl exec -n postgres-operator-dev -it \
  $(kubectl get pod -n postgres-operator-dev -l postgres-operator.crunchydata.com/role=master -o name) \
  -c database -- psql age_dev_db age_dev_user
```

### Multiple Clients on Same Cluster

```bash
# Deploy multiple clients
./scripts/install-client.sh --client acme-corp
./scripts/install-client.sh --client globex-inc
./scripts/install-client.sh --client dev
./scripts/install-client.sh --client staging

# Each gets their own:
# - Namespace: postgres-operator-{client}
# - RBAC: postgres-operator-{client}
# - Resources: {client}-*-age
```

## Management Commands

### List All Client Instances

```bash
./scripts/install-client.sh --list
```

### Remove Client Instance

```bash
# Remove dev environment
./scripts/install-client.sh --client dev --uninstall

# Remove client completely
./scripts/install-client.sh --client acme-corp --uninstall
```

### Dry Run (Preview Changes)

```bash
./scripts/install-client.sh --client new-client --dry-run
```

## Monitoring Multiple Instances

### Check All Instances

```bash
# List all AGE operator deployments
kubectl get deployments -A -l app.kubernetes.io/instance

# List all PostgreSQL clusters across all clients
kubectl get postgresclusters --all-namespaces

# Monitor specific client
kubectl get pods -n postgres-operator-acme-corp -w
```

### Resource Usage

```bash
# Check resource usage per client
kubectl top pods -n postgres-operator-acme-corp
kubectl top pods -n postgres-operator-dev

# Check storage usage
kubectl get pvc -n postgres-operator-acme-corp
kubectl get pvc -n postgres-operator-dev
```

## Customization

### Client-Specific Configuration

Create custom cluster configurations by copying and modifying the examples:

```bash
cp examples/multi-instance/acme-corp-cluster.yaml examples/multi-instance/my-client-cluster.yaml

# Edit the file to:
# 1. Change namespace to: postgres-operator-my-client
# 2. Update labels with: app.kubernetes.io/client: my-client
# 3. Adjust resources and settings as needed
```

### Environment-Specific Settings

| Environment | Replicas | Resources | Storage | Monitoring |
|-------------|----------|-----------|---------|------------|
| Development | 1 | Small | 5Gi | Basic |
| Staging | 2 | Medium | 20Gi | Full |
| Production | 3+ | Large | 100Gi+ | Full + Alerts |

## Troubleshooting

### Common Issues

1. **Namespace already exists**: Use `--uninstall` first
2. **RBAC conflicts**: Each client gets unique RBAC resources
3. **Image pull issues**: Ensure image exists for specified tag
4. **Resource limits**: Check cluster capacity for multiple instances

### Debug Commands

```bash
# Check operator logs for specific client
kubectl logs -n postgres-operator-acme-corp -l app.kubernetes.io/instance=acme-corp-age

# Check cluster status
kubectl describe postgrescluster -n postgres-operator-acme-corp age-cluster

# Check events
kubectl get events -n postgres-operator-acme-corp --sort-by='.lastTimestamp'
```