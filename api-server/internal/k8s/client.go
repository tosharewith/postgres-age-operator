package k8s

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/tosharewith/postgres-age-operator/api-server/internal/models"
)

// Client wraps the Kubernetes client with customer instance operations
type Client struct {
	clientset kubernetes.Interface
	config    *rest.Config
}

// NewClient creates a new Kubernetes client
func NewClient() (*Client, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &Client{
		clientset: clientset,
		config:    config,
	}, nil
}

// getKubernetesConfig returns the Kubernetes configuration
func getKubernetesConfig() (*rest.Config, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Fall back to kubeconfig
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
	}

	return config, nil
}

// CreateCustomerInstance creates a new customer instance
func (c *Client) CreateCustomerInstance(ctx context.Context, customer *models.CustomerInstance) error {
	// Validate customer name
	if err := c.validateCustomerName(customer.Name); err != nil {
		return err
	}

	// Check if instance already exists
	exists, err := c.customerInstanceExists(ctx, customer.Name)
	if err != nil {
		return fmt.Errorf("failed to check if customer exists: %w", err)
	}
	if exists {
		return fmt.Errorf("customer instance '%s' already exists", customer.Name)
	}

	// Create namespace
	if err := c.createNamespace(ctx, customer); err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	// Create service account
	if err := c.createServiceAccount(ctx, customer); err != nil {
		return fmt.Errorf("failed to create service account: %w", err)
	}

	// Create cluster role
	if err := c.createClusterRole(ctx, customer); err != nil {
		return fmt.Errorf("failed to create cluster role: %w", err)
	}

	// Create cluster role binding
	if err := c.createClusterRoleBinding(ctx, customer); err != nil {
		return fmt.Errorf("failed to create cluster role binding: %w", err)
	}

	// Create deployment
	if err := c.createDeployment(ctx, customer); err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	return nil
}

// DeleteCustomerInstance deletes a customer instance
func (c *Client) DeleteCustomerInstance(ctx context.Context, customerName string) error {
	namespace := c.getNamespaceName(customerName)

	// Delete namespace (this will cascade delete most resources)
	err := c.clientset.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete namespace: %w", err)
	}

	// Delete cluster role
	clusterRoleName := c.getClusterRoleName(customerName)
	err = c.clientset.RbacV1().ClusterRoles().Delete(ctx, clusterRoleName, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete cluster role: %w", err)
	}

	// Delete cluster role binding
	clusterRoleBindingName := c.getClusterRoleBindingName(customerName)
	err = c.clientset.RbacV1().ClusterRoleBindings().Delete(ctx, clusterRoleBindingName, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete cluster role binding: %w", err)
	}

	return nil
}

// GetCustomerInstance retrieves a customer instance
func (c *Client) GetCustomerInstance(ctx context.Context, customerName string) (*models.CustomerInstance, error) {
	namespace := c.getNamespaceName(customerName)

	// Check if namespace exists
	ns, err := c.clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, fmt.Errorf("customer instance '%s' not found", customerName)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}

	// Get deployment
	deploymentName := c.getDeploymentName(customerName)
	deployment, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	// Build customer instance from Kubernetes resources
	customer := &models.CustomerInstance{
		ID:        customerName,
		Name:      customerName,
		Namespace: namespace,
		CreatedAt: ns.CreationTimestamp.Time,
		UpdatedAt: deployment.Status.Conditions[len(deployment.Status.Conditions)-1].LastUpdateTime.Time,
		Status: models.InstanceStatus{
			Replicas:      deployment.Status.Replicas,
			ReadyReplicas: deployment.Status.ReadyReplicas,
			Ready:         deployment.Status.ReadyReplicas > 0,
			LastUpdated:   time.Now(),
		},
	}

	// Set status phase based on deployment status
	if deployment.Status.ReadyReplicas > 0 {
		customer.Status.Phase = models.InstancePhaseRunning
	} else if deployment.Status.Replicas > 0 {
		customer.Status.Phase = models.InstancePhaseCreating
	} else {
		customer.Status.Phase = models.InstancePhaseFailed
	}

	// Extract labels and annotations
	if ns.Labels != nil {
		customer.Labels = make(map[string]string)
		for k, v := range ns.Labels {
			if strings.HasPrefix(k, "app.kubernetes.io/") || strings.HasPrefix(k, "postgres-operator.") {
				customer.Labels[k] = v
			}
		}
	}

	// Extract image tag from deployment
	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		image := deployment.Spec.Template.Spec.Containers[0].Image
		parts := strings.Split(image, ":")
		if len(parts) > 1 {
			customer.ImageTag = parts[1]
		}
	}

	return customer, nil
}

