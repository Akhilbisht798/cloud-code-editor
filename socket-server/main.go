package main

import (
	//"flag"
	"log"
	"net/http"
)

//var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	log.Print("Server Starting")
	//flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", echo)
	//log.Fatal(http.ListenAndServe(*addr, nil))
	log.Fatal(http.ListenAndServe(":5000", nil))
}
