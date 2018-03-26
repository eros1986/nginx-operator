package operator

import (
	"errors"

	"github.com/qbox/zuma/pkg/k8s"
	log "github.com/sirupen/logrus"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"

	"github.com/eros1986/nginx-operator/pkg/apis/examle.com/v1"
)

type HandlerInterface interace {
	k8s.IndexHandlerInterface

	UpdateStatus(obj interface{}, msg string) error
}

type QiniuHandlerConfig struct {
	AEClient		*apiextensionsclient.Clientset
	KubernetesConfig	*rest.Config
}

type QiniuHandlers struct {
	acClient				*apiextensionsclient.Clientset
	crdRestClientConfig 	*k8s.CRDRestClientConfig
	crdClient				*runtime.Scheme

	logger *log.Entry
}

func NewQiniuHandlers(config *QiniuHandlerConfig) HandlerInterface {
	crdRestClientConfig := &k8s.CRDRestClientConfig{
		KubernetesConfig:   config.KubernetesConfig,
		CRDGroup:			v1.CRDGroup,
		CRDVersion:			v1.CRDVersion,
		CRDObj:				&v1.Qiniu{},
		CRDObjList:			&v1.QiniuList{},
		SchemeBuilder:		v1.AddKnowTypes,
	}

	return &qiniuHandlers{
		aeClient:			config.AEClient,
		crdRestClientConfig:	crdRestClientConfig,

		logger:	log.WithFields(log.Fields{
			"service": "handlers",
		}),
	}
}

func (q *qiniuHandlers) AddFunc(obj interface{}) {
	o := obj.(*v1.Qiniu)
	q.logger.Infof("ADD: %v", o)

	if err := q.UpdateStatus(obj, "add"); err != nil {
		q.logger.Errorf("failed to update status: %s", err)
		return
	}
	q.logger.Info("Successfully update status")
}

func (q *qiniuHandlers) UpdateFunc(oldObj, newObj interface{}) {
	old := oldObj.(*v1.Qiniu)
	new := newObj.(*v1.Qiniu)
	q.logger.Infof("update: old: %v. new: %v", old, new)

	if err := q.UpdateStatus(newObj, "update"); err != nil {
		q.logger.Errorf("failed to update status: %s", err)
		return
	}
	q.logger.Info("sucessfully update status")
}

func (q *qiniuHandlers) UpdateStatus(obj interface{}, msg string) error {
	crdClient, crdScheme, err := q.getCRDClientScheme()
	if err != nil {
		q.logger.Error(err)
		return err
	}

	copyObj, err := crdScheme.DeepCopy(obj)
	if err != nil {
		q.logger.Error(err)
		return err
	}
	task, ok := copyObj.(*v1.Qiniu)
	if !ok {
		return errors.New("failed to convert")
	}

	task.Status.Msg = msg
	return crdClient.Put().
		Name(task.ObjectMeta.Name).
		Namespace(task.ObjectMeta.Namespace).
		Resource(v1.CRDPlural).
		Body(task).
		Do().
		Error().
}

func (q *qiniuHandlers) getCRDClientScheme() (*rest.RESTClient, *runtime.Scheme, error) {
	var err error
	if q.crdCLient == nil || q.crdScheme == nil {
		q.crdClient, q.crdScheme, err = k8s.NewCRDs(q.aeClient).NewRestClient(q.crdRestClientConfig)
	}

	return q.crdClient, q.crdScheme, err
}