// ListCustomerInstances lists all customer instances
func (c *Client) ListCustomerInstances(ctx context.Context) ([]models.CustomerInstance, error) {
	// List namespaces with customer label
	labelSelector := "app.kubernetes.io/client"
	namespaces, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	var customers []models.CustomerInstance
	for _, ns := range namespaces.Items {
		if !strings.HasPrefix(ns.Name, "postgres-operator-") {
			continue
		}

		// Extract customer name from namespace
		customerName := strings.TrimPrefix(ns.Name, "postgres-operator-")

		customer, err := c.GetCustomerInstance(ctx, customerName)
		if err != nil {
			// Log error but continue with other customers
			continue
		}

		customers = append(customers, *customer)
	}

	return customers, nil
}

// UpdateCustomerInstance updates a customer instance
func (c *Client) UpdateCustomerInstance(ctx context.Context, customerName string, updates *models.UpdateCustomerRequest) error {
	namespace := c.getNamespaceName(customerName)
	deploymentName := c.getDeploymentName(customerName)

	// Get current deployment
	deployment, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	// Update image tag if provided
	if updates.ImageTag != "" {
		for i := range deployment.Spec.Template.Spec.Containers {
			if deployment.Spec.Template.Spec.Containers[i].Name == "operator" {
				imageParts := strings.Split(deployment.Spec.Template.Spec.Containers[i].Image, ":")
				if len(imageParts) > 0 {
					deployment.Spec.Template.Spec.Containers[i].Image = imageParts[0] + ":" + updates.ImageTag
				}
			}
		}
	}

	// Update labels if provided
	if updates.Labels != nil {
		if deployment.Spec.Template.Labels == nil {
			deployment.Spec.Template.Labels = make(map[string]string)
		}
		for k, v := range updates.Labels {
			deployment.Spec.Template.Labels[k] = v
		}
	}

	// Apply updates
	_, err = c.clientset.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}

// Helper methods for resource naming

func (c *Client) getNamespaceName(customerName string) string {
	return fmt.Sprintf("postgres-operator-%s", customerName)
}

func (c *Client) getServiceAccountName(customerName string) string {
	return fmt.Sprintf("%s-pgo-age", customerName)
}

func (c *Client) getClusterRoleName(customerName string) string {
	return fmt.Sprintf("postgres-operator-%s", customerName)
}

func (c *Client) getClusterRoleBindingName(customerName string) string {
	return fmt.Sprintf("postgres-operator-%s", customerName)
}

func (c *Client) getDeploymentName(customerName string) string {
	return fmt.Sprintf("%s-pgo-age", customerName)
}

func (c *Client) validateCustomerName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("customer name cannot be empty")
	}
	if len(name) > 20 {
		return fmt.Errorf("customer name cannot be longer than 20 characters")
	}
	if !isValidKubernetesName(name) {
		return fmt.Errorf("customer name must be a valid Kubernetes name (lowercase alphanumeric with optional hyphens)")
	}
	return nil
}

func (c *Client) customerInstanceExists(ctx context.Context, customerName string) (bool, error) {
	namespace := c.getNamespaceName(customerName)
	_, err := c.clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func isValidKubernetesName(name string) bool {
	// Simple validation - could be more comprehensive
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}
	return !strings.HasPrefix(name, "-") && !strings.HasSuffix(name, "-")
}