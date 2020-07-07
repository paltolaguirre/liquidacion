package main

import (
	"fmt"
	"github.com/xubiosueldos/framework/configuracion"
	"log"
	"net/http"
)

func main() {
	configuracion := configuracion.GetInstance()
	router := newRouter()

	fmt.Println("Microservicio de Liquidación escuchando en el puerto: " + configuracion.Puertomicroservicioliquidacion)
	server := http.ListenAndServe(":"+configuracion.Puertomicroservicioliquidacion, router)

	log.Fatal(server)

}
