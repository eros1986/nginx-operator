// This file implements an operator that uses `IndexedInformer`.

package k8s

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// IndexedHandlerInterface defines the functionalities of IndexedHandler.
type IndexedHandlerInterface interface {
	HandlerInterface
}

// IndexedOperatorConfig is the configuration to construct an operator.
type IndexedOperatorConfig struct {
	// kubeconfig path
	KubeConfigPath string

	// WatchNamespace represents the namespace that the Operator lists/watch.
	// Empty string means any namespace.
	// Default: default
	WatchNamespace string

	// ResyncPeriod is the period that will re-list the objects.
	// If non-zero, it will re-list this often (you will get OnUpdate calls,
	// even if nothing changed). Otherwise, re-list will be deplayed as long
	// as possible (until the upstream source closes the watch or times out,
	// or you stop the controller).
	ResyncPeriod time.Duration

	// Handlers are the implementaion of IndexedHandlerInterface.
	// Note: This is implemented outside the operator.
	Handlers IndexedHandlerInterface

	// Logger is a logrus logger.
	Logger *log.Entry
}

// IndexedOperator is the construct that implements an IndexedOperator.
type IndexedOperator struct {
	// config is the raw IndexedOperatorConfig.
	config *IndexedOperatorConfig

	// kubeconfig is the k8s configuration.
	kubeconfig *rest.Config

	// kclient is the k8s clientset.
	kclient *kubernetes.Clientset
	// aeclient is the apiextensions clientset.
	aeclient *apiextensionsclient.Clientset

	// crdI is the CRD interface.
	crdI CRDInterface

	// logger is the root logger for this operator.
	Logger *log.Entry

	// Indexer is an cache.Indexer used to index objects in the cache.
	Indexer cache.Indexer
}

func init() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
}

// NewIndexedOperator news an IndexedOperator.
func NewIndexedOperator(config *IndexedOperatorConfig) (OperatorInterface, error) {
	if err := checkIndexedOperatorConfig(config); err != nil {
		return nil, err
	}

	kubeconfig, err := BuildKuberentesConfig(config.KubeConfigPath)
	if err != nil {
		return nil, err
	}

	clientsets := NewClientSets(kubeconfig)
	kclient, err := clientsets.GetKClientset()
	if err != nil {
		return nil, err
	}
	aeclient, err := clientsets.GetAPIExtensionClientset()
	if err != nil {
		return nil, err
	}

	crdI := NewCRDs(aeclient)

	return &IndexedOperator{
		config:     config,
		kubeconfig: kubeconfig,
		kclient:    kclient,
		aeclient:   aeclient,
		crdI:       crdI,
		Logger:     config.Logger,
	}, nil
}

// CreateCRD creates a CRD resource.
func (i *IndexedOperator) CreateCRD(crd *CRD) error {
	crdData := &CRDData{
		Name: crd.Name,
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   crd.Group,
			Version: crd.Version,
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Kind:   crd.Kind,
				Plural: crd.Plural,
			},
			Scope: crd.Scope,
		},
	}
	crdConfig := i.crdI.MakeConfig(crdData)
	_, err := i.crdI.Create(crdConfig)

	return err
}

// DeleteCRD deletes a CRD resource.
func (i *IndexedOperator) DeleteCRD(name string, options *metav1.DeleteOptions) error {
	return i.crdI.Delete(name, options)
}

// GetCRD gets a CRD resource
func (i *IndexedOperator) GetCRD(name string) (*apiextensionsv1beta1.CustomResourceDefinition, error) {
	return i.crdI.Get(name)
}

// WatchEvents starts watching the CRD's ADD/DELETE/UPDATE events.
func (i *IndexedOperator) WatchEvents(ctx context.Context, crd *CRD) error {
	crdRestClientConfig := &CRDRestClientConfig{
		KubernetesConfig: i.kubeconfig,
		CRDGroup:         crd.Group,
		CRDVersion:       crd.Version,
		CRDObj:           crd.Obj,
		CRDObjList:       crd.ObjList,
		SchemeBuilder:    crd.SchemeBuilder,
	}
	crdClient, _, err := i.crdI.NewRestClient(crdRestClientConfig)
	if err != nil {
		return err
	}

	source := cache.NewListWatchFromClient(
		crdClient,
		crd.Plural,
		i.config.WatchNamespace,
		fields.Everything(),
	)
	indexer, controller := cache.NewIndexerInformer(
		source,
		crd.Obj,
		i.config.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    i.config.Handlers.AddFunc,
			UpdateFunc: i.config.Handlers.UpdateFunc,
			DeleteFunc: i.config.Handlers.DeleteFunc,
		},
		cache.Indexers{},
	)
	i.Indexer = indexer

	go controller.Run(ctx.Done())

	return nil
}

func checkIndexedOperatorConfig(config *IndexedOperatorConfig) error {
	if config.Logger == nil {
		config.Logger = log.WithFields(log.Fields{
			"app": "indexed-operator",
		})
	}
	if config.WatchNamespace == "" {
		config.WatchNamespace = DefaultNamespace
	}

	return nil
}
