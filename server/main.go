package main

import (
	"fmt"
	"net/http"
)

func int32ptr(i int32) *int32 {
	return &i
}

func main() {
	// dbConnect()
	fmt.Println("Server Starting at Port 3000")
	http.HandleFunc("/api", rootHandler)
	http.HandleFunc("/api/kube", kubeHandler)
	http.HandleFunc("/api/kube/delete", closeResource)

	http.HandleFunc("/api/register", register)
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/user", getUser)
	http.HandleFunc("/api/logout", logout)

	http.HandleFunc("/api/getUserFiles", S3PresignedGetURLHandler)
	http.ListenAndServe(":3000", nil)
}
