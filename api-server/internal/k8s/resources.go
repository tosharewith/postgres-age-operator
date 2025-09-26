package k8s

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/tosharewith/postgres-age-operator/api-server/internal/models"
)

// createNamespace creates a namespace for the customer instance
func (c *Client) createNamespace(ctx context.Context, customer *models.CustomerInstance) error {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: c.getNamespaceName(customer.Name),
			Labels: map[string]string{
				"app.kubernetes.io/name":     "postgres-age-operator",
				"app.kubernetes.io/instance": customer.Name + "-age",
				"app.kubernetes.io/client":   customer.Name,
				"app.kubernetes.io/managed-by": "age-api-server",
			},
		},
	}

	// Add custom labels if provided
	if customer.Labels != nil {
		for k, v := range customer.Labels {
			namespace.Labels[k] = v
		}
	}

	_, err := c.clientset.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	return err
}

// createServiceAccount creates a service account for the customer instance
func (c *Client) createServiceAccount(ctx context.Context, customer *models.CustomerInstance) error {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.getServiceAccountName(customer.Name),
			Namespace: c.getNamespaceName(customer.Name),
			Labels: map[string]string{
				"app.kubernetes.io/name":     "postgres-age-operator",
				"app.kubernetes.io/instance": customer.Name + "-age",
				"app.kubernetes.io/client":   customer.Name,
			},
		},
	}

	_, err := c.clientset.CoreV1().ServiceAccounts(c.getNamespaceName(customer.Name)).Create(ctx, serviceAccount, metav1.CreateOptions{})
	return err
}

// createClusterRole creates a cluster role for the customer instance
func (c *Client) createClusterRole(ctx context.Context, customer *models.CustomerInstance) error {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: c.getClusterRoleName(customer.Name),
			Labels: map[string]string{
				"app.kubernetes.io/name":     "postgres-age-operator",
				"app.kubernetes.io/instance": customer.Name + "-age",
				"app.kubernetes.io/client":   customer.Name,
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{
					"configmaps",
					"endpoints",
					"events",
					"persistentvolumeclaims",
					"pods",
					"secrets",
					"serviceaccounts",
					"services",
				},
				Verbs: []string{"*"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{
					"daemonsets",
					"deployments",
					"replicasets",
					"statefulsets",
				},
				Verbs: []string{"*"},
			},
			{
				APIGroups: []string{"batch"},
				Resources: []string{
					"cronjobs",
					"jobs",
				},
				Verbs: []string{"*"},
			},
			{
				APIGroups: []string{"postgres-operator.crunchydata.com"},
				Resources: []string{
					"postgresclusters",
					"pgadmins",
					"pgupgrades",
				},
				Verbs: []string{"*"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{"get", "list"},
			},
			{
				APIGroups: []string{"policy"},
				Resources: []string{"poddisruptionbudgets"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{
					"rolebindings",
					"roles",
				},
				Verbs: []string{"*"},
			},
		},
	}

	_, err := c.clientset.RbacV1().ClusterRoles().Create(ctx, clusterRole, metav1.CreateOptions{})
	return err
}

// createClusterRoleBinding creates a cluster role binding for the customer instance
func (c *Client) createClusterRoleBinding(ctx context.Context, customer *models.CustomerInstance) error {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: c.getClusterRoleBindingName(customer.Name),
			Labels: map[string]string{
				"app.kubernetes.io/name":     "postgres-age-operator",
				"app.kubernetes.io/instance": customer.Name + "-age",
				"app.kubernetes.io/client":   customer.Name,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     c.getClusterRoleName(customer.Name),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      c.getServiceAccountName(customer.Name),
				Namespace: c.getNamespaceName(customer.Name),
			},
		},
	}

	_, err := c.clientset.RbacV1().ClusterRoleBindings().Create(ctx, clusterRoleBinding, metav1.CreateOptions{})
	return err
}

