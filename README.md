# PostgreSQL AGE Operator

<p align="center">
  <img width="200" src="data:image/png;base64,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" alt="Apache AGE"/>
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
