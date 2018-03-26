package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type ConfigMapData struct {
	ConfigMapName string

	Data map[string]string
}

type ConfigMapInterface interface {
	MakeConfig(*ConfigMapData) *apiv1.ConfigMap
	Create(*apiv1.ConfigMap) (*apiv1.ConfigMap, error)
	Update(*apiv1.ConfigMap) (*apiv1.ConfigMap, error)
	Delete(name string) error
	Get(name string) (*apiv1.ConfigMap, error)
}

type configMaps struct {
	client corev1.ConfigMapInterface
	ns     string
}

func NewConfigMaps(kclient *kubernetes.Clientset, ns string) ConfigMapInterface {
	client := kclient.CoreV1().ConfigMaps(ns)
	return &configMaps{
		client: client,
		ns:     ns,
	}
}

func (c *configMaps) MakeConfig(rawData *ConfigMapData) *apiv1.ConfigMap {
	return &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rawData.ConfigMapName,
			Namespace: c.ns,
		},
		Data: rawData.Data,
	}
}

func (c *configMaps) Create(config *apiv1.ConfigMap) (*apiv1.ConfigMap, error) {
	return c.client.Create(config)
}

func (c *configMaps) Update(config *apiv1.ConfigMap) (*apiv1.ConfigMap, error) {
	return c.client.Update(config)
}

func (c *configMaps) Delete(name string) error {
	return c.client.Delete(name, nil)
}

func (c *configMaps) Get(name string) (*apiv1.ConfigMap, error) {
	return c.client.Get(name, metav1.GetOptions{})
}
