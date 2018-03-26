package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	extensionsv1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type DaemonSetData struct {
	Name   string
	Labels map[string]string

	Spec extensionsv1beta1.DaemonSetSpec
}

type DaemonSetInterface interface {
	MakeConfig(*DaemonSetData) *extensionsv1beta1.DaemonSet
	Create(*extensionsv1beta1.DaemonSet) (*extensionsv1beta1.DaemonSet, error)
	Get(string) (*extensionsv1beta1.DaemonSet, error)
	Delete(string, *metav1.DeleteOptions) error
}

type daemonSets struct {
	client v1beta1.DaemonSetInterface
}

func NewDaemonSets(kclient *kubernetes.Clientset, ns string) DaemonSetInterface {
	return &daemonSets{
		client: kclient.ExtensionsV1beta1().DaemonSets(ns),
	}
}

func (d *daemonSets) MakeConfig(rawData *DaemonSetData) *extensionsv1beta1.DaemonSet {
	return &extensionsv1beta1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:   rawData.Name,
			Labels: rawData.Labels,
		},
		Spec: rawData.Spec,
	}
}

func (d *daemonSets) Create(config *extensionsv1beta1.DaemonSet) (*extensionsv1beta1.DaemonSet, error) {
	return d.client.Create(config)
}

func (d *daemonSets) Get(name string) (*extensionsv1beta1.DaemonSet, error) {
	return d.client.Get(name, metav1.GetOptions{})
}

func (d *daemonSets) Delete(name string, options *metav1.DeleteOptions) error {
	return d.client.Delete(name, options)
}
