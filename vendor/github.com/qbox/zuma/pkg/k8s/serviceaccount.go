package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type ServiceAccountData struct {
	Name                         string
	AutomountServiceAccountToken bool
}

type ServiceAccountInterface interface {
	MakeConfig(*ServiceAccountData) *apiv1.ServiceAccount
	Create(*apiv1.ServiceAccount) (*apiv1.ServiceAccount, error)
	Get(string) (*apiv1.ServiceAccount, error)
	Delete(string, *metav1.DeleteOptions) error
}

type serviceaccounts struct {
	client corev1.ServiceAccountInterface
}

func NewServiceAccounts(kclient *kubernetes.Clientset, ns string) ServiceAccountInterface {
	return &serviceaccounts{
		client: kclient.CoreV1().ServiceAccounts(ns),
	}
}

func (s *serviceaccounts) MakeConfig(rawData *ServiceAccountData) *apiv1.ServiceAccount {
	return &apiv1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: rawData.Name,
		},
		AutomountServiceAccountToken: &(rawData.AutomountServiceAccountToken),
	}
}

func (s *serviceaccounts) Create(config *apiv1.ServiceAccount) (*apiv1.ServiceAccount, error) {
	return s.client.Create(config)
}

func (s *serviceaccounts) Get(name string) (*apiv1.ServiceAccount, error) {
	return s.client.Get(name, metav1.GetOptions{})
}

func (s *serviceaccounts) Delete(name string, options *metav1.DeleteOptions) error {
	return s.client.Delete(name, options)
}
