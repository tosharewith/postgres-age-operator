# Multi-Client Deployment Guide

The PostgreSQL AGE Operator supports deploying multiple isolated instances for different clients, environments, or teams on the same Kubernetes cluster.

## 🚀 Quick Start

Deploy separate operator instances for different clients:

```bash
# Deploy for ACME Corporation
./scripts/install-client.sh --client acme-corp

# Deploy for development team
./scripts/install-client.sh --client dev --image-tag dev-latest

# Deploy for staging environment
./scripts/install-client.sh --client staging
```

Each client gets completely isolated resources with no conflicts.

## 📋 What Gets Created Per Client

For client `acme-corp`, the script creates:

| Resource Type | Name | Namespace |
|---------------|------|-----------|
| Namespace | `postgres-operator-acme-corp` | - |
| ServiceAccount | `acme-corp-pgo-age` | `postgres-operator-acme-corp` |
| ClusterRole | `postgres-operator-acme-corp` | - |
| ClusterRoleBinding | `postgres-operator-acme-corp` | - |
| Deployment | `acme-corp-pgo-age` | `postgres-operator-acme-corp` |

## 🎯 Use Cases

### Multi-Tenant SaaS Platform
```bash
# Deploy operator per customer
./scripts/install-client.sh --client customer-a
./scripts/install-client.sh --client customer-b
./scripts/install-client.sh --client customer-c

# Each customer gets isolated PostgreSQL AGE clusters
```

### Environment Separation
```bash
# Different environments with different image tags
./scripts/install-client.sh --client dev --image-tag dev-latest
./scripts/install-client.sh --client staging --image-tag v1.2.3
./scripts/install-client.sh --client prod --image-tag v1.2.3
```

### Team-Based Deployment
```bash
# Different teams get their own instances
./scripts/install-client.sh --client team-alpha
./scripts/install-client.sh --client team-beta
./scripts/install-client.sh --client team-gamma
```

## 🔧 Installation Script Options

```bash
./scripts/install-client.sh --help
```

### Common Commands

```bash
# Install new client
./scripts/install-client.sh --client my-client

# Preview what will be created (dry-run)
./scripts/install-client.sh --client my-client --dry-run

# Use specific image tag
./scripts/install-client.sh --client my-client --image-tag v1.5.0

# List all deployed clients
./scripts/install-client.sh --list

# Remove client instance
./scripts/install-client.sh --client my-client --uninstall
```

## 📦 Deploying PostgreSQL Clusters

After installing the operator for a client, deploy PostgreSQL AGE clusters:

```bash
# 1. Install operator for client
./scripts/install-client.sh --client acme-corp

# 2. Deploy cluster in client namespace
kubectl apply -f examples/multi-instance/acme-corp-cluster.yaml

# 3. Verify
kubectl get pods -n postgres-operator-acme-corp
```

## 🔍 Managing Multiple Instances

### List All Client Resources

```bash
# List all client namespaces
kubectl get namespaces -l app.kubernetes.io/client

# List all client deployments
kubectl get deployments -A -l app.kubernetes.io/client

# List all PostgreSQL clusters across clients
kubectl get postgresclusters --all-namespaces
```

### Monitor Specific Client

```bash
# Monitor ACME Corp resources
kubectl get pods -n postgres-operator-acme-corp -w
kubectl logs -n postgres-operator-acme-corp -l app.kubernetes.io/instance=acme-corp-age

# Check cluster status
kubectl describe postgrescluster -n postgres-operator-acme-corp
```

### Resource Usage Per Client

```bash
# CPU/Memory usage
kubectl top pods -n postgres-operator-acme-corp

# Storage usage
kubectl get pvc -n postgres-operator-acme-corp

# Events
kubectl get events -n postgres-operator-acme-corp --sort-by='.lastTimestamp'
```

## 🏗️ Architecture

### Isolation Levels

1. **Namespace Isolation**: Each client operates in its own namespace
2. **RBAC Isolation**: Unique ClusterRole/ClusterRoleBinding per client
3. **Resource Isolation**: All resources prefixed with client name
4. **Network Policies**: Optional network isolation between clients

### Resource Naming Pattern

```
Client: acme-corp
├── Namespace: postgres-operator-acme-corp
├── ServiceAccount: acme-corp-pgo-age
├── ClusterRole: postgres-operator-acme-corp
├── Deployment: acme-corp-pgo-age
└── Labels:
    ├── app.kubernetes.io/client: acme-corp
    ├── app.kubernetes.io/instance: acme-corp-age
    └── postgres-operator.crunchydata.com/control-plane: postgres-operator-acme-corp
```

## ⚠️ Limitations & Considerations

### Cluster-Wide Resources
- **CRDs**: Shared across all clients (same API definitions)
- **ClusterRoles**: Unique per client but all have similar permissions
- **Admission Controllers**: May affect all clients

### Resource Limits
- Each client instance consumes cluster resources
- Consider setting resource quotas per client namespace
- Monitor total resource usage across all clients

### Networking
- All clients share the same cluster network by default
- Consider NetworkPolicies for strict isolation
- Service names are namespace-scoped (naturally isolated)

## 🔐 Security Best Practices

### Namespace Security

```yaml
# Apply resource quotas per client
apiVersion: v1
kind: ResourceQuota
metadata:
  name: acme-corp-quota
  namespace: postgres-operator-acme-corp
spec:
  hard:
    requests.cpu: "4"
    requests.memory: 8Gi
    limits.cpu: "8"
    limits.memory: 16Gi
    persistentvolumeclaims: "10"
```

### Network Policies

```yaml
# Isolate client network traffic
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: acme-corp-isolation
  namespace: postgres-operator-acme-corp
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          app.kubernetes.io/client: acme-corp
```

## 🔄 Migration & Upgrades

### Upgrading Client Instances

```bash
# Upgrade specific client to new image
./scripts/install-client.sh --client acme-corp --image-tag v1.6.0

# The script will update the existing deployment
```

### Migrating Between Clients

```bash
# 1. Backup data from source client
kubectl exec -n postgres-operator-old-client ...

# 2. Deploy new client instance
./scripts/install-client.sh --client new-client

# 3. Restore data to new client
kubectl exec -n postgres-operator-new-client ...

# 4. Remove old client
./scripts/install-client.sh --client old-client --uninstall
```

## 🔧 Troubleshooting

### Common Issues

1. **Client name conflicts**: Use unique, DNS-compliant names
2. **Resource exhaustion**: Monitor cluster capacity
3. **RBAC conflicts**: Each client gets unique RBAC resources
4. **Image pull failures**: Ensure image exists with specified tag

### Debug Commands

```bash
# Check client installation status
./scripts/install-client.sh --list

# Check operator logs
kubectl logs -n postgres-operator-acme-corp deployment/acme-corp-pgo-age

# Check cluster events
kubectl get events -n postgres-operator-acme-corp

# Validate client configuration
kubectl get clusterrole postgres-operator-acme-corp
kubectl get clusterrolebinding postgres-operator-acme-corp
```

## 📈 Scaling Considerations

### Horizontal Scaling
- Deploy clients across multiple Kubernetes clusters
- Use external DNS for cross-cluster communication
- Implement cluster federation for management

### Vertical Scaling
- Adjust resource quotas per client namespace
- Monitor and tune operator resource requests/limits
- Scale PostgreSQL clusters independently per client

## 🤝 Contributing

To improve the multi-client deployment system:

1. Test with different client naming patterns
2. Add support for additional configuration options
3. Improve resource management and monitoring
4. Add integration with service mesh for advanced networking

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.