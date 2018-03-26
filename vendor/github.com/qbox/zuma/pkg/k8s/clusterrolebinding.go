package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	"k8s.io/client-go/pkg/apis/rbac/v1beta1"
)

type ClusterRoleBindingData struct {
	Name string

	Subjects []v1beta1.Subject
	RoleRef  v1beta1.RoleRef
}

type ClusterRoleBindingInterface interface {
	MakeConfig(*ClusterRoleBindingData) *v1beta1.ClusterRoleBinding
	Create(*v1beta1.ClusterRoleBinding) (*v1beta1.ClusterRoleBinding, error)
	Get(string) (*v1beta1.ClusterRoleBinding, error)
	Delete(string, *metav1.DeleteOptions) error
}

type clusterrolebindings struct {
	client apiv1beta1.ClusterRoleBindingInterface
}

func NewClusterRoleBindings(kclient *kubernetes.Clientset) ClusterRoleBindingInterface {
	return &clusterrolebindings{
		client: kclient.RbacV1beta1().ClusterRoleBindings(),
	}
}

func (c *clusterrolebindings) MakeConfig(rawData *ClusterRoleBindingData) *v1beta1.ClusterRoleBinding {
	return &v1beta1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: rawData.Name,
		},
		Subjects: rawData.Subjects,
		RoleRef:  rawData.RoleRef,
	}
}

func (c *clusterrolebindings) Create(config *v1beta1.ClusterRoleBinding) (*v1beta1.ClusterRoleBinding, error) {
	return c.client.Create(config)
}

func (c *clusterrolebindings) Get(name string) (*v1beta1.ClusterRoleBinding, error) {
	return c.client.Get(name, metav1.GetOptions{})
}

func (c *clusterrolebindings) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete(name, options)
}
