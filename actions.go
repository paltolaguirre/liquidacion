package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/xubiosueldos/framework/configuracion"

	"sueldos-liquidacion/structLiquidacion"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/autenticacion/publico"
	"github.com/xubiosueldos/conexionBD/apiclientconexionbd"
	"github.com/xubiosueldos/framework"
)

type strhelper struct {
	//	gorm.Model
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
	//	Activo      int    `json:"activo"`
}

type strResponse struct {
	//	gorm.Model
	Exists string `json:"exists"`
}

type strHlprServlet struct {
	//	gorm.Model
	Username       string `json:"username"`
	Tenant         string `json:"tenant"`
	Token          string `json:"token"`
	Options        string `json:"options"`
	CuentaContable int    `json:"cuentacontable"`
}

type requestMono struct {
	Value interface{}
	Error error
}

type strIdsLiquidacionesAContabilizar struct {
	Idsliquidacionesacontabilizar []int  `json:"idsliquidacionesacontabilizar"`
	Descripcion                   string `json:"descripcion"`
}

type strTransaccionesIdsAsientosContablesManuales struct {
	Transaccionesidsasientoscontablesmanuales []int `json:"transaccionesidsasientoscontablesmanuales"`
}

type strCuentaImporte struct {
	Cuentaid      int     `json:"cuentaid"`
	Importecuenta float32 `json:"importecuenta"`
}

type strLiquidacionContabilizarDescontabilizar struct {
	Username                string             `json:"username"`
	Tenant                  string             `json:"tenant"`
	Token                   string             `json:"token"`
	Descripcion             string             `json:"descripcion"`
	Cuentasimportes         []strCuentaImporte `json:"cuentasimportes"`
	Asientocontablemanualid int                `json:"asientocontablemanualid"`
}

type respJson struct {
	Codigo    int    `json:"codigo"`
	Respuesta string `json:"respuesta"`
}

type IdsAEliminar struct {
	Ids []int `json:"ids"`
}

type StrDatosAsientoContableManual struct {
	Asientocontablemanualid     int    `json:"asientocontablemanualid"`
	Asientocontablemanualnombre string `json:"asientocontablemanualnombre"`
}

type strCheckLiquidacionesNoContabilizadas struct {
	Cantidadliquidacionesnocontabilizadas int `json:"cantidadliquidacionesnocontabilizadas"`
}

var nombreMicroservicio string = "liquidacion"

// Sirve para controlar si el server esta OK
func Healthy(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Healthy"))
}

func (s *requestMono) requestMonolitico(options string, w http.ResponseWriter, r *http.Request, liquidacion_data structLiquidacion.Liquidacion, tokenAutenticacion *publico.Security, codigo string) *requestMono {

	var strHlprSrv strHlprServlet
	token := *tokenAutenticacion

	strHlprSrv.Options = options
	strHlprSrv.Tenant = token.Tenant
	strHlprSrv.Token = token.Token
	strHlprSrv.Username = token.Username
	strHlprSrv.CuentaContable = *liquidacion_data.Cuentabanco
	pagesJson, err := json.Marshal(strHlprSrv)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url := configuracion.GetUrlMonolitico() + codigo + "GoServlet"

	fmt.Println("URL:>", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pagesJson))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	str := string(body)
	fmt.Println("BYTES RECIBIDOS :", len(str))

	if str == "0" {
		framework.RespondError(w, http.StatusNotFound, "Cuenta Inexistente")
		s.Error = errors.New("Cuenta Inexistente")
	}
	return s
}

