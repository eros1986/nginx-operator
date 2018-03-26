package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type PVData struct {
	Name        string
	Labels      map[string]string
	Annotations map[string]string

	Spec apiv1.PersistentVolumeSpec
}

type PVInterface interface {
	MakeConfig(*PVData) *apiv1.PersistentVolume
	Create(*apiv1.PersistentVolume) (*apiv1.PersistentVolume, error)
	Get(string) (*apiv1.PersistentVolume, error)
	Delete(string, *metav1.DeleteOptions) error
}

type pvs struct {
	client corev1.PersistentVolumeInterface
}

func NewPVs(kclient *kubernetes.Clientset) PVInterface {
	return &pvs{
		client: kclient.CoreV1().PersistentVolumes(),
	}
}

func (p *pvs) MakeConfig(rawData *PVData) *apiv1.PersistentVolume {
	return &apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:        rawData.Name,
			Labels:      rawData.Labels,
			Annotations: rawData.Annotations,
		},
		Spec: rawData.Spec,
	}
}

func (p *pvs) Create(config *apiv1.PersistentVolume) (*apiv1.PersistentVolume, error) {
	return p.client.Create(config)
}

func (p *pvs) Get(name string) (*apiv1.PersistentVolume, error) {
	return p.client.Get(name, metav1.GetOptions{})
}

func (p *pvs) Delete(name string, options *metav1.DeleteOptions) error {
	return p.client.Delete(name, options)
}
