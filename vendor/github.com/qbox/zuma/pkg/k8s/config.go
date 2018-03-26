package k8s

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func BuildKuberentesConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		// run out-of-cluster
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	// run in-cluster
	return rest.InClusterConfig()
}
