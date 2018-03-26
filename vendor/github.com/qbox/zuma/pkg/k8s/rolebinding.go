package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	"k8s.io/client-go/pkg/apis/rbac/v1beta1"
)

type RoleBindingData struct {
	Name string

	Subjects []v1beta1.Subject
	RoleRef  v1beta1.RoleRef
}

type RoleBindingInterface interface {
	MakeConfig(*RoleBindingData) *v1beta1.RoleBinding
	Create(*v1beta1.RoleBinding) (*v1beta1.RoleBinding, error)
	Get(string) (*v1beta1.RoleBinding, error)
	Delete(string, *metav1.DeleteOptions) error
}

type rolebindings struct {
	client apiv1beta1.RoleBindingInterface
}

func NewRoleBindings(kclient *kubernetes.Clientset, ns string) RoleBindingInterface {
	return &rolebindings{
		client: kclient.RbacV1beta1().RoleBindings(ns),
	}
}

func (r *rolebindings) MakeConfig(rawData *RoleBindingData) *v1beta1.RoleBinding {
	return &v1beta1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: rawData.Name,
		},
		Subjects: rawData.Subjects,
		RoleRef:  rawData.RoleRef,
	}
}

func (r *rolebindings) Create(config *v1beta1.RoleBinding) (*v1beta1.RoleBinding, error) {
	return r.client.Create(config)
}

func (r *rolebindings) Get(name string) (*v1beta1.RoleBinding, error) {
	return r.client.Get(name, metav1.GetOptions{})
}

func (r *rolebindings) Delete(name string, options *metav1.DeleteOptions) error {
	return r.client.Delete(name, options)
}
