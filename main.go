package main

import (
	"log"
	"net/http"
)

func main() {

	router := newRouter()

	server := http.ListenAndServe(":8085", router)

	log.Fatal(server)

}
