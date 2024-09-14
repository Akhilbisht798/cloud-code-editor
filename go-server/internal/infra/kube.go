package infra

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeHandler struct {
	Client *kubernetes.Clientset
}

func int32ptr(i int32) *int32 {
	return &i
}

// TODO: change it according to production
func GetClient() (KubeHandler, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error: getting home directory ", err)
		panic(err.Error())
	}

	kubeConfig := flag.String("kubeconfig", filepath.Join(homeDir, ".kube", "config"), "")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	client := KubeHandler{
		Client: clientset,
	}
	return client, nil
}

func (c *KubeHandler) CreateContainerAndService(userId string, projectId string) (string, error) {
	label := map[string]string{
		"app":       "socket-server",
		"userID":    userId,
		"projectID": projectId,
	}

	deploymentClient := c.Client.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("user-%s-project-%s", userId, projectId),
			Labels: label,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: label,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: label,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "socket-server",
							Image: "akhilbisht798/socket-server",
							Env: []apiv1.EnvVar{
								{
									Name:  "userId",
									Value: userId,
								},
								{
									Name:  "projectId",
									Value: projectId,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("user-%s-project-%s", userId, projectId),
			Labels: label,
		},
		Spec: apiv1.ServiceSpec{
			Selector: label,
			Ports: []apiv1.ServicePort{
				{
					Port:       5000,
					TargetPort: intstr.FromInt(5000),
				},
			},
			Type: apiv1.ServiceTypeLoadBalancer,
		},
	}
	serviceClient := c.Client.CoreV1().Services(apiv1.NamespaceDefault)
	_, err = serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	//Getting the ip back
	serviceName := fmt.Sprintf("user-%s-project-%s", userId, projectId)
	svc, err := serviceClient.Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return svc.Spec.ClusterIP, nil
}

// TODO:Proper handling of error
func (c *KubeHandler) CloseContainerAndService(userId string, projectId string) error {
	labels := fmt.Sprintf("user-%s-project-%s", userId, projectId)

	deploymentClient := c.Client.AppsV1().Deployments(apiv1.NamespaceDefault)
	err := deploymentClient.Delete(context.TODO(), labels, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	serviceClient := c.Client.CoreV1().Services(apiv1.NamespaceDefault)
	err = serviceClient.Delete(context.TODO(), labels, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
