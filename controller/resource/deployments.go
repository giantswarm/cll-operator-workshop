package resource

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/cll-operator-workshop/controller/key"
	workshopv1alpha1 "github.com/giantswarm/cll-operator-workshop/pkg/apis/workshop/v1alpha1"
	"github.com/giantswarm/cll-operator-workshop/pkg/logger"
)

const (
	deploymentsName = "deployments"
)

type DeploymentsConfig struct {
	K8sClient kubernetes.Interface
}

type Deployments struct {
	logger    micrologger.Logger
	k8sClient kubernetes.Interface
}

func NewDeployments(config DeploymentsConfig) (*Deployments, error) {
	d := &Deployments{
		logger:    logger.Default,
		k8sClient: config.K8sClient,
	}

	return d, nil
}
func (d *Deployments) Name() string {
	return deploymentsName
}

func (d *Deployments) EnsureCreated(ctx context.Context, obj interface{}) error {
	memcachedConfig := obj.(*workshopv1alpha1.MemcachedConfig).DeepCopy()

	managedDeployments, err := d.k8sClient.AppsV1().Deployments(memcachedConfig.Namespace).List(metav1.ListOptions{
		LabelSelector: key.LabelSelectorManagedBy,
	})
	if err != nil {
		return microerror.Mask(err)
	}

	currentReplicas := len(managedDeployments.Items)
	desiredReplicas := memcachedConfig.Spec.Replicas

	d.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("currentReplicas = %d desiredReplicas = %d", currentReplicas, desiredReplicas))

	// Update existing deployments and scale up if necessary.
	for i := 0; i < desiredReplicas; i++ {
		err := d.ensureReplicaCreated(ctx, memcachedConfig, i)
		if err != nil {
			return microerror.Mask(err)
		}

	}

	// Scale down if necessary.
	for i := desiredReplicas; i < currentReplicas; i++ {
		err := d.ensureReplicaDeleted(ctx, memcachedConfig, i)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (d *Deployments) EnsureDeleted(ctx context.Context, obj interface{}) error {
	memcachedConfig := obj.(*workshopv1alpha1.MemcachedConfig).DeepCopy()

	managedDeployments, err := d.k8sClient.AppsV1().Deployments(memcachedConfig.Namespace).List(metav1.ListOptions{
		LabelSelector: key.LabelSelectorManagedBy,
	})
	if err != nil {
		return microerror.Mask(err)
	}

	currentReplicas := len(managedDeployments.Items)

	d.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("currentReplicas = %d", currentReplicas))

	// Delete existing deploymentes.
	for i := 0; i < currentReplicas; i++ {
		err := d.ensureReplicaDeleted(ctx, memcachedConfig, i)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (d *Deployments) ensureReplicaCreated(ctx context.Context, m *workshopv1alpha1.MemcachedConfig, replica int) error {
	desired, err := newDesiredDeployment(m, replica)
	if err != nil {
		return microerror.Mask(err)
	}

	current, err := d.k8sClient.AppsV1().Deployments(desired.Namespace).Get(desired.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		// Just make sure current is nil when not found.
		current = nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	if current == nil {
		_, err = d.k8sClient.AppsV1().Deployments(desired.Namespace).Create(desired)
		if err != nil {
			return microerror.Mask(err)
		}

		d.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deployment %s/%s created", desired.Namespace, desired.Name))
	} else {
		desired.ResourceVersion = current.ResourceVersion

		_, err = d.k8sClient.AppsV1().Deployments(desired.Namespace).Update(desired)
		if err != nil {
			return microerror.Mask(err)
		}

		d.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deployment %s/%s updated", desired.Namespace, desired.Name))
	}

	return nil
}

func (d *Deployments) ensureReplicaDeleted(ctx context.Context, m *workshopv1alpha1.MemcachedConfig, replica int) error {
	name := key.ReplicaName(replica)
	namespace := key.Namespace(m)

	err := d.k8sClient.AppsV1().Deployments(namespace).Delete(name, &metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		d.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deployment %s/%s already deleted", namespace, name))
	} else if err != nil {
		return microerror.Mask(err)
	} else {
		d.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("deployment %s/%s deleted", namespace, name))
	}

	return nil
}

func newDesiredDeployment(memcachedConfig *workshopv1alpha1.MemcachedConfig, replica int) (*appsv1.Deployment, error) {
	name := key.ReplicaName(replica)
	namespace := key.Namespace(memcachedConfig)

	memoryMB, err := key.MemoryMB(memcachedConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	labels := key.CommonLabels(memcachedConfig, replica)

	replicas := int32(1)

	d := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "memcached:" + key.MemcachedVersion,
						Name:  "memcached",
						Command: []string{
							"memcached",
							"-v",
							fmt.Sprintf("-m=%d", memoryMB),
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: key.MemcachedPort,
							Name:          "memcached",
						}},
					}},
				},
			},
		},
	}

	return d, nil
}
