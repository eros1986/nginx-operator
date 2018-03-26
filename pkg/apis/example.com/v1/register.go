package v1

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	CRDKind    = "Qiniu"
	CRDPlural  = strings.ToLower(CRDKind) + "s"
	CRDGroup   = "example.com"
	CRDVersion = "v1"
	CRDName    = fmt.Sprintf("%s.%s", CRDPlural, CRDGroup)
)

func AddKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		schema.GroupVersion{
			Group:   CRDGroup,
			Version: CRDVersion,
		},
		&Qiniu{},
		&QiniuList{},
	)
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{
		Group:   CRDGroup,
		Version: CRDVersion,
	})
	return nil
}