func LiquidacionList(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {
		queries := r.URL.Query()

		versionMicroservicio := obtenerVersionLiquidacion()
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		var liquidaciones []structLiquidacion.Liquidacion

		if queries["fechadesde"] == nil && queries["fechahasta"] == nil {
			db.Find(&liquidaciones)
		} else {
			var p_fechadesde string = r.URL.Query()["fechadesde"][0]
			var p_fechahasta string = r.URL.Query()["fechahasta"][0]
			db.Where("fechaperiodoliquidacion BETWEEN ? AND ?", p_fechadesde, p_fechahasta).Find(&liquidaciones)
		}

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
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

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
		//&liquidacion_data para decirle que es la var que no tiene datos y va a tener que rellenar
		if err := decoder.Decode(&liquidacion_data); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		versionMicroservicio := obtenerVersionLiquidacion()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		var requestMono requestMono

		if err := requestMono.requestMonolitico("CANQUERY", w, r, liquidacion_data, tokenAutenticacion, "cuenta").Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

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

		versionMicroservicio := obtenerVersionLiquidacion()
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)
		defer apiclientconexionbd.CerrarDB(db)

		if !liquidacionContabilizada(p_liquidacionid, db) {
			decoder := json.NewDecoder(r.Body)

			var liquidacion_data structLiquidacion.Liquidacion

			if err := decoder.Decode(&liquidacion_data); err != nil {
				framework.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			defer r.Body.Close()

			liquidacionid := liquidacion_data.ID

			var requestMono requestMono

			if err := requestMono.requestMonolitico("CANQUERY", w, r, liquidacion_data, tokenAutenticacion, "cuenta").Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			if p_liquidacionid == liquidacionid || liquidacionid == 0 {

				liquidacion_data.ID = p_liquidacionid

				//abro una transacción para que si hay un error no persista en la DB
				tx := db.Begin()

				//modifico el legajo de acuerdo a lo enviado en el json
				if err := tx.Save(&liquidacion_data).Error; err != nil {
					tx.Rollback()
					framework.RespondError(w, http.StatusInternalServerError, err.Error())
					return
				}

				//despues de modificar, recorro los descuentos asociados a la liquidacion para ver si alguno fue eliminado logicamente y lo elimino de la BD
				if err := tx.Model(structLiquidacion.Descuento{}).Unscoped().Where("liquidacionid = ? AND deleted_at is not null", liquidacionid).Delete(structLiquidacion.Descuento{}).Error; err != nil {
					tx.Rollback()
					framework.RespondError(w, http.StatusInternalServerError, err.Error())
					return
				}

				//despues de modificar, recorro los importes remunerativos asociados a la liquidacion para ver si fue eliminado logicamente y lo elimino de la BD
				if err := tx.Model(structLiquidacion.Importenoremunerativo{}).Unscoped().Where("liquidacionid = ? AND deleted_at is not null", liquidacionid).Delete(structLiquidacion.Importenoremunerativo{}).Error; err != nil {
					tx.Rollback()
					framework.RespondError(w, http.StatusInternalServerError, err.Error())
					return
				}

				//despues de modificar, recorro los importes no remunerativos asociados a la liquidacion para ver si fue eliminado logicamente y lo elimino de la BD
				if err := tx.Model(structLiquidacion.Importenoremunerativo{}).Unscoped().Where("liquidacionid = ? AND deleted_at is not null", liquidacionid).Delete(structLiquidacion.Importenoremunerativo{}).Error; err != nil {
					tx.Rollback()
					framework.RespondError(w, http.StatusInternalServerError, err.Error())
					return
				}

				//despues de modificar, recorro las retenciones asociadas a la liquidacion para ver si fue eliminado logicamente y lo elimino de la BD
				if err := tx.Model(structLiquidacion.Retencion{}).Unscoped().Where("liquidacionid = ? AND deleted_at is not null", liquidacionid).Delete(structLiquidacion.Retencion{}).Error; err != nil {
					tx.Rollback()
					framework.RespondError(w, http.StatusInternalServerError, err.Error())
					return
				}

				tx.Commit()

				framework.RespondJSON(w, http.StatusOK, liquidacion_data)

			} else {
				framework.RespondError(w, http.StatusNotFound, framework.IdParametroDistintoStruct)
				return
			}
		} else {
			framework.RespondError(w, http.StatusNotFound, framework.Modificarliquidacioncontabilizada)
			return
		}
	}

}

func LiquidacionRemove(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		//Para obtener los parametros por la url
		params := mux.Vars(r)
		param_liquidacionid, _ := strconv.ParseInt(params["id"], 10, 64)
		p_liquidacionid := int(param_liquidacionid)

		versionMicroservicio := obtenerVersionLiquidacion()
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		if !liquidacionContabilizada(p_liquidacionid, db) {

			//--Borrado Fisico
			if err := db.Unscoped().Where("id = ?", p_liquidacionid).Delete(structLiquidacion.Liquidacion{}).Error; err != nil {

				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			framework.RespondJSON(w, http.StatusOK, framework.Liquidacion+strconv.Itoa(p_liquidacionid)+framework.MicroservicioEliminado)
		} else {
			framework.RespondError(w, http.StatusNotFound, framework.Eliminarliquidacioncontabilizada)
			return
		}
	}

}

