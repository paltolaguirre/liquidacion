package main

import (
	"fmt"
	"log"
	"net/http"
	"sueldos-liquidacion/structLiquidacion"

	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework/configuracion"
)

func main() {
	configuracion := configuracion.GetInstance()
	var tokenAutenticacion publico.Security
	tokenAutenticacion.Tenant = "public"

	tenant := apiclientautenticacion.ObtenerTenant(&tokenAutenticacion)
	apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, obtenerVersionLiquidacion(), AutomigrateTablasPublicas)

	router := newRouter()

	server := http.ListenAndServe(":"+configuracion.Puertomicroservicioliquidacion, router)
	fmt.Println("Microservicio de Liquidación escuchando en el puerto: " + configuracion.Puertomicroservicioliquidacion)

	log.Fatal(server)

}

func AutomigrateTablasPublicas(db *gorm.DB) {

	//para actualizar tablas...agrega columnas e indices, pero no elimina
	db.AutoMigrate(&structLiquidacion.Liquidacioncondicionpago{}, &structLiquidacion.Liquidaciontipo{})

}
