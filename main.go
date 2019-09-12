package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/conexionBD/Autenticacion/structAutenticacion"
	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework/configuracion"
)

func main() {
	configuracion := configuracion.GetInstance()
	var tokenAutenticacion structAutenticacion.Security
	tokenAutenticacion.Tenant = "public"

	tenant := apiclientautenticacion.ObtenerTenant(&tokenAutenticacion)
	apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, obtenerVersionLiquidacion())

	router := newRouter()

	server := http.ListenAndServe(":"+configuracion.Puertomicroservicioliquidacion, router)
	fmt.Println("Microservicio de Liquidación escuchando en el puerto: " + configuracion.Puertomicroservicioliquidacion)

	log.Fatal(server)

}
