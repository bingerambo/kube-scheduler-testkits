package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	// KubeConfigFile is a string flag which indicates kubeconfig filepath
	KubeConfigFile string

	// Clientset init on start, used by others to create k8s resources
	Clientset *kubernetes.Clientset

	// KubeConfig init on start, used by others to create k8s rest client
	KubeConfig *rest.Config
)

func MustInit() {
	var err error
	KubeConfig, err = clientcmd.BuildConfigFromFlags("", KubeConfigFile)
	if err != nil {
		panic(err)
	}
	Clientset = kubernetes.NewForConfigOrDie(KubeConfig)
}

func GetClient() *kubernetes.Clientset {
	return Clientset
}
