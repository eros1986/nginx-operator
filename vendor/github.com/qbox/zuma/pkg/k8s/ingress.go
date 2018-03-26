package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	extensionsv1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type IngressData struct {
	Name string
	Spec extensionsv1beta1.IngressSpec
}

type IngressInterface interface {
	MakeConfig(*IngressData) *extensionsv1beta1.Ingress
	Create(*extensionsv1beta1.Ingress) (*extensionsv1beta1.Ingress, error)
	Get(string) (*extensionsv1beta1.Ingress, error)
	Delete(string, *metav1.DeleteOptions) error
}

type ingresses struct {
	client v1beta1.IngressInterface
}

func NewIngresses(kclient *kubernetes.Clientset, ns string) IngressInterface {
	return &ingresses{
		client: kclient.ExtensionsV1beta1().Ingresses(ns),
	}
}

func (i *ingresses) MakeConfig(rawData *IngressData) *extensionsv1beta1.Ingress {
	return &extensionsv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: rawData.Name,
		},
		Spec: rawData.Spec,
	}
}

func (i *ingresses) Create(config *extensionsv1beta1.Ingress) (*extensionsv1beta1.Ingress, error) {
	return i.client.Create(config)
}

func (i *ingresses) Get(name string) (*extensionsv1beta1.Ingress, error) {
	return i.client.Get(name, metav1.GetOptions{})
}

func (i *ingresses) Delete(name string, options *metav1.DeleteOptions) error {
	return i.client.Delete(name, options)
}
