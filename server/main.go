package main

import (
	"fmt"
	"net/http"
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
