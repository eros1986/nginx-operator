package k8s

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type PodTemplateSpecData struct {
	Name   string
	Labels map[string]string

	Spec apiv1.PodSpec
}

type PodTemplateSpecInterface interface {
	MakeConfig(*PodTemplateSpecData) *apiv1.PodTemplateSpec
	SetInitContainersAnnotations(*apiv1.PodTemplateSpec) error
}

type podTemplateSpecs struct {
	templateSpec *apiv1.PodTemplateSpec
}

func NewPodTemplateSpecs() PodTemplateSpecInterface {
	return &podTemplateSpecs{}
}

func (p *podTemplateSpecs) MakeConfig(rawData *PodTemplateSpecData) *apiv1.PodTemplateSpec {
	p.templateSpec = &apiv1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:   rawData.Name,
			Labels: rawData.Labels,
		},
		Spec: rawData.Spec,
	}

	return p.templateSpec
}

// TODO: add more annotation fields
func (p *podTemplateSpecs) SetInitContainersAnnotations(spec *apiv1.PodTemplateSpec) error {
	if len(spec.Spec.InitContainers) > 0 {
		value, err := json.Marshal(spec.Spec.InitContainers)
		if err != nil {
			return err
		}
		if spec.Annotations == nil {
			spec.Annotations = make(map[string]string)
		}
		spec.Annotations[apiv1.PodInitContainersBetaAnnotationKey] = string(value)
	}
	return nil
}

type PodData struct {
	Name   string
	Labels map[string]string

	Spec apiv1.PodSpec
}

type PodInterface interface {
	MakeConfig(*PodData) *apiv1.Pod
	Create(*apiv1.Pod) (*apiv1.Pod, error)
	Delete(string, *metav1.DeleteOptions) error
	Get(string) (*apiv1.Pod, error)
	List(metav1.ListOptions) (*apiv1.PodList, error)
}

type pods struct {
	client corev1.PodInterface
	ns     string
}

func NewPods(kclient *kubernetes.Clientset, ns string) PodInterface {
	return &pods{
		client: kclient.CoreV1().Pods(ns),
		ns:     ns,
	}
}

func (p *pods) MakeConfig(rawData *PodData) *apiv1.Pod {
	return &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rawData.Name,
			Namespace: p.ns,
			Labels:    rawData.Labels,
		},
		Spec: rawData.Spec,
	}
}

func (p *pods) Create(config *apiv1.Pod) (*apiv1.Pod, error) {
	return p.client.Create(config)
}

func (p *pods) Delete(name string, options *metav1.DeleteOptions) error {
	return p.client.Delete(name, options)
}

func (p *pods) Get(name string) (*apiv1.Pod, error) {
	return p.client.Get(name, metav1.GetOptions{})
}

func (p *pods) List(options metav1.ListOptions) (*apiv1.PodList, error) {
	return p.client.List(options)
}
