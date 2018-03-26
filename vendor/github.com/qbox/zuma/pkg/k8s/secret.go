package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type SecreteData struct {
	Name string

	Data map[string]string
}

type SecretInterface interface {
	MakeConfig(*SecreteData) *apiv1.Secret
	Create(*apiv1.Secret) (*apiv1.Secret, error)
	Delete(name string) error
	Get(name string) (*apiv1.Secret, error)
}

type secrets struct {
	client corev1.SecretInterface
	ns     string
}

func NewSecrets(kclient *kubernetes.Clientset, ns string) SecretInterface {
	return &secrets{
		client: kclient.CoreV1().Secrets(ns),
		ns:     ns,
	}
}

func (s *secrets) MakeConfig(rawData *SecreteData) *apiv1.Secret {
	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rawData.Name,
			Namespace: s.ns,
		},
		StringData: rawData.Data,
	}
}

func (s *secrets) Create(config *apiv1.Secret) (*apiv1.Secret, error) {
	return s.client.Create(config)
}

func (s *secrets) Delete(name string) error {
	return s.client.Delete(name, nil)
}

func (s *secrets) Get(name string) (*apiv1.Secret, error) {
	return s.client.Get(name, metav1.GetOptions{})
}
