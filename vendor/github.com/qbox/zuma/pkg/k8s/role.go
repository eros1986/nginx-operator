package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	"k8s.io/client-go/pkg/apis/rbac/v1beta1"
)

type RoleData struct {
	Name  string
	Rules []v1beta1.PolicyRule
}

type RoleInterface interface {
	MakeConfig(*RoleData) *v1beta1.Role
	Create(*v1beta1.Role) (*v1beta1.Role, error)
	Get(string) (*v1beta1.Role, error)
	Delete(string, *metav1.DeleteOptions) error
}

type roles struct {
	client apiv1beta1.RoleInterface
}

func NewRoles(kclient *kubernetes.Clientset, ns string) RoleInterface {
	return &roles{
		client: kclient.RbacV1beta1().Roles(ns),
	}
}

func (r *roles) MakeConfig(rawData *RoleData) *v1beta1.Role {
	return &v1beta1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: rawData.Name,
		},
		Rules: rawData.Rules,
	}
}

func (r *roles) Create(config *v1beta1.Role) (*v1beta1.Role, error) {
	return r.client.Create(config)
}

func (r *roles) Get(name string) (*v1beta1.Role, error) {
	return r.client.Get(name, metav1.GetOptions{})
}

func (r *roles) Delete(name string, options *metav1.DeleteOptions) error {
	return r.client.Delete(name, options)
}
