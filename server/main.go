package main

import (
	"fmt"
	"net/http"
)

func int32ptr(i int32) *int32 {
	return &i
}

func main() {
	dbConnect()
	fmt.Println("Server Starting at Port 3000")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/kube", kubeHandler)
	http.HandleFunc("/delete", closeResource)

	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/user", getUser)
	http.ListenAndServe(":3000", nil)
}