func LiquidacionContabilizar(w http.ResponseWriter, r *http.Request) {
	var mapCuentasImportes = make(map[int]float32)
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		decoder := json.NewDecoder(r.Body)

		var strIdsLiquidaciones strIdsLiquidacionesAContabilizar

		if err := decoder.Decode(&strIdsLiquidaciones); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		versionMicroservicio := obtenerVersionLiquidacion()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		defer apiclientconexionbd.CerrarDB(db)
		var liquidaciones_ids string
		descripcion_asiento := strIdsLiquidaciones.Descripcion
		if len(strIdsLiquidaciones.Idsliquidacionesacontabilizar) > 0 {
			liquidaciones_ids = "(" + strings.Trim(strings.Replace(fmt.Sprint(strIdsLiquidaciones.Idsliquidacionesacontabilizar), " ", ",", -1), "[]") + ")"
		} else {
			framework.RespondError(w, http.StatusNotFound, framework.Seleccioneunregistro)
		}
		var liquidaciones []structLiquidacion.Liquidacion
		db.Set("gorm:auto_preload", true).Find(&liquidaciones, "id IN "+liquidaciones_ids)

		if len(liquidaciones) > 0 {
			if checkLiquidacionesNoContabilizadas(liquidaciones, liquidaciones_ids, db) {
				for i := 0; i < len(liquidaciones); i++ {
					agruparLasCuentasDeLasGrillasYSusImportes(liquidaciones[i], mapCuentasImportes)
				}
			} else {
				framework.RespondError(w, http.StatusNotFound, framework.Seleccionaronliquidacionescontabilizadas)
				return
			}
		}

		generarAsientoManualDesdeMonolitico(w, r, liquidaciones, mapCuentasImportes, tokenAutenticacion, descripcion_asiento, 0, db)

	}

}

func liquidacionContabilizada(liquidacion_id int, db *gorm.DB) bool {
	var liquidacion structLiquidacion.Liquidacion
	db.Set("gorm:auto_preload", true).First(&liquidacion, "id = ?", liquidacion_id)
	return liquidacion.Estacontabilizada

}

func checkLiquidacionesNoContabilizadas(liquidaciones []structLiquidacion.Liquidacion, liquidaciones_ids string, db *gorm.DB) bool {
	var strCheckLiquidacionesNoContabilizadas strCheckLiquidacionesNoContabilizadas

	db.Raw("SELECT COUNT(ID) AS cantidadliquidacionesnocontabilizadas FROM LIQUIDACION WHERE ID IN " + liquidaciones_ids + " AND ESTACONTABILIZADA = " + strconv.FormatBool(false)).Scan(&strCheckLiquidacionesNoContabilizadas)
	return len(liquidaciones) == strCheckLiquidacionesNoContabilizadas.Cantidadliquidacionesnocontabilizadas
}

func generarAsientoManualDesdeMonolitico(w http.ResponseWriter, r *http.Request, liquidaciones []structLiquidacion.Liquidacion, mapCuentasImportes map[int]float32, tokenAutenticacion *publico.Security, descripcion string, asientomanualtransaccionid int, db *gorm.DB) {

	resp := requestMonoliticoContabilizarDescontabilizarLiquidaciones(r, mapCuentasImportes, tokenAutenticacion, descripcion, asientomanualtransaccionid, db)

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	var datosAsientoContableManual StrDatosAsientoContableManual

	defer resp.Body.Close()

	json.Unmarshal(body, &datosAsientoContableManual)

	if resp.StatusCode == http.StatusOK {
		marcarLiquidacionesComoContabilizadas(liquidaciones, datosAsientoContableManual, db)
		var respuestaJson respJson
		respuestaJson.Codigo = http.StatusOK
		respuestaJson.Respuesta = "Se contabilizaron correctamente " + strconv.Itoa(len(liquidaciones)) + " liquidaciones"
		framework.RespondJSON(w, http.StatusOK, respuestaJson)
	} else {
		str := string(body)
		framework.RespondError(w, http.StatusNotFound, str)
	}

}

