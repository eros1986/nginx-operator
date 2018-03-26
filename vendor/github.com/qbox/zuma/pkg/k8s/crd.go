package k8s

import (
	"fmt"
	"time"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
)

type CRDData struct {
	Name string

	Spec apiextensionsv1beta1.CustomResourceDefinitionSpec
}

type CRDRestClientConfig struct {
	KubernetesConfig *rest.Config

	// Note: The prefix of the `CRDGroup` should be different among the CRD
	// resources according to the experience.
	CRDGroup   string
	CRDVersion string
	CRDObj     runtime.Object
	CRDObjList runtime.Object

	SchemeBuilder func(*runtime.Scheme) error
}

type CRDInstanceData struct {
	CRDPlural    string
	CRDNamespace string

	CRDObj    runtime.Object
	CRDReturn runtime.Object

	InstanceName string
}

type CRDInterface interface {
	MakeConfig(*CRDData) *apiextensionsv1beta1.CustomResourceDefinition
	Create(*apiextensionsv1beta1.CustomResourceDefinition) (*apiextensionsv1beta1.CustomResourceDefinition, error)
	Get(name string) (*apiextensionsv1beta1.CustomResourceDefinition, error)
	Delete(name string, options *metav1.DeleteOptions) error

	NewRestClient(*CRDRestClientConfig) (*rest.RESTClient, *runtime.Scheme, error)
	CreateInstance(*rest.RESTClient, *CRDInstanceData) error
	GetInstance(*rest.RESTClient, *CRDInstanceData) error
	DeleteInstance(*rest.RESTClient, *CRDInstanceData) error
}

type crds struct {
	client v1beta1.CustomResourceDefinitionInterface
}

func NewCRDs(clientset apiextensionsclient.Interface) CRDInterface {
	return &crds{
		client: clientset.ApiextensionsV1beta1().CustomResourceDefinitions(),
	}
}

func (c *crds) MakeConfig(rawData *CRDData) *apiextensionsv1beta1.CustomResourceDefinition {
	return &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: rawData.Name,
		},
		Spec: rawData.Spec,
	}
}

func (c *crds) Create(config *apiextensionsv1beta1.CustomResourceDefinition) (
	*apiextensionsv1beta1.CustomResourceDefinition,
	error) {
	crd, err := c.client.Create(config)
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			return c.client.Get(config.ObjectMeta.Name, metav1.GetOptions{})
		}

		return nil, err
	}

	// wait for CRD being established
	err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		crd, err = c.client.Get(config.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, cond := range crd.Status.Conditions {
			switch cond.Type {
			case apiextensionsv1beta1.Established:
				if cond.Status == apiextensionsv1beta1.ConditionTrue {
					return true, err
				}
			case apiextensionsv1beta1.NamesAccepted:
				if cond.Status == apiextensionsv1beta1.ConditionFalse {
					fmt.Printf("Name conflict: %v\n", cond.Reason)
				}
			}
		}

		return false, err
	})
	if err != nil {
		if deleteErr := c.client.Delete(config.ObjectMeta.Name, nil); deleteErr != nil {
			return nil, errors.NewAggregate([]error{err, deleteErr})
		}

		return nil, err
	}

	return crd, nil
}

func (c *crds) Get(name string) (*apiextensionsv1beta1.CustomResourceDefinition, error) {
	return c.client.Get(name, metav1.GetOptions{})
}

func (c *crds) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete(name, options)
}

func (c *crds) NewRestClient(config *CRDRestClientConfig) (*rest.RESTClient, *runtime.Scheme, error) {
	schemeBuilder := runtime.NewSchemeBuilder(config.SchemeBuilder)
	addToScheme := schemeBuilder.AddToScheme

	scheme := runtime.NewScheme()
	if err := addToScheme(scheme); err != nil {
		return nil, nil, err
	}

	cfg := *config.KubernetesConfig
	cfg.GroupVersion = &schema.GroupVersion{
		Group:   config.CRDGroup,
		Version: config.CRDVersion,
	}
	cfg.APIPath = "/apis"
	cfg.ContentType = runtime.ContentTypeJSON
	cfg.NegotiatedSerializer = serializer.DirectCodecFactory{
		CodecFactory: serializer.NewCodecFactory(scheme),
	}

	client, err := rest.RESTClientFor(&cfg)
	if err != nil {
		return nil, nil, err
	}

	return client, scheme, nil
}

func (c *crds) CreateInstance(client *rest.RESTClient, instance *CRDInstanceData) error {
	return client.Post().
		Resource(instance.CRDPlural).
		Namespace(instance.CRDNamespace).
		Body(instance.CRDObj).
		Do().
		Into(instance.CRDReturn)
}

func (c *crds) GetInstance(client *rest.RESTClient, instance *CRDInstanceData) error {
	return client.Get().
		Resource(instance.CRDPlural).
		Namespace(instance.CRDNamespace).
		Name(instance.InstanceName).
		Do().
		Into(instance.CRDReturn)
}

func (c *crds) DeleteInstance(client *rest.RESTClient, instance *CRDInstanceData) error {
	return client.Delete().
		Resource(instance.CRDPlural).
		Namespace(instance.CRDNamespace).
		Name(instance.InstanceName).
		Do().
		Error()
}
