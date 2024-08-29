package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"golang.org/x/crypto/bcrypt"

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

type KubeHandlerRequest struct {
	UserId    string `json:"userId"`
	ProjectId string `json:"projectId"`
}

var clientset *kubernetes.Clientset

const secretKey = "secret"

func kubeHandler(w http.ResponseWriter, r *http.Request) {
	var req KubeHandlerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if clientset == nil {
		clientset, err = getClient()
		if err != nil {
			panic(err.Error())
		}
	}
	label := map[string]string{
		"app":       "socket-server",
		"userID":    req.UserId,
		"projectID": req.ProjectId,
	}

	deploymentClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("user-%s-project-%s", req.UserId, req.ProjectId),
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
									Value: req.UserId,
								},
								{
									Name:  "projectId",
									Value: req.ProjectId,
								},
							},
						},
					},
				},
			},
		},
	}

	deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("user-%s-project-%s", req.UserId, req.ProjectId),
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

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)
	user := User{
		Email:    data["email"],
		Password: string(password),
	}
	DB.Create(&user)
	jsonData, _ := json.Marshal(user)
	w.Write(jsonData)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	var user User

	DB.Where("email = ?", data["email"]).First(&user)
	if user.Id == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"]))
	if err != nil {
		http.Error(w, "wrong password", http.StatusBadRequest)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer: strconv.Itoa(user.Id),
		ExpiresAt: &jwt.Time{
			Time: time.Now().Add(time.Hour * 24),
		},
	})

	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		http.Error(w, "Could'nt login", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user User
	DB.Where("id = ?", claims.Issuer).First(&user)

	jsonData, err := json.Marshal(&user)
	w.Write([]byte(jsonData))
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "/",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

func S3PresignedGetURLHandler(w http.ResponseWriter, r *http.Request) {
	var req KubeHandlerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// prefix := fmt.Sprintf("userId-%s/projectId-%s", req.UserId, req.ProjectId)
	prefix := fmt.Sprintf("userId-%s/client", req.UserId)
	fmt.Printf("prefix: %s\n", prefix)

	client := getS3ClientDevelopment()
	getPresignerClient()
	urls, err := client.GetPresignedUrls("project", prefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}
