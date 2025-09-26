# Multi-Instance Installation Guide

This guide explains how to install multiple instances of the PostgreSQL AGE Operator in the same Kubernetes cluster without conflicts.

## Current Limitations

The default installation has several hardcoded names that prevent multiple instances:

### Cluster-Level Conflicts
- ClusterRole: `postgres-operator`
- ClusterRoleBinding: `postgres-operator`
- CRDs: All use `postgres-operator.crunchydata.com` API group

### Namespace-Level Conflicts
- Default namespace: `postgres-operator`
- ServiceAccount: `pgo`
- Deployment: `pgo`
- Control plane label: `postgres-operator.crunchydata.com/control-plane: postgres-operator`

## Solutions

### Option 1: Different Namespaces (Recommended)

Install each instance in a different namespace with unique naming:

```bash
# Instance 1 - Default
kubectl apply --server-side -k config/default

# Instance 2 - AGE Development
kustomize build config/default | sed 's/postgres-operator/postgres-operator-age-dev/g' | kubectl apply -f -

# Instance 3 - AGE Staging
kustomize build config/default | sed 's/postgres-operator/postgres-operator-age-staging/g' | kubectl apply -f -
```

### Option 2: Custom Kustomization Overlays

Create environment-specific overlays:

**Create `config/overlays/age-dev/kustomization.yaml`:**
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- ../../default

namespace: postgres-operator-age-dev

namePrefix: age-dev-
nameSuffix: -v2

labels:
- includeSelectors: true
  pairs:
    postgres-operator.crunchydata.com/control-plane: postgres-operator-age-dev
    app.kubernetes.io/instance: age-dev

images:
- name: postgres-operator
  newName: localhost/postgres-age-operator
  newTag: dev-latest
```

**Deploy:**
```bash
kubectl apply --server-side -k config/overlays/age-dev
```

### Option 3: Helm Chart Approach

Convert to Helm chart with templated values:

**`values.yaml`:**
```yaml
nameOverride: ""
fullnameOverride: ""

namespace: postgres-operator-age

operator:
  name: pgo-age
  image:
    repository: localhost/postgres-age-operator
    tag: latest

controlPlane:
  name: postgres-operator-age

rbac:
  clusterRoleName: postgres-operator-age
  clusterRoleBindingName: postgres-operator-age
```

### Option 4: Different API Groups (Breaking Change)

**⚠️ Warning**: This requires code changes and creates API incompatibility.

1. Change API group from `postgres-operator.crunchydata.com` to `postgres-operator.gregdata.com`
2. Regenerate CRDs
3. Update all Go import paths
4. Update client configurations

## Recommended Multi-Instance Strategy

### For Development/Testing:
```bash
# Production instance
kubectl apply --server-side -k config/default

# Development instance
kubectl create namespace postgres-operator-dev
kustomize build config/default | sed 's/namespace: postgres-operator/namespace: postgres-operator-dev/g' | sed 's/name: postgres-operator/name: postgres-operator-dev/g' | kubectl apply -f -
```

### For Production Multi-Tenancy:
1. Use Option 2 (Kustomization Overlays)
2. Create separate namespaces per environment
3. Use different image tags for each environment
4. Configure resource quotas per namespace

## Migration Path

If you have existing PGO installations:

1. **Backup existing clusters**
2. **Create new namespace** for AGE operator
3. **Deploy AGE operator** with unique naming
4. **Migrate data** from old clusters to new AGE-enabled clusters
5. **Decommission** old clusters after verification

## Monitoring Multiple Instances

Each instance needs separate monitoring configuration:

```bash
# Monitor specific instance
kubectl get pods -n postgres-operator-age-dev -l postgres-operator.crunchydata.com/control-plane=postgres-operator-age-dev

# Check all AGE clusters across instances
kubectl get postgresclusters --all-namespaces -l postgres-operator.crunchydata.com/cluster
```

## Troubleshooting

### Common Issues:
1. **CRD conflicts**: Only one set of CRDs can exist cluster-wide
2. **RBAC conflicts**: ClusterRole names must be unique
3. **Webhook conflicts**: Admission controllers may conflict

### Solutions:
- Use unique naming prefixes
- Deploy in separate namespaces
- Monitor resource conflicts with `kubectl describe`