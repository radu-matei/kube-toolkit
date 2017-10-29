package portforwarder

import (
	"fmt"

	"github.com/kubernetes/helm/pkg/kube"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
)

//New returns a tunnel to the Gotham pod.
func New(clientset *kubernetes.Clientset, config *restclient.Config, namespace string) (*kube.Tunnel, error) {
	podName, err := getGothamPodName(clientset, namespace)
	if err != nil {
		return nil, err
	}
	const gothamPort = 10000
	t := kube.NewTunnel(clientset.CoreV1().RESTClient(), config, namespace, podName, gothamPort)
	return t, t.ForwardPort()
}

func getGothamPodName(clientset *kubernetes.Clientset, namespace string) (string, error) {
	// TODO use a const for labels
	selector := labels.Set{"app": "joker", "name": "gotham"}.AsSelector()
	pod, err := getFirstRunningPod(clientset, selector, namespace)
	if err != nil {
		return "", err
	}
	return pod.ObjectMeta.GetName(), nil
}

func getFirstRunningPod(clientset *kubernetes.Clientset, selector labels.Selector, namespace string) (*v1.Pod, error) {
	options := metav1.ListOptions{LabelSelector: selector.String()}
	pods, err := clientset.CoreV1().Pods(namespace).List(options)
	if err != nil {
		return nil, err
	}
	if len(pods.Items) < 1 {
		return nil, fmt.Errorf("could not find gotham")
	}
	for _, p := range pods.Items {
		if podutil.IsPodReady(&p) {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("could not find a ready gotham pod")
}
