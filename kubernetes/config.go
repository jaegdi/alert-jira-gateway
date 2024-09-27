package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func GetConfigMap(clientset *kubernetes.Clientset, namespace, name string) (*v1.ConfigMap, error) {
	return clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func GetSecret(clientset *kubernetes.Clientset, namespace, name string) (*v1.Secret, error) {
	return clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
