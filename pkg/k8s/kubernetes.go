package k8s

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeClient is a convenience method for creating kubernetes config and client
// for a given kubeconfig
func GetKubeClient(kubeconfig string) (*kubernetes.Clientset, *restclient.Config, error) {

	log.Debugf("getting kube clientset from kubeconfig %v", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get kubernetes config from kubeconfig '%s': %v", kubeconfig, err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get kubernetes client: %s", err)
	}
	return clientset, config, nil
}

// CreateDeployment creates a new deployment in the cluster
func CreateDeployment(kubeconfig, image, name string) error {

	clientset, _, err := GetKubeClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("cannot get clientset: %v", err)
	}

	deploymentsClient := clientset.AppsV1beta1().Deployments(v1.NamespaceDefault)

	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  "ktkd",
						"name": "ktkd",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  name,
							Image: image,
						},
					},
				},
			},
		},
	}

	_, err = deploymentsClient.Create(deployment)
	if err != nil {
		return fmt.Errorf("cannot create deployment: %v", err)
	}

	return nil
}

func int32Ptr(i int32) *int32 { return &i }
