# PostgreSQL AGE Operator Integration Guide

This document provides a complete guide for installing, configuring, and managing PostgreSQL clusters with Apache AGE (A Graph Extension) using the modified Crunchy PostgreSQL Operator.

## Table of Contents
- [Quick Start](#quick-start)
- [Installation Options](#installation-options)
  - [Kind (Kubernetes in Docker)](#kind-kubernetes-in-docker)
  - [Standard Kubernetes](#standard-kubernetes)
  - [Minikube](#minikube)
- [Operator Management](#operator-management)
- [Deployment Configurations](#deployment-configurations)
- [Advanced Operations](#advanced-operations)

## Quick Start

Get a PostgreSQL + AGE cluster running in under 5 minutes:

```bash
# Clone the repository
git clone https://github.com/gregoriomomm/postgres-operator postgres-age-operator
cd postgres-age-operator

# Build the custom image
docker build -f Dockerfile.age -t localhost/postgres-age-patroni .

# For Kind: Load the image
kind load docker-image localhost/postgres-age-patroni

# Deploy operator
kubectl apply --server-side -k config/default

# Deploy AGE cluster
kubectl apply -k examples/age-cluster/

# Check status
kubectl get pods -n postgres-operator
```

## Overview

Apache AGE is an extension for PostgreSQL that provides graph database functionality. This fork of the Crunchy PostgreSQL Operator has been modified to support deploying and managing PostgreSQL clusters with AGE enabled.

## Key Modifications

### 1. Security Context Changes

**File Modified**: `internal/initialize/security.go`

The operator's default security context has been modified to allow containers to run as root, which is required by the Apache AGE Docker image:

```go
// Line 42 in RestrictedSecurityContext()
RunAsNonRoot: Bool(false),  // Changed from Bool(true)
```

**Impact**: This change relaxes the Kubernetes Pod Security Standards from "Restricted" to "Baseline" profile. Production deployments should implement additional security measures to compensate for this relaxation.

### 2. Custom Docker Image

A custom Docker image (`Dockerfile.age`) has been created that combines Apache AGE with the required operator components:

```dockerfile
FROM apache/age

# Install required components for operator compatibility
RUN apt-get update && \
    apt-get install -y python3 python3-pip python3-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install Patroni for HA orchestration
RUN pip3 install --break-system-packages --no-cache-dir \
    patroni[kubernetes]==3.3.4 \
    psycopg2-binary \
    python-etcd \
    boto3

# Install PostgreSQL tools and extensions
RUN apt-get update && \
    apt-get install -y \
    postgresql-client-16 \
    postgresql-16-pgaudit \
    pgbackrest \
    curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Fix libnss_wrapper path for Debian compatibility
RUN ln -s /usr/lib/x86_64-linux-gnu/libnss_wrapper.so /usr/lib64/libnss_wrapper.so
```

### 3. Operator Configuration Updates

#### Environment Variables
**File**: `config/manager/manager.yaml`

The `RELATED_IMAGE_POSTGRES_16` environment variable has been updated to use the custom AGE image:

```yaml
- name: RELATED_IMAGE_POSTGRES_16
  value: "localhost/postgres-age-patroni"
```

#### Kustomization Configuration
**File**: `config/default/kustomization.yaml`

The operator image has been updated to use the local custom build:

```yaml
images:
- name: postgres-operator
  newName: localhost/postgres-age-operator
  newTag: latest
```

### 4. AGE Cluster Configuration

**Directory**: `examples/age-cluster/`

A new example configuration has been created for deploying AGE-enabled clusters:

#### PostgresCluster Resource
**File**: `examples/age-cluster/postgres-age-cluster.yaml`

```yaml
apiVersion: postgres-operator.crunchydata.com/v1beta1
kind: PostgresCluster
metadata:
  name: age-cluster
spec:
  postgresVersion: 16
  image: localhost/postgres-age-patroni
  imagePullPolicy: Never
  
  instances:
    - name: instance1
      replicas: 3
      dataVolumeClaimSpec:
        accessModes:
        - "ReadWriteOnce"
        resources:
          requests:
            storage: 1Gi
            
  backups:
    pgbackrest:
      repos:
      - name: repo1
        volume:
          volumeClaimSpec:
            accessModes:
            - "ReadWriteOnce"
            resources:
              requests:
                storage: 1Gi
                
  proxy:
    pgBouncer:
      replicas: 1
```

#### AGE Initialization
**File**: `config/age-init.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: age-init-sql
data:
  init.sql: |
    -- Load Apache AGE extension
    CREATE EXTENSION IF NOT EXISTS age;
    LOAD 'age';
    -- Set search path to include ag_catalog
    ALTER DATABASE postgres SET search_path = ag_catalog, "$user", public;
```

## Installation Options

### Prerequisites

- Docker (for building images)
- kubectl (Kubernetes CLI)
- A Kubernetes cluster or Kind/Minikube for local development

### Kind (Kubernetes in Docker)

Kind is perfect for local development and testing.

#### Installing Kind

```bash
# macOS
brew install kind

# Linux
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# Windows (using Chocolatey)
choco install kind
```

#### Create Kind Cluster

```bash
# Create a single-node cluster
kind create cluster --name age-cluster

# Or create a multi-node cluster for HA testing
cat <<EOF | kind create cluster --name age-cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
EOF

# Verify cluster
kubectl cluster-info --context kind-age-cluster
```

#### Deploy on Kind

```bash
# Build the custom AGE image
docker build -f Dockerfile.age -t localhost/postgres-age-patroni .

# Load image into Kind
kind load docker-image localhost/postgres-age-patroni --name age-cluster

# Deploy the operator
kubectl apply --server-side -k config/default

# Wait for operator to be ready
kubectl wait --for=condition=Ready pod -l postgres-operator.crunchydata.com/control-plane=postgres-operator -n postgres-operator --timeout=120s

# Deploy an AGE cluster
kubectl apply -k examples/age-cluster/
```

### Standard Kubernetes

For production Kubernetes clusters (EKS, GKE, AKS, or on-premises).

#### Push Image to Registry

```bash
# Build and tag for your registry
docker build -f Dockerfile.age -t your-registry.com/postgres-age-patroni:latest .

# Push to registry
docker push your-registry.com/postgres-age-patroni:latest

# Update image references in config/manager/manager.yaml
sed -i 's|localhost/postgres-age-patroni|your-registry.com/postgres-age-patroni:latest|g' config/manager/manager.yaml

# Update cluster configuration
sed -i 's|localhost/postgres-age-patroni|your-registry.com/postgres-age-patroni:latest|g' examples/age-cluster/postgres-age-cluster.yaml
sed -i 's|imagePullPolicy: Never|imagePullPolicy: Always|g' examples/age-cluster/postgres-age-cluster.yaml
```

#### Deploy on Kubernetes

```bash
# Deploy the operator
kubectl apply --server-side -k config/default

# Verify operator deployment
kubectl get deployment -n postgres-operator pgo

# Deploy AGE cluster
kubectl apply -k examples/age-cluster/
```

### Minikube

For local development with Minikube.

```bash
# Start Minikube with sufficient resources
minikube start --cpus=4 --memory=8192 --disk-size=20g

# Use Minikube's Docker daemon
eval $(minikube docker-env)

# Build image directly in Minikube
docker build -f Dockerfile.age -t localhost/postgres-age-patroni .

# Deploy operator and cluster
kubectl apply --server-side -k config/default
kubectl apply -k examples/age-cluster/
```

## Operator Management

### Understanding the Operator

The PostgreSQL AGE Operator manages the lifecycle of PostgreSQL clusters with AGE extension. It handles:

- **Cluster Creation**: Automated provisioning of PostgreSQL instances
- **High Availability**: Automatic failover using Patroni
- **Backup Management**: Scheduled and on-demand backups via pgBackRest
- **Scaling**: Dynamic scaling of replicas
- **Updates**: Rolling updates with minimal downtime
- **Monitoring**: Integration with monitoring tools

### Operator Commands

#### Check Operator Status

```bash
# View operator deployment
kubectl get deployment -n postgres-operator pgo

# Check operator logs
kubectl logs -n postgres-operator deployment/pgo

# View operator configuration
kubectl describe deployment -n postgres-operator pgo
```

#### Manage Operator Lifecycle

```bash
# Update operator
kubectl apply --server-side -k config/default

# Restart operator (forces reconciliation)
kubectl rollout restart deployment/pgo -n postgres-operator

# Scale operator replicas (for HA)
kubectl scale deployment/pgo -n postgres-operator --replicas=2

# Delete operator (keeps clusters running)
kubectl delete -k config/default
```

### Operator Configuration

#### Environment Variables

Key environment variables for customizing operator behavior:

```yaml
# Edit config/manager/manager.yaml
env:
  - name: PGO_TARGET_NAMESPACE  # Namespace to watch (empty = all namespaces)
    value: ""
  - name: CRUNCHY_DEBUG         # Enable debug logging
    value: "false"
  - name: RELATED_IMAGE_POSTGRES_16
    value: "localhost/postgres-age-patroni"
  - name: CHECK_FOR_UPGRADES    # Disable upgrade checks
    value: "false"
```

#### Resource Limits

Configure operator resource consumption:

```yaml
# In config/manager/manager.yaml
resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

## Deployment Configurations

### Single Instance (No Replicas)

Perfect for development or non-critical workloads.

#### Create Single Instance Deployment

```yaml
# Save as single-instance.yaml
apiVersion: postgres-operator.crunchydata.com/v1beta1
kind: PostgresCluster
metadata:
  name: age-single
spec:
  postgresVersion: 16
  image: localhost/postgres-age-patroni
  imagePullPolicy: Never
  
  instances:
    - name: instance1
      replicas: 1  # Single instance
      dataVolumeClaimSpec:
        accessModes:
        - "ReadWriteOnce"
        resources:
          requests:
            storage: 10Gi
            
  # Minimal backup configuration
  backups:
    pgbackrest:
      repos:
      - name: repo1
        volume:
          volumeClaimSpec:
            accessModes:
            - "ReadWriteOnce"
            resources:
              requests:
                storage: 10Gi
```

Deploy:
```bash
kubectl apply -f single-instance.yaml
```

### High Availability with Replicas

For production workloads requiring high availability.

#### Create HA Deployment

```yaml
# Save as ha-cluster.yaml
apiVersion: postgres-operator.crunchydata.com/v1beta1
kind: PostgresCluster
metadata:
  name: age-ha
spec:
  postgresVersion: 16
  image: localhost/postgres-age-patroni
  imagePullPolicy: Never
  
  instances:
    - name: instance1
      replicas: 3  # 1 primary + 2 replicas
      dataVolumeClaimSpec:
        accessModes:
        - "ReadWriteOnce"
        resources:
          requests:
            storage: 20Gi
      
      # Resource allocation
      resources:
        limits:
          cpu: 2
          memory: 4Gi
        requests:
          cpu: 500m
          memory: 1Gi
      
      # Pod anti-affinity for spreading across nodes
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 1
            podAffinityTerm:
              topologyKey: kubernetes.io/hostname
              labelSelector:
                matchLabels:
                  postgres-operator.crunchydata.com/cluster: age-ha
                  postgres-operator.crunchydata.com/instance-set: instance1
  
  # Connection pooling
  proxy:
    pgBouncer:
      replicas: 2
      resources:
        limits:
          cpu: 200m
          memory: 256Mi
        requests:
          cpu: 10m
          memory: 64Mi
          
  # Backup configuration with retention
  backups:
    pgbackrest:
      configuration:
      - secret:
          name: age-ha-pgbackrest-secrets
      global:
        repo1-retention-full: "14"
        repo1-retention-full-type: "time"
      repos:
      - name: repo1
        schedules:
          full: "0 1 * * 0"    # Weekly full backup
          incremental: "0 1 * * 1-6"  # Daily incremental
        volume:
          volumeClaimSpec:
            accessModes:
            - "ReadWriteOnce"
            resources:
              requests:
                storage: 50Gi
```

Deploy:
```bash
kubectl apply -f ha-cluster.yaml
```

### Multi-Region/Multi-Zone Deployment

For geographic distribution and disaster recovery.

```yaml
# Save as multi-zone.yaml
apiVersion: postgres-operator.crunchydata.com/v1beta1
kind: PostgresCluster
metadata:
  name: age-multi-zone
spec:
  postgresVersion: 16
  image: localhost/postgres-age-patroni
  
  instances:
    # Zone A instances
    - name: zone-a
      replicas: 2
      dataVolumeClaimSpec:
        accessModes:
        - "ReadWriteOnce"
        resources:
          requests:
            storage: 20Gi
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: topology.kubernetes.io/zone
                operator: In
                values:
                - zone-a
    
    # Zone B instances
    - name: zone-b
      replicas: 2
      dataVolumeClaimSpec:
        accessModes:
        - "ReadWriteOnce"
        resources:
          requests:
            storage: 20Gi
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: topology.kubernetes.io/zone
                operator: In
                values:
                - zone-b
  
  # ... rest of configuration
```

### Custom AGE Configuration

#### With Initialization Script

```yaml
# Save as age-with-init.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: age-init-script
data:
  init.sql: |
    -- Create AGE extension
    CREATE EXTENSION IF NOT EXISTS age;
    LOAD 'age';
    
    -- Create application database
    CREATE DATABASE graph_app;
    \c graph_app
    CREATE EXTENSION IF NOT EXISTS age;
    LOAD 'age';
    SET search_path = ag_catalog, "$user", public;
    
    -- Create initial graph
    SELECT create_graph('app_graph');
    
    -- Create initial schema
    SELECT * FROM cypher('app_graph', $$
      CREATE (n:Config {key: 'version', value: '1.0'})
    $$) as (n agtype);
---
apiVersion: postgres-operator.crunchydata.com/v1beta1
kind: PostgresCluster
metadata:
  name: age-custom
spec:
  postgresVersion: 16
  image: localhost/postgres-age-patroni
  
  databaseInitSQL:
    name: age-init-script
    key: init.sql
    
  instances:
    - name: instance1
      replicas: 2
      dataVolumeClaimSpec:
        accessModes:
        - "ReadWriteOnce"
        resources:
          requests:
            storage: 10Gi
```

### Managing Multiple Clusters

Deploy multiple AGE clusters with different configurations:

```bash
# Development cluster
kubectl apply -f deployments/dev-cluster.yaml

# Staging cluster with replicas
kubectl apply -f deployments/staging-cluster.yaml

# Production cluster with full HA
kubectl apply -f deployments/prod-cluster.yaml

# List all clusters
kubectl get postgrescluster -n postgres-operator

# Get detailed status
kubectl describe postgrescluster -n postgres-operator
```

### Dynamic Scaling

Scale clusters without downtime:

```bash
# Scale up replicas
kubectl patch postgrescluster age-cluster -n postgres-operator \
  --type='json' \
  -p='[{"op": "replace", "path": "/spec/instances/0/replicas", "value":5}]'

# Scale down replicas
kubectl patch postgrescluster age-cluster -n postgres-operator \
  --type='json' \
  -p='[{"op": "replace", "path": "/spec/instances/0/replicas", "value":2}]'

# Add pgBouncer for connection pooling
kubectl patch postgrescluster age-cluster -n postgres-operator \
  --type='merge' \
  -p='{"spec":{"proxy":{"pgBouncer":{"replicas":2}}}}'
```

## Using AGE in Your Cluster

### Connect to the Primary Database

```bash
# Get the primary pod
PRIMARY_POD=$(kubectl get pods -n postgres-operator -l postgres-operator.crunchydata.com/cluster=age-cluster,postgres-operator.crunchydata.com/role=master -o jsonpath='{.items[0].metadata.name}')

# Connect to PostgreSQL
kubectl exec -n postgres-operator $PRIMARY_POD -c database -- psql
```

### Create and Use a Graph

```sql
-- The AGE extension should already be loaded
-- Create a graph
SELECT * FROM ag_catalog.create_graph('my_graph');

-- Set search path
SET search_path = ag_catalog, "$user", public;

-- Create vertices
SELECT * FROM cypher('my_graph', $$
    CREATE (n:Person {name: 'John', age: 30})
    RETURN n
$$) as (n agtype);

-- Create edges
SELECT * FROM cypher('my_graph', $$
    MATCH (a:Person {name: 'John'})
    CREATE (b:Person {name: 'Jane', age: 28})
    CREATE (a)-[:KNOWS]->(b)
    RETURN a, b
$$) as (a agtype, b agtype);

-- Query the graph
SELECT * FROM cypher('my_graph', $$
    MATCH (n:Person)-[:KNOWS]->(m:Person)
    RETURN n.name, m.name
$$) as (person1 agtype, person2 agtype);
```

## High Availability and Replication

The modified operator maintains full HA capabilities with AGE:

- **Patroni Integration**: Manages automatic failover and replica synchronization
- **AGE Replication**: Graph data is automatically replicated to standby nodes
- **pgBackRest**: Provides backup and restore capabilities for AGE data

### Verify Replication Status

```bash
# Check Patroni cluster status
kubectl exec -n postgres-operator $PRIMARY_POD -c database -- patronictl list

# Example output:
# + Cluster: age-cluster-ha (7545681931211747459) -------+---------+-----------+----+-----------+
# | Member                       | Host                 | Role    | State     | TL | Lag in MB |
# +------------------------------+----------------------+---------+-----------+----+-----------+
# | age-cluster-instance1-66rx-0 | age-cluster-pods     | Replica | streaming |  1 |         0 |
# | age-cluster-instance1-hml4-0 | age-cluster-pods     | Replica | streaming |  1 |         0 |
# | age-cluster-instance1-nhff-0 | age-cluster-pods     | Leader  | running   |  1 |           |
# +------------------------------+----------------------+---------+-----------+----+-----------+
```

## Known Issues and Limitations

### 1. Security Context
- Containers run as root due to AGE image requirements
- Not compliant with Kubernetes "Restricted" Pod Security Standards
- Recommended: Implement network policies and RBAC restrictions

### 2. Initialization Timing
- Initial replica setup may take longer due to pgBackRest configuration
- Transient "repo1-path cannot be set multiple times" errors may appear during initialization
- These errors typically resolve automatically through operator reconciliation

### 3. CRD Size Limitation
- PostgresCluster CRD may exceed Kubernetes annotation size limits
- Workaround: Use `--server-side` flag when applying CRDs

### 4. Namespace Deployment
- Clusters are deployed in the `postgres-operator` namespace by default
- Cross-namespace deployments require additional RBAC configuration

## Troubleshooting

### Pod Not Ready

If a pod shows 3/4 containers ready:

```bash
# Check pod events
kubectl describe pod -n postgres-operator <pod-name>

# Check pgBackRest logs
kubectl logs -n postgres-operator <pod-name> -c pgbackrest

# Common issue: pgBackRest configuration conflict
# Solution: Wait for automatic reconciliation (typically resolves within 2-3 minutes)
```

### AGE Extension Not Found

```bash
# Verify AGE is installed in the image
kubectl exec -n postgres-operator <pod-name> -c database -- psql -c "SELECT * FROM pg_available_extensions WHERE name = 'age';"

# Manually create the extension if needed
kubectl exec -n postgres-operator <pod-name> -c database -- psql -c "CREATE EXTENSION age;"
```

### Replica Lag Issues

```bash
# Check replication status
kubectl exec -n postgres-operator <primary-pod> -c database -- psql -c "SELECT * FROM pg_stat_replication;"

# Check Patroni logs
kubectl logs -n postgres-operator <pod-name> -c database | grep -i patroni
```

## Advanced Operations

### Accessing Clusters

#### Port Forwarding (Development)

```bash
# Get cluster service
kubectl get svc -n postgres-operator

# Port forward to primary service
kubectl port-forward -n postgres-operator svc/age-cluster-primary 5432:5432

# Connect with psql
psql -h localhost -p 5432 -U postgres -d postgres
```

#### Using pgBouncer (Production)

```bash
# Port forward to pgBouncer
kubectl port-forward -n postgres-operator svc/age-cluster-pgbouncer 6432:5432

# Connect through connection pooler
psql -h localhost -p 6432 -U postgres -d postgres
```

#### Get Connection Credentials

```bash
# Get password
kubectl get secret -n postgres-operator age-cluster-pguser-postgres -o jsonpath='{.data.password}' | base64 -d

# Get connection string
kubectl get secret -n postgres-operator age-cluster-pguser-postgres -o jsonpath='{.data.uri}' | base64 -d

# Get all connection details
kubectl get secret -n postgres-operator age-cluster-pguser-postgres -o yaml
```

### Monitoring Clusters

#### Cluster Health

```bash
# Check cluster status
kubectl get postgrescluster -n postgres-operator age-cluster -o jsonpath='{.status.conditions[*]}'

# View Patroni status
kubectl exec -n postgres-operator -it $(kubectl get pod -n postgres-operator -l postgres-operator.crunchydata.com/cluster=age-cluster,postgres-operator.crunchydata.com/role=master -o jsonpath='{.items[0].metadata.name}') -c database -- patronictl list

# Check replication lag
kubectl exec -n postgres-operator -it $(kubectl get pod -n postgres-operator -l postgres-operator.crunchydata.com/cluster=age-cluster,postgres-operator.crunchydata.com/role=master -o jsonpath='{.items[0].metadata.name}') -c database -- psql -c "SELECT client_addr, state, sync_state, replay_lag FROM pg_stat_replication;"
```

#### Resource Usage

```bash
# Check pod resource usage
kubectl top pods -n postgres-operator -l postgres-operator.crunchydata.com/cluster=age-cluster

# View PVC usage
kubectl get pvc -n postgres-operator -l postgres-operator.crunchydata.com/cluster=age-cluster

# Database size
kubectl exec -n postgres-operator -it $(kubectl get pod -n postgres-operator -l postgres-operator.crunchydata.com/cluster=age-cluster,postgres-operator.crunchydata.com/role=master -o jsonpath='{.items[0].metadata.name}') -c database -- psql -c "SELECT pg_database.datname, pg_size_pretty(pg_database_size(pg_database.datname)) AS size FROM pg_database;"
```

### Performing Updates

#### Update PostgreSQL Configuration

```yaml
# Save as postgres-config.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: age-cluster-config
data:
  postgresql.conf: |
    max_connections = 200
    shared_buffers = 256MB
    effective_cache_size = 1GB
    maintenance_work_mem = 64MB
    checkpoint_completion_target = 0.9
    wal_buffers = 16MB
    default_statistics_target = 100
    random_page_cost = 1.1
    effective_io_concurrency = 200
    work_mem = 4MB
    min_wal_size = 1GB
    max_wal_size = 4GB
```

Apply configuration:
```bash
# Apply config
kubectl apply -f postgres-config.yaml

# Update cluster to use config
kubectl patch postgrescluster age-cluster -n postgres-operator --type='merge' -p='{"spec":{"customReplicationClientTLS":{"config":"age-cluster-config"}}}'

# Restart to apply changes
kubectl rollout restart statefulset -n postgres-operator -l postgres-operator.crunchydata.com/cluster=age-cluster
```

#### Upgrade AGE Image

```bash
# Build new image
docker build -f Dockerfile.age -t localhost/postgres-age-patroni:v2 .

# Load into cluster (Kind)
kind load docker-image localhost/postgres-age-patroni:v2

# Update cluster
kubectl patch postgrescluster age-cluster -n postgres-operator --type='merge' -p='{"spec":{"image":"localhost/postgres-age-patroni:v2"}}'
```

### Backup Operations

#### Manual Backup

```bash
# Trigger full backup
kubectl exec -n postgres-operator -it $(kubectl get pod -n postgres-operator -l postgres-operator.crunchydata.com/cluster=age-cluster,postgres-operator.crunchydata.com/role=master -o jsonpath='{.items[0].metadata.name}') -c database -- pgbackrest backup --stanza=db --type=full

# Trigger incremental backup
kubectl exec -n postgres-operator -it $(kubectl get pod -n postgres-operator -l postgres-operator.crunchydata.com/cluster=age-cluster,postgres-operator.crunchydata.com/role=master -o jsonpath='{.items[0].metadata.name}') -c database -- pgbackrest backup --stanza=db --type=incr

# Check backup info
kubectl exec -n postgres-operator -it $(kubectl get pod -n postgres-operator -l postgres-operator.crunchydata.com/cluster=age-cluster,postgres-operator.crunchydata.com/role=master -o jsonpath='{.items[0].metadata.name}') -c database -- pgbackrest info
```

#### Restore from Backup

```yaml
# Save as restore-cluster.yaml
apiVersion: postgres-operator.crunchydata.com/v1beta1
kind: PostgresCluster
metadata:
  name: age-restored
spec:
  postgresVersion: 16
  image: localhost/postgres-age-patroni
  
  dataSource:
    pgbackrest:
      stanza: db
      configuration:
      - secret:
          name: age-cluster-pgbackrest-secrets
      global:
        repo1-path: /pgbackrest/repo1
      repo:
        name: repo1
        volume:
          volumeClaimSpec:
            accessModes:
            - "ReadWriteOnce"
            resources:
              requests:
                storage: 10Gi
  
  instances:
    - name: instance1
      replicas: 1
      dataVolumeClaimSpec:
        accessModes:
        - "ReadWriteOnce"
        resources:
          requests:
            storage: 10Gi
```

### Troubleshooting Operations

#### Debug Pod Issues

```bash
# Get pod logs
kubectl logs -n postgres-operator <pod-name> -c database --tail=50

# Get all container logs
kubectl logs -n postgres-operator <pod-name> --all-containers=true

# Execute commands in pod
kubectl exec -n postgres-operator -it <pod-name> -c database -- bash

# Check events
kubectl get events -n postgres-operator --sort-by='.lastTimestamp' | grep age-cluster
```

#### Common Issues and Solutions

1. **Pod Stuck in Pending**
```bash
# Check PVC status
kubectl get pvc -n postgres-operator

# Check node resources
kubectl describe nodes

# Check pod events
kubectl describe pod -n postgres-operator <pod-name>
```

2. **Replica Lag**
```bash
# Check replication status
kubectl exec -n postgres-operator -it <primary-pod> -c database -- psql -c "SELECT * FROM pg_stat_replication;"

# Force replica rebuild
kubectl delete pod -n postgres-operator <replica-pod>
```

3. **Connection Issues**
```bash
# Test internal connectivity
kubectl run -n postgres-operator test-connection --image=postgres:16 --rm -it --restart=Never -- psql -h age-cluster-primary -U postgres

# Check service endpoints
kubectl get endpoints -n postgres-operator age-cluster-primary
```

## Maintenance and Operations

### Scaling the Cluster

```bash
# Edit the PostgresCluster resource
kubectl edit postgrescluster -n postgres-operator age-cluster

# Or use kubectl patch
kubectl patch postgrescluster -n postgres-operator age-cluster --type='json' -p='[{"op": "replace", "path": "/spec/instances/0/replicas", "value":5}]'
```

### Backup and Restore

```bash
# Manual backup
kubectl exec -n postgres-operator <primary-pod> -c database -- pgbackrest backup --type=full --stanza=db

# Check backup info
kubectl exec -n postgres-operator <primary-pod> -c database -- pgbackrest info
```

### Updating the AGE Image

1. Build new image with updated components
2. Push to registry or load into cluster
3. Update PostgresCluster spec with new image
4. Operator will perform rolling update

## Production Recommendations

1. **Security Hardening**
   - Implement Pod Security Policies or Pod Security Standards at namespace level
   - Use network policies to restrict pod-to-pod communication
   - Enable audit logging for database operations

2. **Resource Management**
   - Set appropriate resource requests and limits
   - Configure storage classes for production workloads
   - Enable monitoring with pgMonitor

3. **Backup Strategy**
   - Configure automated backups with retention policies
   - Test restore procedures regularly
   - Consider multi-repository backup configuration

4. **High Availability**
   - Deploy at least 3 replicas for production
   - Configure synchronous replication for critical workloads
   - Use anti-affinity rules to spread pods across nodes

## Contributing

To contribute to this AGE-enabled operator fork:

1. Report issues specific to AGE integration
2. Test AGE functionality with different PostgreSQL versions
3. Improve security posture while maintaining AGE compatibility
4. Add AGE-specific monitoring and metrics

## License

This fork maintains the original Apache 2.0 license of the Crunchy PostgreSQL Operator. Apache AGE is also licensed under Apache 2.0.

## References

- [Apache AGE Documentation](https://age.apache.org/)
- [Crunchy PostgreSQL Operator Documentation](https://access.crunchydata.com/documentation/postgres-operator/)
- [Kubernetes Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
- [Patroni Documentation](https://patroni.readthedocs.io/)
- [pgBackRest Documentation](https://pgbackrest.org/)