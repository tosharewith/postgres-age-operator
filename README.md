# PostgreSQL AGE Operator

<p align="center">
  <img width="200" src="https://age.apache.org/age-manual/master/_static/age_logo.png" alt="Apache AGE"/>
</p>

[![Go Report Card](https://goreportcard.com/badge/github.com/gregoriomomm/postgres-age-operator)](https://goreportcard.com/report/github.com/gregoriomomm/postgres-age-operator)
![GitHub Repo stars](https://img.shields.io/github/stars/gregoriomomm/postgres-age-operator)
[![License](https://img.shields.io/github/license/gregoriomomm/postgres-age-operator)](LICENSE.md)

# Production-Ready Graph Database on Kubernetes

The **PostgreSQL AGE Operator** brings the power of [Apache AGE](https://age.apache.org/) (A Graph Extension) to Kubernetes, providing a **declarative graph database** solution that automatically manages your [PostgreSQL](https://www.postgresql.org) clusters with native graph capabilities.

Based on the proven architecture of the Crunchy Data PostgreSQL Operator, this operator extends PostgreSQL with Apache AGE to deliver enterprise-grade graph database functionality with the reliability and features you expect from PostgreSQL.

## What is Apache AGE?

[Apache AGE](https://age.apache.org/) is an extension for PostgreSQL that provides graph database functionality. It allows you to:
- Store and query graph data using Cypher query language
- Combine SQL and Cypher in the same query
- Leverage PostgreSQL's reliability, ACID compliance, and ecosystem
- Build modern applications with relationships at their core

## Why PostgreSQL AGE Operator?

This operator makes it easy to deploy and manage PostgreSQL clusters with AGE on Kubernetes, providing:

✅ **Graph + Relational**: Best of both worlds - graph queries with Cypher, relational queries with SQL  
✅ **Cloud Native**: Designed for Kubernetes from the ground up  
✅ **High Availability**: Automatic failover and replica management  
✅ **Production Ready**: Based on battle-tested Crunchy PostgreSQL Operator  
✅ **GitOps Friendly**: Declarative configuration for your entire database stack

## Quick Start

Get a graph database running in under 5 minutes:

```bash
# Clone the repository
git clone https://github.com/gregoriomomm/postgres-age-operator
cd postgres-age-operator

# Build the AGE-enabled image
docker build -f Dockerfile.age -t localhost/postgres-age-patroni .

# For Kind clusters: Load the image
kind load docker-image localhost/postgres-age-patroni

# Deploy the operator
kubectl apply --server-side -k config/default

# Create an AGE cluster
kubectl apply -k examples/age-cluster/

# Connect and start using graphs!
kubectl exec -it -n postgres-operator \
  $(kubectl get pod -n postgres-operator -l postgres-operator.crunchydata.com/role=master -o name) \
  -c database -- psql
```  

## Features

### 🎯 Graph Database Capabilities

- **Cypher Query Language**: Industry-standard graph query language
- **Hybrid Queries**: Combine SQL and Cypher in the same query
- **Graph Algorithms**: Built-in support for common graph algorithms
- **Visual Data Modeling**: Natural representation of connected data

### 🚀 Enterprise Features

#### PostgreSQL Cluster [Provisioning][provisioning]

[Create, Scale, & Delete PostgreSQL clusters with ease][provisioning], while fully customizing your
Pods and PostgreSQL configuration!

#### [High Availability][high-availability]

Safe, automated failover backed by a [distributed consensus high availability solution][high-availability].
Uses [Pod Anti-Affinity][k8s-anti-affinity] to help resiliency; you can configure how aggressive this can be!
Failed primaries automatically heal, allowing for faster recovery time.

Support for [standby PostgreSQL clusters][multiple-cluster] that work both within and across [multiple Kubernetes clusters][multiple-cluster].

#### [Disaster Recovery][disaster-recovery]

[Backups][backups] and [restores][disaster-recovery] leverage the open source [pgBackRest][] utility and
[includes support for full, incremental, and differential backups as well as efficient delta restores][backups].
Set how long you to retain your backups. Works great with very large databases!

#### Security and [TLS][tls]

PGO enforces that all connections are over [TLS][tls]. You can also [bring your own TLS infrastructure][tls] if you do not want to use the defaults provided by PGO.

PGO runs containers with locked-down settings and provides Postgres credentials in a secure, convenient way for connecting your applications to your data.

#### [Monitoring][monitoring]

[Track the health of your PostgreSQL clusters][monitoring] using the open source [pgMonitor][] library.

#### [Upgrade Management][update-postgres]

Safely [apply PostgreSQL updates][update-postgres] with minimal impact to the availability of your PostgreSQL clusters.

#### Advanced Replication Support

Choose between [asynchronous][high-availability] and synchronous replication
for workloads that are sensitive to losing transactions.

#### [Clone][clone]

[Create new clusters from your existing clusters or backups][clone] with efficient data cloning.

#### [Connection Pooling][pool]

Advanced [connection pooling][pool] support using [pgBouncer][].

#### Pod Anti-Affinity, Node Affinity, Pod Tolerations

Have your PostgreSQL clusters deployed to [Kubernetes Nodes][k8s-nodes] of your preference. Set your [pod anti-affinity][k8s-anti-affinity], node affinity, Pod tolerations, and more rules to customize your deployment topology!

#### [Scheduled Backups][backup-management]

Choose the type of backup (full, incremental, differential) and [how frequently you want it to occur][backup-management] on each PostgreSQL cluster.

#### Backup to Local Storage, [S3][backups-s3], [GCS][backups-gcs], [Azure][backups-azure], or a Combo!

[Store your backups in Amazon S3][backups-s3] or any object storage system that supports
the S3 protocol. You can also store backups in [Google Cloud Storage][backups-gcs] and [Azure Blob Storage][backups-azure].

You can also [mix-and-match][backups-multi]: PGO lets you [store backups in multiple locations][backups-multi].

#### [Full Customizability][customize-cluster]

PGO makes it easy to fully customize your Postgres cluster to tailor to your workload:

- Choose the resources for your Postgres cluster: [container resources and storage size][resize-cluster]. [Resize at any time][resize-cluster] with minimal disruption.
- - Use your own container image repository, including support `imagePullSecrets` and private repositories
- [Customize your PostgreSQL configuration][customize-cluster]

#### [Namespaces][k8s-namespaces]

Deploy PGO to watch Postgres clusters in all of your [namespaces][k8s-namespaces], or [restrict which namespaces][single-namespace] you want PGO to manage Postgres clusters in!

[backups]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/backups-disaster-recovery/backups
[backups-s3]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/backups-disaster-recovery/backups#using-s3
[backups-gcs]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/backups-disaster-recovery/backups#using-google-cloud-storage-gcs
[backups-azure]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/backups-disaster-recovery/backups#using-azure-blob-storage
[backups-multi]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/backups-disaster-recovery/backups#set-up-multiple-backup-repositories
[backup-management]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/backups-disaster-recovery/backup-management
[clone]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/backups-disaster-recovery/disaster-recovery#clone-a-postgres-cluster
[customize-cluster]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/day-two/customize-cluster
[disaster-recovery]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/backups-disaster-recovery/disaster-recovery
[high-availability]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/day-two/high-availability/
[monitoring]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/day-two/monitoring/
[multiple-cluster]: https://access.crunchydata.com/documentation/postgres-operator/v5/architecture/disaster-recovery/#standby-cluster-overview
[pool]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/basic-setup/connection-pooling/
[provisioning]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/basic-setup/create-cluster/
[resize-cluster]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/cluster-management/resize-cluster/
[single-namespace]: https://access.crunchydata.com/documentation/postgres-operator/v5/installation/kustomize/#installation-mode
[tls]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/day-two/customize-cluster#customize-tls
[update-postgres]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/cluster-management/update-cluster
[k8s-anti-affinity]: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#inter-pod-affinity-and-anti-affinity
[k8s-namespaces]: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
[k8s-nodes]: https://kubernetes.io/docs/concepts/architecture/nodes/
[pgBackRest]: https://www.pgbackrest.org
[pgBouncer]: https://access.crunchydata.com/documentation/postgres-operator/v5/tutorials/basic-setup/connection-pooling/
[pgMonitor]: https://github.com/CrunchyData/pgmonitor

## Included Components

[PostgreSQL containers](https://github.com/CrunchyData/crunchy-containers) deployed with the PostgreSQL Operator include the following components:

- [PostgreSQL](https://www.postgresql.org)
  - [PostgreSQL Contrib Modules](https://www.postgresql.org/docs/current/contrib.html)
  - [PL/Python + PL/Python 3](https://www.postgresql.org/docs/current/plpython.html)
  - [PL/Perl](https://www.postgresql.org/docs/current/plperl.html)
  - [PL/Tcl](https://www.postgresql.org/docs/current/pltcl.html)
  - [pgAudit](https://www.pgaudit.org/)
  - [pgAudit Analyze](https://github.com/pgaudit/pgaudit_analyze)
  - [pg_cron](https://github.com/citusdata/pg_cron)
  - [pg_partman](https://github.com/pgpartman/pg_partman)
  - [pgnodemx](https://github.com/CrunchyData/pgnodemx)
  - [set_user](https://github.com/pgaudit/set_user)
  - [TimescaleDB](https://github.com/timescale/timescaledb) (Apache-licensed community edition)
  - [wal2json](https://github.com/eulerto/wal2json)
- [pgBackRest](https://pgbackrest.org/)
- [pgBouncer](http://pgbouncer.github.io/)
- [pgAdmin 4](https://www.pgadmin.org/)
- [pgMonitor](https://github.com/CrunchyData/pgmonitor)
- [Patroni](https://patroni.readthedocs.io/)
- [LLVM](https://llvm.org/) (for [JIT compilation](https://www.postgresql.org/docs/current/jit.html))

In addition to the above, the geospatially enhanced PostgreSQL + PostGIS container adds the following components:

- [PostGIS](http://postgis.net/)
- [pgRouting](https://pgrouting.org/)

[PostgreSQL Operator Monitoring](https://access.crunchydata.com/documentation/postgres-operator/latest/architecture/monitoring/) uses the following components:

- [pgMonitor](https://github.com/CrunchyData/pgmonitor)
- [Prometheus](https://github.com/prometheus/prometheus)
- [Grafana](https://github.com/grafana/grafana)
- [Alertmanager](https://github.com/prometheus/alertmanager)

For more information about which versions of the PostgreSQL Operator include which components, please visit the [compatibility](https://access.crunchydata.com/documentation/postgres-operator/v5/references/components/) section of the documentation.

## Supported Platforms

The PostgreSQL AGE Operator is tested on the following platforms:

- Kubernetes 1.21+
- OpenShift 4.8+
- Rancher
- Google Kubernetes Engine (GKE)
- Amazon EKS
- Microsoft AKS
- VMware Tanzu
- Kind (for local development)
- Minikube

## Installation

### Prerequisites

- Kubernetes 1.21+ or OpenShift 4.8+
- kubectl configured
- Docker for building images

### Quick Install on Kind

```bash
# Create a Kind cluster
kind create cluster --name age-demo

# Build and load the image
docker build -f Dockerfile.age -t localhost/postgres-age-patroni .
kind load docker-image localhost/postgres-age-patroni --name age-demo

# Deploy operator
kubectl apply --server-side -k config/default

# Create your first graph database
kubectl apply -k examples/age-cluster/
```

For production deployments, see [AGE-INTEGRATION.md](AGE-INTEGRATION.md).

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Areas where we need help:
- Testing on different Kubernetes distributions
- Additional graph algorithm implementations
- Performance benchmarking
- Documentation improvements
- Example applications

## Support

- **Issues**: [GitHub Issues](https://github.com/gregoriomomm/postgres-age-operator/issues)
- **Discussions**: [GitHub Discussions](https://github.com/gregoriomomm/postgres-age-operator/discussions)
- **Apache AGE**: [AGE Documentation](https://age.apache.org/)

## Documentation

- [Installation Guide](AGE-INTEGRATION.md#installation-options) - Detailed setup instructions
- [Operator Management](AGE-INTEGRATION.md#operator-management) - Managing the operator
- [Deployment Configurations](AGE-INTEGRATION.md#deployment-configurations) - Various deployment patterns
- [Advanced Operations](AGE-INTEGRATION.md#advanced-operations) - Day-2 operations
- [Troubleshooting](AGE-INTEGRATION.md#troubleshooting-operations) - Common issues and solutions

# Releases

When a PostgreSQL Operator general availability (GA) release occurs, the container images are distributed on the following platforms in order:

- [Crunchy Data Customer Portal](https://access.crunchydata.com/)
- [Crunchy Data Developer Portal](https://www.crunchydata.com/developers)

The image rollout can occur over the course of several days.

To stay up-to-date on when releases are made available in the [Crunchy Data Developer Portal](https://www.crunchydata.com/developers), please sign up for the [Crunchy Data Developer Program Newsletter](https://www.crunchydata.com/developers#email). You can also [join the PGO project community discord](https://discord.gg/a7vWKG8Ec9)

# FAQs, License and Terms

For more information regarding PGO, the Postgres Operator project from Crunchy Data, and Crunchy Postgres for Kubernetes, please see the [frequently asked questions](https://access.crunchydata.com/documentation/postgres-operator/latest/faq). 

The installation instructions provided in this repo are designed for the use of PGO along with Crunchy Data's Postgres distribution, Crunchy Postgres, as Crunchy Postgres for Kubernetes. The unmodified use of these installation instructions will result in downloading container images from Crunchy Data repositories - specifically the Crunchy Data Developer Portal. The use of container images downloaded from the Crunchy Data Developer Portal are subject to the [Crunchy Data Developer Program terms](https://www.crunchydata.com/developers/terms-of-use).  

The PGO Postgres Operator project source code is available subject to the [Apache 2.0 license](LICENSE.md) with the PGO logo and branding assets covered by [our trademark guidelines](docs/static/logos/TRADEMARKS.md).
