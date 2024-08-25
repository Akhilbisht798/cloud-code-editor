package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

type kubeHandlerResponse struct {
	Status int    `json:"status"`
	IP     string `json:"ip"`
}

var clientset *kubernetes.Clientset

func kubeHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	if clientset == nil {
		clientset, err = getClient()
		if err != nil {
			panic(err.Error())
		}
	}
	//TODO: Later user request to get these.
	userId := "1"
	projectId := "1"
	label := map[string]string{
		"app":       "socket-server",
		"userID":    userId,
		"projectID": projectId,
	}

	deploymentClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

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
						},
					},
				},
			},
		},
	}

	deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})

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
	serviceClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})

	serviceName := "user-1-project-1"
	svc, err := serviceClient.Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("error: ", err)
	}
	resp := kubeHandlerResponse{
		Status: 200,
		IP:     svc.Spec.ClusterIP,
	}
	respData, err := json.Marshal(resp)
	if err != nil {
		panic(err.Error())
	}

	w.Write(respData)
}

// TODO: save the file to s3 before clearing the resoure.
func closeResource(w http.ResponseWriter, r *http.Request) {
	var err error
	if clientset == nil {
		clientset, err = getClient()
		if err != nil {
			panic(err.Error())
		}
	}

	labels := fmt.Sprintf("user-1-project-1")
	deploymentClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deploymentClient.Delete(context.TODO(), labels, metav1.DeleteOptions{})

	serviceClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	serviceClient.Delete(context.TODO(), labels, metav1.DeleteOptions{})

	w.Write([]byte("deleted the thing"))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
