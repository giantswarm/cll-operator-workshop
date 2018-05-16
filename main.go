package main

import (
	"fmt"

	"github.com/giantswarm/microerror"
	flag "github.com/spf13/pflag"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/cll-operator-workshop/controller"
	"github.com/giantswarm/cll-operator-workshop/pkg/clientset/versioned"
	"github.com/giantswarm/cll-operator-workshop/pkg/logger"
	"github.com/giantswarm/operatorkit/client/k8srestconfig"
)

var (
	k8sAddress   string
	k8sInCluster bool
	k8sCAFile    string
	k8sCrtFile   string
	k8sKeyFile   string
)

func init() {
	flag.StringVar(&k8sAddress, "kubernetes.address", "", "Address used to connect to Kubernetes.")
	flag.BoolVar(&k8sInCluster, "kubernetes.incluster", true, "Whether to use the in-cluster config to authenticate with Kubernetes.")
	flag.StringVar(&k8sCAFile, "kubernetes.ca", "", "Certificate authority file path to use to authenticate with Kubernetes.")
	flag.StringVar(&k8sCrtFile, "kubernetes.crt", "", "Certificate file path to use to authenticate with Kubernetes.")
	flag.StringVar(&k8sKeyFile, "kubernetes.key", "", "Key file path to use to authenticate with Kubernetes.")
	flag.Parse()
}

func main() {
	err := mainWithError()
	if err != nil {
		panic(fmt.Sprintf("%#v\n", err))
	}
}

func mainWithError() error {
	var err error

	var restConfig *rest.Config
	{
		c := k8srestconfig.Config{
			Logger: logger.Default,

			Address:   k8sAddress,
			InCluster: k8sInCluster,
			TLS: k8srestconfig.TLSClientConfig{
				CAFile:  k8sCAFile,
				CrtFile: k8sCrtFile,
				KeyFile: k8sKeyFile,
			},
		}

		restConfig, err = k8srestconfig.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	k8sClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return microerror.Mask(err)
	}

	k8sExtClient, err := apiextensionsclient.NewForConfig(restConfig)
	if err != nil {
		return microerror.Mask(err)
	}

	cllClient, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return microerror.Mask(err)
	}

	var memcachedController *controller.Memcached
	{
		c := controller.Config{
			K8sClient:    k8sClient,
			K8sExtClient: k8sExtClient,
			CLLClient:    cllClient,
		}

		memcachedController, err = controller.NewMemcached(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	memcachedController.Boot()

	return nil
}
