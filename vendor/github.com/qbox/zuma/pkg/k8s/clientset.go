package k8s

import (
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ClientInterface interface {
	CreateKClientset() (*kubernetes.Clientset, error)
	CreateAPIExtensionClientset() (*apiextensionsclient.Clientset, error)

	GetKClientset() (*kubernetes.Clientset, error)
	GetAPIExtensionClientset() (*apiextensionsclient.Clientset, error)
}

type clientSets struct {
	config *rest.Config

	kclient  *kubernetes.Clientset
	aeClient *apiextensionsclient.Clientset
}

func NewClientSets(config *rest.Config) ClientInterface {
	return &clientSets{
		config: config,
	}
}

func (c *clientSets) CreateKClientset() (*kubernetes.Clientset, error) {
	kclient, err := kubernetes.NewForConfig(c.config)
	if err != nil {
		return nil, err
	}
	c.kclient = kclient

	return kclient, nil
}

func (c *clientSets) CreateAPIExtensionClientset() (*apiextensionsclient.Clientset, error) {
	aeClient, err := apiextensionsclient.NewForConfig(c.config)
	if err != nil {
		return nil, err
	}
	c.aeClient = aeClient

	return aeClient, nil
}

func (c *clientSets) GetKClientset() (*kubernetes.Clientset, error) {
	if c.kclient != nil {
		return c.kclient, nil
	}

	return c.CreateKClientset()
}

func (c *clientSets) GetAPIExtensionClientset() (*apiextensionsclient.Clientset, error) {
	if c.aeClient != nil {
		return c.aeClient, nil
	}

	return c.CreateAPIExtensionClientset()
}
