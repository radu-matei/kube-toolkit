package portforwarder

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/radu-matei/kube-toolkit/pkg/k8s"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
)

//New returns a tunnel to the ktkd pod.
func New(clientset *kubernetes.Clientset, config *restclient.Config, namespace string, remotePort, localPort int) (*k8s.Tunnel, error) {
	podName, err := getKTKDPodName(clientset, namespace)
	if err != nil {
		return nil, err
	}
	log.Debugf("found pod: %s", podName)

	t := k8s.NewTunnel(clientset.CoreV1().RESTClient(), config, namespace, podName, remotePort)
	return t, t.ForwardPort(localPort)
}

func getKTKDPodName(clientset *kubernetes.Clientset, namespace string) (string, error) {
	// TODO use a const for labels
	selector := labels.Set{"app": "ktkd", "name": "ktkd"}.AsSelector()
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
		return nil, fmt.Errorf("could not find ktkd")
	}
	for _, p := range pods.Items {
		if podutil.IsPodReady(&p) {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("could not find a ready ktkd pod")
}
