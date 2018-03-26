package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	storagev1 "k8s.io/client-go/kubernetes/typed/storage/v1"
	storagetypev1 "k8s.io/client-go/pkg/apis/storage/v1"
)

type StorageClassData struct {
	Name   string
	Labels map[string]string

	Provisioner string

	Parameters map[string]string
}

type StorageClassInterface interface {
	MakeConfig(*StorageClassData) *storagetypev1.StorageClass
	Create(*storagetypev1.StorageClass) (*storagetypev1.StorageClass, error)
	Get(string) (*storagetypev1.StorageClass, error)
	Delete(string, *metav1.DeleteOptions) error
}

type storageClasses struct {
	client storagev1.StorageClassInterface
}

func NewStorageClasses(kclient *kubernetes.Clientset) StorageClassInterface {
	return &storageClasses{
		client: kclient.StorageV1().StorageClasses(),
	}
}

func (s *storageClasses) MakeConfig(rawData *StorageClassData) *storagetypev1.StorageClass {
	return &storagetypev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:   rawData.Name,
			Labels: rawData.Labels,
		},
		Provisioner: rawData.Provisioner,
		Parameters:  rawData.Parameters,
	}
}

func (s *storageClasses) Create(config *storagetypev1.StorageClass) (*storagetypev1.StorageClass, error) {
	return s.client.Create(config)
}

func (s *storageClasses) Get(name string) (*storagetypev1.StorageClass, error) {
	return s.client.Get(name, metav1.GetOptions{})
}

func (s *storageClasses) Delete(name string, options *metav1.DeleteOptions) error {
	return s.client.Delete(name, options)
}
