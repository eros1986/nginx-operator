package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type NamespaceInterface interface {
	MakeConfig(name string) *apiv1.Namespace
	Create(*apiv1.Namespace) (*apiv1.Namespace, error)
	Get(name string) (*apiv1.Namespace, error)
	Delete(name string) error
}

type namesapce struct {
	client corev1.NamespaceInterface
}

func NewNamespace(kclient *kubernetes.Clientset) NamespaceInterface {
	return &namesapce{
		client: kclient.CoreV1().Namespaces(),
	}
}

func (n *namesapce) MakeConfig(name string) *apiv1.Namespace {
	return &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func (n *namesapce) Create(config *apiv1.Namespace) (*apiv1.Namespace, error) {
	return n.client.Create(config)
}

func (n *namesapce) Get(name string) (*apiv1.Namespace, error) {
	return n.client.Get(name, metav1.GetOptions{})
}

func (n *namesapce) Delete(name string) error {
	return n.client.Delete(name, nil)
}