func requestMonoliticoContabilizarDescontabilizarLiquidaciones(r *http.Request, mapCuentasImportes map[int]float32, tokenAutenticacion *publico.Security, descripcion string, asientomanualtransaccionid int, db *gorm.DB) *http.Response {

	var strLiquidacionContabilizarDescontabilizar strLiquidacionContabilizarDescontabilizar
	token := *tokenAutenticacion

	strLiquidacionContabilizarDescontabilizar.Tenant = token.Tenant
	strLiquidacionContabilizarDescontabilizar.Token = token.Token
	strLiquidacionContabilizarDescontabilizar.Username = token.Username
	if asientomanualtransaccionid == 0 {
		if descripcion == "" {
			descripcion = framework.Descripcionasientomanualcontableliquidacionescontabilizadas
		}
		strLiquidacionContabilizarDescontabilizar.Descripcion = descripcion
		strLiquidacionContabilizarDescontabilizar.Cuentasimportes = obtenerCuentasImportesLiquidacion(mapCuentasImportes)
	} else {
		strLiquidacionContabilizarDescontabilizar.Asientocontablemanualid = asientomanualtransaccionid
	}
	pagesJson, err := json.Marshal(strLiquidacionContabilizarDescontabilizar)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url := configuracion.GetUrlMonolitico() + "ContabilizarLiquidacionServlet"

	fmt.Println("URL:>", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pagesJson))

	if err != nil {
		fmt.Println("Error: ", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	return resp
}

func marcarLiquidacionesComoContabilizadas(liquidaciones []structLiquidacion.Liquidacion, datosAsientoContableManual StrDatosAsientoContableManual, db *gorm.DB) {
	for i := 0; i < len(liquidaciones); i++ {
		db.Model(&liquidaciones[i]).Update("Estacontabilizada", true)
		db.Model(&liquidaciones[i]).Update("Asientomanualtransaccionid", datosAsientoContableManual.Asientocontablemanualid)
		db.Model(&liquidaciones[i]).Update("Asientomanualnombre", datosAsientoContableManual.Asientocontablemanualnombre)
	}
}

func agruparLasCuentasDeLasGrillasYSusImportes(liquidacion structLiquidacion.Liquidacion, mapCuentasImportes map[int]float32) {

	var cuentaContable *int

	for i := 0; i < len(liquidacion.Descuentos); i++ {
		cuentaContable = liquidacion.Descuentos[i].Concepto.CuentaContable
		importeUnitario := *liquidacion.Descuentos[i].Importeunitario

		importe := mapCuentasImportes[*cuentaContable]
		mapCuentasImportes[*cuentaContable] = importe + importeUnitario

	}

	for j := 0; j < len(liquidacion.Importesnoremunerativos); j++ {
		cuentaContable = liquidacion.Importesnoremunerativos[j].Concepto.CuentaContable
		importeUnitario := *liquidacion.Importesnoremunerativos[j].Importeunitario

		importe := mapCuentasImportes[*cuentaContable]
		mapCuentasImportes[*cuentaContable] = importe + importeUnitario
	}

	for k := 0; k < len(liquidacion.Importesremunerativos); k++ {
		cuentaContable = liquidacion.Importesremunerativos[k].Concepto.CuentaContable
		importeUnitario := *liquidacion.Importesremunerativos[k].Importeunitario

		importe := mapCuentasImportes[*cuentaContable]
		mapCuentasImportes[*cuentaContable] = importe + importeUnitario
	}

	for m := 0; m < len(liquidacion.Retenciones); m++ {
		cuentaContable = liquidacion.Retenciones[m].Concepto.CuentaContable
		importeUnitario := *liquidacion.Retenciones[m].Importeunitario

		importe := mapCuentasImportes[*cuentaContable]
		mapCuentasImportes[*cuentaContable] = importe + importeUnitario
	}

}

func obtenerCuentasImportesLiquidacion(mapCuentasImportes map[int]float32) []strCuentaImporte {
	var arrayStrCuentaImporte []strCuentaImporte

	for cuenta, importe := range mapCuentasImportes {
		var strcuentaimporte strCuentaImporte
		strcuentaimporte.Cuentaid = cuenta
		strcuentaimporte.Importecuenta = importe
		arrayStrCuentaImporte = append(arrayStrCuentaImporte, strcuentaimporte)
	}

	return arrayStrCuentaImporte
}

/*func LiquidacionDesContabilizar(w http.ResponseWriter, r *http.Request) {
	var respuestaDescontabilizar = make(map[int]respJson)
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		decoder := json.NewDecoder(r.Body)

		var strTransaccionesIdsAsientosContablesManuales strTransaccionesIdsAsientosContablesManuales

		if err := decoder.Decode(&strTransaccionesIdsAsientosContablesManuales); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		versionMicroservicio := obtenerVersionLiquidacion()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		defer apiclientconexionbd.CerrarDB(db)

		for i := 0; i < len(strTransaccionesIdsAsientosContablesManuales.Transaccionesidsasientoscontablesmanuales); i++ {
			asientomanualtransaccionid := strTransaccionesIdsAsientosContablesManuales.Transaccionesidsasientoscontablesmanuales[i]
			if checkAsientoManualTransaccionID(w, asientomanualtransaccionid, respuestaDescontabilizar, db) {
				descontabilizarLiquidaciones(w, r, asientomanualtransaccionid, tokenAutenticacion, respuestaDescontabilizar, db)
			} else {
				var respuestaJson respJson
				respuestaJson.Codigo = http.StatusOK
				respuestaJson.Respuesta = "No se encuentra contabilizando ninguna liquidación"
				respuestaDescontabilizar[asientomanualtransaccionid] = respuestaJson
			}
		}
	}

	framework(w, http.StatusOK, respuestaDescontabilizar)
}
func checkAsientoManualTransaccionID(w http.ResponseWriter, asientomanualtransaccionid int, respuestaDescontabilizar map[int]respJson, db *gorm.DB) bool {

	liquidaciones := buscarLiquidacionesAsientoManualTransaccion(asientomanualtransaccionid, respuestaDescontabilizar, db)
	return len(liquidaciones) > 0
}

func buscarLiquidacionesAsientoManualTransaccion(asientomanualtransaccionid int, w http.ResponseWriter, db *gorm.DB) []structLiquidacion.Liquidacion {
	var liquidaciones []structLiquidacion.Liquidacion

	if err := db.Find(&liquidaciones, "asientomanualtransaccionid = ?", asientomanualtransaccionid).Error; gorm.IsRecordNotFoundError(err) {
		framework.RespondError(w, http.StatusNotFound, err.Error())

	}

	return &liquidaciones
}

func descontabilizarLiquidaciones(w http.ResponseWriter, r *http.Request, asientomanualtransaccionid int, tokenAutenticacion *publico.Security, respuestaDescontabilizar map[int]string, db *gorm.DB) {

	resp := requestMonoliticoContabilizarDescontabilizarLiquidaciones(r, nil, tokenAutenticacion, "", asientomanualtransaccionid, db)
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		blanquearAsientoManualTransaccionYNombreEnLiquidaciones(w, asientomanualtransaccionid, db)
		var respuestaJson respJson
		respuestaJson.Codigo = http.StatusOK
		respuestaJson.Respuesta = ""
		framework.RespondJSON(w, http.StatusOK, respuestaJson)
	} else {
		str := string(body)
		framework.RespondError(w, http.StatusNotFound, str)
	}
}

func blanquearAsientoManualTransaccionYNombreEnLiquidaciones(w http.ResponseWriter, asientocontablemanualid int, db *gorm.DB) {

	db.Model(&liquidaciones).Where("asientomanualtransaccionid = " + strconv.Itoa(asientocontablemanualid)).Updates(structLiquidacion.Liquidacion{Asientomanualtransaccionid: 0, Asientomanualnombre: ""})
}
*/
func LiquidacionesRemoveMasivo(w http.ResponseWriter, r *http.Request) {
	var resultadoDeEliminacion = make(map[int]string)
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		var idsEliminar IdsAEliminar
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&idsEliminar); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		versionMicroservicio := obtenerVersionLiquidacion()
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		defer apiclientconexionbd.CerrarDB(db)

		if len(idsEliminar.Ids) > 0 {
			for i := 0; i < len(idsEliminar.Ids); i++ {
				liquidacion_id := idsEliminar.Ids[i]
				if err := db.Unscoped().Where("id = ?", liquidacion_id).Delete(structLiquidacion.Liquidacion{}).Error; err != nil {
					//framework.RespondError(w, http.StatusInternalServerError, err.Error())
					resultadoDeEliminacion[liquidacion_id] = string(err.Error())

				} else {
					resultadoDeEliminacion[liquidacion_id] = "Fue eliminado con exito"
				}
			}
		} else {
			framework.RespondError(w, http.StatusInternalServerError, framework.Seleccioneunregistro)
		}

		framework.RespondJSON(w, http.StatusOK, resultadoDeEliminacion)
	}

}

func AutomigrateTablasPrivadas(db *gorm.DB) {

	//para actualizar tablas...agrega columnas e indices, pero no elimina
	db.AutoMigrate(&structLiquidacion.Descuento{}, &structLiquidacion.Importenoremunerativo{}, &structLiquidacion.Importeremunerativo{}, &structLiquidacion.Retencion{}, &structLiquidacion.Liquidacion{})

	db.Model(&structLiquidacion.Descuento{}).AddForeignKey("liquidacionid", "liquidacion(id)", "CASCADE", "CASCADE")
	db.Model(&structLiquidacion.Importenoremunerativo{}).AddForeignKey("liquidacionid", "liquidacion(id)", "CASCADE", "CASCADE")
	db.Model(&structLiquidacion.Importeremunerativo{}).AddForeignKey("liquidacionid", "liquidacion(id)", "CASCADE", "CASCADE")
	db.Model(&structLiquidacion.Retencion{}).AddForeignKey("liquidacionid", "liquidacion(id)", "CASCADE", "CASCADE")

}

func obtenerVersionLiquidacion() int {
	configuracion := configuracion.GetInstance()

	return configuracion.Versionliquidacion
}
