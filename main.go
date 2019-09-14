package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/xubiosueldos/framework/configuracion"
)

func main() {
	configuracion := configuracion.GetInstance()
	router := newRouter()

	server := http.ListenAndServe(":"+configuracion.Puertomicroservicioliquidacion, router)
	fmt.Println("Microservicio de Liquidación escuchando en el puerto: " + configuracion.Puertomicroservicioliquidacion)

	log.Fatal(server)

}
