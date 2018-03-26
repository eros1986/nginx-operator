package operator

import (
	"context"
	"time"

	"github.com/qbox/zuma/pkg/k8s"
	log "github.com/sirupsen/logrus"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	"github.com/eros1986/nginx-operator/pkg/apis/example.com/v1"
)

type OperatorInterface interface {
	Run(ctx context.Context, stopc <-chan struct{}) error
}

type OperatorConfig struct {
	KubeConfigPath string
	WatchNamespace string
	ResyncPeriod   time.Duration
}

type operator struct {
	op  k8s.OperatorInterface
	crd *k8s.CRD

	logger *log.Entry
}

func NewOperator(config *OperatorConfig) (OperatorInterface, error) {
	kubernetesConfig, err := k8s.BuildKuberentesConfig(config.KubeConfigPath)
	if err != nil {
		return nil, err
	}

	clientsets := k8s.NewClientSets(kubernetesConfig)
	aeClient, err := clientsets.GetAPIExtensionClientset()
	if err != nil {
		return nil, err
	}

	handlerConfig := &QiniuHandlersConfig{
		AEClient:         aeClient,
		KubernetesConfig: kubernetesConfig,
	}
	handlers := NewQiniuHandlers(handlerConfig)

	indexedOperatorConfig := &k8s.IndexedOperatorConfig{
		KubeConfigPath: config.KubeConfigPath,
		WatchNamespace: "",
		ResyncPeriod:   config.ResyncPeriod,
		Handlers:       handlers,
	}
	op, err := k8s.NewIndexedOperator(indexedOperatorConfig)
	if err != nil {
		return nil, err
	}

	crd := &k8s.CRD{
		Name:    v1.CRDName,
		Kind:    v1.CRDKind,
		Plural:  v1.CRDPlural,
		Group:   v1.CRDGroup,
		Version: v1.CRDVersion,
		Scope:   apiextensionsv1beta1.NamespaceScoped,

		Obj:           &v1.Qiniu{},
		ObjList:       &v1.QiniuList{},
		SchemeBuilder: v1.AddKnownTypes,
	}

	return &operator{
		op:  op,
		crd: crd,
		logger: log.WithFields(log.Fields{
			"app": "operator",
		}),
	}, nil
}

func (o *operator) Run(ctx context.Context, stopc <-chan struct{}) error {
	// Create CRD.
	o.logger.Info("Begin to create crd.")
	if err := o.op.CreateCRD(o.crd); err != nil {
		return err
	}
	o.logger.Info("Successfully create crd.")

	// Watch CRD events.
	o.logger.Info("Begin to watch events. ")
	if err := o.op.WatchEvents(ctx, o.crd); err != nil {
		return err
	}

	<-stopc

	o.logger.Info("Bye.")

	return nil
}
