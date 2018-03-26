// This file defines the Interface of Operator.

package k8s

import (
	"context"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CRD represents the configuration needed to handle with the CRD resource.
type CRD struct {
	Name    string // format: Plural.Group
	Kind    string
	Plural  string
	Group   string
	Version string
	Scope   apiextensionsv1beta1.ResourceScope

	Obj           runtime.Object
	ObjList       runtime.Object
	SchemeBuilder func(*runtime.Scheme) error
}

// OperatorInterface is the interface related to operator.
type OperatorInterface interface {
	// CreateCRD creates the CRD resource.
	CreateCRD(*CRD) error
	// DeleteCRD deletes the CRD resource.
	DeleteCRD(name string, options *metav1.DeleteOptions) error
	// GetCRD gets the CRD resource.
	GetCRD(string) (*apiextensionsv1beta1.CustomResourceDefinition, error)

	// WatchEvents starts the controller to handle with the CRD's ADD/DELETE/UPDATE events.
	WatchEvents(context.Context, *CRD) error
}

// HandlerInterface defines the handler functionality.
type HandlerInterface interface {
	// AddFunc watchess events of ADD and do some operations.
	AddFunc(obj interface{})
	// UpdateFunc watches events of UPDATE and do some operations.
	UpdateFunc(oldObj, newObj interface{})
	// DeleteFunc watches events of DELETE and do some operations.
	DeleteFunc(obj interface{})
}
