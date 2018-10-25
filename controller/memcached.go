package controller

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/client/k8scrdclient"
	"github.com/giantswarm/operatorkit/controller"
	"github.com/giantswarm/operatorkit/informer"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"

	"github.com/giantswarm/cll-operator-workshop/controller/resource"
	workshopv1alpha1 "github.com/giantswarm/cll-operator-workshop/pkg/apis/workshop/v1alpha1"
	"github.com/giantswarm/cll-operator-workshop/pkg/clientset/versioned"
	"github.com/giantswarm/cll-operator-workshop/pkg/logger"
)

const name = "memcached-operator"

// Config represents the configuration used to create a new memcached controller.
type Config struct {
	K8sClient    kubernetes.Interface
	K8sExtClient apiextensionsclient.Interface
	CLLClient    versioned.Interface
}

// Memcached is a type containing the OperatorKit controller.
type Memcached struct {
	*controller.Controller
}

// New creates a new memcached controller.
func NewMemcached(config Config) (*Memcached, error) {
	var err error

	var (
		crd        = workshopv1alpha1.NewMemcachedConfigCRD()
		restClient = config.CLLClient.WorkshopV1alpha1().RESTClient()
		watcher    = config.CLLClient.WorkshopV1alpha1().MemcachedConfigs("")
	)

	var deploymentsResource controller.Resource
	{
		c := resource.DeploymentsConfig{
			K8sClient: config.K8sClient,
		}

		deploymentsResource, err = resource.NewDeployments(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var servicesResource controller.Resource
	{
		c := resource.ServicesConfig{
			K8sClient: config.K8sClient,
		}

		servicesResource, err = resource.NewServices(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []controller.Resource{
		deploymentsResource,
		servicesResource,
	}

	// Below is a common controller wiring. This code doesn't change unless
	// you want to reconcile non-CRD objects or to use more sophisticated
	// object routing.

	var crdClient *k8scrdclient.CRDClient
	{
		c := k8scrdclient.Config{
			Logger: logger.Default,

			K8sExtClient: config.K8sExtClient,
		}

		crdClient, err = k8scrdclient.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

	}

	resourceRouter, err := newSimpleResourceRouter(resources)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var memcachedInformer *informer.Informer
	{
		c := informer.Config{
			Logger: logger.Default,

			Watcher: watcher,
		}

		memcachedInformer, err = informer.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

	}

	var underlying *controller.Controller
	{
		c := controller.Config{
			Logger: logger.Default,
			Name:   name,

			CRD:        crd,
			CRDClient:  crdClient,
			Informer:   memcachedInformer,
			RESTClient: restClient,

			ResourceRouter: resourceRouter,
		}

		underlying, err = controller.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	m := &Memcached{
		Controller: underlying,
	}

	return m, nil
}
