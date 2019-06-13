package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/xubiosueldos/framework/configuracion"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework"
	"github.com/xubiosueldos/liquidacion/structLiquidacion"
)

var nombreMicroservicio string = "liquidacion"

// Sirve para controlar si el server esta OK
func Healthy(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Healthy"))
}

func LiquidacionList(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		versionMicroservicio := obtenerVersionLiquidacion()
		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		defer db.Close()

		var liquidaciones []structLiquidacion.Liquidacion

		db.Find(&liquidaciones)

		framework.RespondJSON(w, http.StatusOK, liquidaciones)
	}

}

func LiquidacionShow(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)
		liquidacion_id := params["id"]

		var liquidacion structLiquidacion.Liquidacion

		versionMicroservicio := obtenerVersionLiquidacion()
		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		defer db.Close()

		//gorm:auto_preload se usa para que complete todos los struct con su informacion
		if err := db.Set("gorm:auto_preload", true).First(&liquidacion, "id = ?", liquidacion_id).Error; gorm.IsRecordNotFoundError(err) {
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, liquidacion)
	}

}

func LiquidacionAdd(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		decoder := json.NewDecoder(r.Body)

		var liquidacion_data structLiquidacion.Liquidacion

		if err := decoder.Decode(&liquidacion_data); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		versionMicroservicio := obtenerVersionLiquidacion()
		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		defer db.Close()

		if err := db.Create(&liquidacion_data).Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusCreated, liquidacion_data)
	}
}

func LiquidacionUpdate(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		params := mux.Vars(r)
		//se convirtió el string en uint para poder comparar
		param_liquidacionid, _ := strconv.ParseInt(params["id"], 10, 64)
		p_liquidacionid := int(param_liquidacionid)

		if p_liquidacionid == 0 {
			framework.RespondError(w, http.StatusNotFound, framework.IdParametroVacio)
			return
		}

		decoder := json.NewDecoder(r.Body)

		var liquidacion_data structLiquidacion.Liquidacion

		if err := decoder.Decode(&liquidacion_data); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		liquidacionid := liquidacion_data.ID

		if p_liquidacionid == liquidacionid || liquidacionid == 0 {

			liquidacion_data.ID = p_liquidacionid

			versionMicroservicio := obtenerVersionLiquidacion()
			db := apiclientconexionbd.ObtenerDB(tokenAutenticacion, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

			defer db.Close()

			if err := db.Save(&liquidacion_data).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			framework.RespondJSON(w, http.StatusOK, liquidacion_data)

		} else {
			framework.RespondError(w, http.StatusNotFound, framework.IdParametroDistintoStruct)
			return
		}
	}

}

func LiquidacionRemove(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		//Para obtener los parametros por la url
		params := mux.Vars(r)
		liquidacion_id := params["id"]

		versionMicroservicio := obtenerVersionLiquidacion()
		db := apiclientconexionbd.ObtenerDB(tokenAutenticacion, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		defer db.Close()

		//--Borrado Fisico
		if err := db.Unscoped().Where("id = ?", liquidacion_id).Delete(structLiquidacion.Liquidacion{}).Error; err != nil {

			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, framework.Liquidacion+liquidacion_id+framework.MicroservicioEliminado)
	}

}

func AutomigrateTablasPrivadas(db *gorm.DB) {

	//para actualizar tablas...agrega columnas e indices, pero no elimina
	db.AutoMigrate(&structLiquidacion.Liquidacion{})

}

func obtenerVersionLiquidacion() int {
	configuracion := configuracion.GetInstance()

	return configuracion.Versionliquidacion
}