// createDeployment creates the operator deployment for the customer instance
func (c *Client) createDeployment(ctx context.Context, customer *models.CustomerInstance) error {
	imageTag := "latest"
	if customer.ImageTag != "" {
		imageTag = customer.ImageTag
	}

	// Default resource requirements
	requests := corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse("100m"),
		corev1.ResourceMemory: resource.MustParse("128Mi"),
	}
	limits := corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse("500m"),
		corev1.ResourceMemory: resource.MustParse("512Mi"),
	}

	// Apply custom resource requirements if specified
	if customer.Config.Resources.Requests.CPU != "" {
		requests[corev1.ResourceCPU] = resource.MustParse(customer.Config.Resources.Requests.CPU)
	}
	if customer.Config.Resources.Requests.Memory != "" {
		requests[corev1.ResourceMemory] = resource.MustParse(customer.Config.Resources.Requests.Memory)
	}
	if customer.Config.Resources.Limits.CPU != "" {
		limits[corev1.ResourceCPU] = resource.MustParse(customer.Config.Resources.Limits.CPU)
	}
	if customer.Config.Resources.Limits.Memory != "" {
		limits[corev1.ResourceMemory] = resource.MustParse(customer.Config.Resources.Limits.Memory)
	}

	replicas := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.getDeploymentName(customer.Name),
			Namespace: c.getNamespaceName(customer.Name),
			Labels: map[string]string{
				"app.kubernetes.io/name":     "postgres-age-operator",
				"app.kubernetes.io/instance": customer.Name + "-age",
				"app.kubernetes.io/client":   customer.Name,
				"postgres-operator.crunchydata.com/control-plane": "postgres-operator-" + customer.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":     "postgres-age-operator",
					"app.kubernetes.io/instance": customer.Name + "-age",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":     "postgres-age-operator",
						"app.kubernetes.io/instance": customer.Name + "-age",
						"app.kubernetes.io/client":   customer.Name,
						"postgres-operator.crunchydata.com/control-plane": "postgres-operator-" + customer.Name,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: c.getServiceAccountName(customer.Name),
					Containers: []corev1.Container{
						{
							Name:            "operator",
							Image:           fmt.Sprintf("localhost/postgres-age-operator:%s", imageTag),
							ImagePullPolicy: corev1.PullNever,
							Resources: corev1.ResourceRequirements{
								Requests: requests,
								Limits:   limits,
							},
							Env: []corev1.EnvVar{
								{
									Name:  "PGO_INSTALLER",
									Value: "api-server",
								},
								{
									Name:  "PGO_INSTALLER_ORIGIN",
									Value: fmt.Sprintf("api-server-customer-%s", customer.Name),
								},
								{
									Name: "PGO_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name:  "CRUNCHY_DEBUG",
									Value: "false",
								},
								{
									Name:  "RELATED_IMAGE_POSTGRES_16",
									Value: "localhost/postgres-age-patroni",
								},
								{
									Name:  "RELATED_IMAGE_POSTGRES_17",
									Value: "localhost/postgres-age-patroni",
								},
								{
									Name:  "RELATED_IMAGE_PGBACKREST",
									Value: "registry.developers.crunchydata.com/crunchydata/crunchy-pgbackrest:ubi9-2.54.2-2520",
								},
								{
									Name:  "RELATED_IMAGE_PGBOUNCER",
									Value: "registry.developers.crunchydata.com/crunchydata/crunchy-pgbouncer:ubi9-1.24-2520",
								},
								{
									Name:  "RELATED_IMAGE_PGEXPORTER",
									Value: "registry.developers.crunchydata.com/crunchydata/crunchy-postgres-exporter:ubi9-0.17.1-2520",
								},
								{
									Name:  "RELATED_IMAGE_PGUPGRADE",
									Value: "registry.developers.crunchydata.com/crunchydata/crunchy-upgrade:ubi9-17.5-2520",
								},
								{
									Name:  "RELATED_IMAGE_STANDALONE_PGADMIN",
									Value: "registry.developers.crunchydata.com/crunchydata/crunchy-pgadmin4:ubi9-9.2-2520",
								},
								{
									Name:  "RELATED_IMAGE_COLLECTOR",
									Value: "registry.developers.crunchydata.com/crunchydata/postgres-operator:ubi9-5.8.2-0",
								},
							},
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: &[]bool{false}[0],
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
								},
								ReadOnlyRootFilesystem: &[]bool{true}[0],
								RunAsNonRoot:          &[]bool{true}[0],
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/healthz",
										Port: intstr.FromInt(8081),
									},
								},
								InitialDelaySeconds: 15,
								PeriodSeconds:       20,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/readyz",
										Port: intstr.FromInt(8081),
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       10,
							},
						},
					},
				},
			},
		},
	}

	_, err := c.clientset.AppsV1().Deployments(c.getNamespaceName(customer.Name)).Create(ctx, deployment, metav1.CreateOptions{})
	return err
}