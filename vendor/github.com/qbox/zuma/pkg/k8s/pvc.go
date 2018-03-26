package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type PVCData struct {
	Name   string
	Labels map[string]string

	Spec apiv1.PersistentVolumeClaimSpec
}

type PVCInterface interface {
	MakeConfig(*PVCData) *apiv1.PersistentVolumeClaim
	Create(*apiv1.PersistentVolumeClaim) (*apiv1.PersistentVolumeClaim, error)
	Get(string) (*apiv1.PersistentVolumeClaim, error)
	Delete(string, *metav1.DeleteOptions) error
}

type pvcs struct {
	client corev1.PersistentVolumeClaimInterface
	ns     string
}

func NewPVCs(kclient *kubernetes.Clientset, ns string) PVCInterface {
	return &pvcs{
		client: kclient.CoreV1().PersistentVolumeClaims(ns),
		ns:     ns,
	}
}

func (p *pvcs) MakeConfig(rawData *PVCData) *apiv1.PersistentVolumeClaim {
	return &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rawData.Name,
			Namespace: p.ns,
			Labels:    rawData.Labels,
		},
		Spec: rawData.Spec,
	}
}

func (p *pvcs) Create(config *apiv1.PersistentVolumeClaim) (*apiv1.PersistentVolumeClaim, error) {
	return p.client.Create(config)
}

func (p *pvcs) Get(name string) (*apiv1.PersistentVolumeClaim, error) {
	return p.client.Get(name, metav1.GetOptions{})
}

func (p *pvcs) Delete(name string, options *metav1.DeleteOptions) error {
	return p.client.Delete(name, options)
}
