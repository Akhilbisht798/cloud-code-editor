package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func int32ptr(i int32) *int32 {
	return &i
}

func main() {
	fmt.Println("Server Starting at Port 3000")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/kube", kubeHandler)
	http.ListenAndServe(":3000", nil)
}

func kubeHandler(w http.ResponseWriter, r *http.Request) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error: getting home directory ", err)
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

	//TODO: Later user request to get these.
	userId := "8"
	projectId := "8"
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
	// if err != nil {
	// 	log.Fatal("failed to create deployment: %v", err)
	// } else {
	// 	fmt.Println("deployment created successfully!")
	// }

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
					NodePort:   30000,
				},
			},
			Type: apiv1.ServiceTypeLoadBalancer,
		},
	}
	serviceClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})

	serviceClient.Get(context.TODO(), "", metav1.GetOptions{})

	// if err != nil {
	// 	log.Fatal("failed to create service: %v", err)
	// } else {
	// 	fmt.Println("service created successfully!")
	// }

	w.Write([]byte("Hello world"))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
