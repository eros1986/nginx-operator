package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1beta1 "k8s.io/client-go/kubernetes/typed/rbac/v1beta1"
	"k8s.io/client-go/pkg/apis/rbac/v1beta1"
)

type ClusterRoleData struct {
	Name  string
	Rules []v1beta1.PolicyRule
}

type ClusterRoleInterface interface {
	MakeConfig(*ClusterRoleData) *v1beta1.ClusterRole
	Create(*v1beta1.ClusterRole) (*v1beta1.ClusterRole, error)
	Get(string) (*v1beta1.ClusterRole, error)
	Delete(string, *metav1.DeleteOptions) error
}

type clusterroles struct {
	client apiv1beta1.ClusterRoleInterface
}

func NewClusterRoles(kclient *kubernetes.Clientset) ClusterRoleInterface {
	return &clusterroles{
		client: kclient.RbacV1beta1().ClusterRoles(),
	}
}

func (c *clusterroles) MakeConfig(rawData *ClusterRoleData) *v1beta1.ClusterRole {
	return &v1beta1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: rawData.Name,
		},
		Rules: rawData.Rules,
	}
}

func (c *clusterroles) Create(config *v1beta1.ClusterRole) (*v1beta1.ClusterRole, error) {
	return c.client.Create(config)
}

func (c *clusterroles) Get(name string) (*v1beta1.ClusterRole, error) {
	return c.client.Get(name, metav1.GetOptions{})
}

func (c *clusterroles) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete(name, options)
}
