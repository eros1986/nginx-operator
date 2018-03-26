package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
)

type StatefulSetData struct {
	Name string

	Spec v1beta1.StatefulSetSpec
}

type StatefulSetInterface interface {
	MakeConfig(*StatefulSetData) *v1beta1.StatefulSet
	Create(*v1beta1.StatefulSet) (*v1beta1.StatefulSet, error)
	Delete(string, *metav1.DeleteOptions) error
	Get(name string) (*v1beta1.StatefulSet, error)
}

type statefulSets struct {
	client appsv1beta1.StatefulSetInterface
	ns     string
}

func NewStatefulSets(kclient *kubernetes.Clientset, ns string) StatefulSetInterface {
	return &statefulSets{
		client: kclient.Apps().StatefulSets(ns),
		ns:     ns,
	}
}

func (s *statefulSets) MakeConfig(rawData *StatefulSetData) *v1beta1.StatefulSet {
	return &v1beta1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: rawData.Name,
		},
		Spec: rawData.Spec,
	}
}

func (s *statefulSets) Create(config *v1beta1.StatefulSet) (*v1beta1.StatefulSet, error) {
	return s.client.Create(config)
}

func (s *statefulSets) Delete(name string, options *metav1.DeleteOptions) error {
	return s.client.Delete(name, options)
}

func (s *statefulSets) Get(name string) (*v1beta1.StatefulSet, error) {
	return s.client.Get(name, metav1.GetOptions{})
}
