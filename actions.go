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

	"github.com/xubiosueldos/concepto/structConcepto"
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
	Username                   string             `json:"username"`
	Tenant                     string             `json:"tenant"`
	Token                      string             `json:"token"`
	Descripcion                string             `json:"descripcion"`
	Cuentasimportes            []strCuentaImporte `json:"cuentasimportes"`
	Asientomanualtransaccionid int                `json:"asientomanualtransaccionid"`
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
			db.Set("gorm:auto_preload", true).Find(&liquidaciones)
		} else {
			var p_fechadesde string = r.URL.Query()["fechadesde"][0]
			var p_fechahasta string = r.URL.Query()["fechahasta"][0]
			db.Set("gorm:auto_preload", true).Where("fechaperiodoliquidacion BETWEEN ? AND ?", p_fechadesde, p_fechahasta).Find(&liquidaciones)
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
	fmt.Println("La URL accedida: " + r.URL.String())
	var mapCuentasImportes = make(map[int]float32)
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {
		fmt.Println("Token valido")
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
					agruparLasCuentasDeLasGrillasYSusImportes(liquidaciones[i], mapCuentasImportes, r)
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
		strLiquidacionContabilizarDescontabilizar.Asientomanualtransaccionid = asientomanualtransaccionid
	}
	pagesJson, err := json.Marshal(strLiquidacionContabilizarDescontabilizar)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url := configuracion.GetUrlMonolitico() + "ContabilizarLiquidacionServlet"

	fmt.Println("Se hace un request al monolitico con la siguiente URL:>", url)
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

func agruparLasCuentasDeLasGrillasYSusImportes(liquidacion structLiquidacion.Liquidacion, mapCuentasImportes map[int]float32, r *http.Request) {
	fmt.Println("Se agrupan las cuentas para la Liquidacion: " + strconv.Itoa(liquidacion.ID))
	var cuentaContable *int

	for i := 0; i < len(liquidacion.Descuentos); i++ {

		descuento := liquidacion.Descuentos[i]
		concepto := obtenerConcepto(*descuento.Conceptoid, r)
		cuentaContable = concepto.CuentaContable
		importeUnitario := *descuento.Importeunitario

		importe := mapCuentasImportes[*cuentaContable]
		mapCuentasImportes[*cuentaContable] = importe + importeUnitario

	}

	for j := 0; j < len(liquidacion.Importesnoremunerativos); j++ {

		importenoremunerativo := liquidacion.Importesnoremunerativos[j]
		concepto := obtenerConcepto(*importenoremunerativo.Conceptoid, r)
		cuentaContable = concepto.CuentaContable
		importeUnitario := *importenoremunerativo.Importeunitario

		importe := mapCuentasImportes[*cuentaContable]
		mapCuentasImportes[*cuentaContable] = importe + importeUnitario
	}

	for k := 0; k < len(liquidacion.Importesremunerativos); k++ {

		importeremunerativo := liquidacion.Importesremunerativos[k]
		concepto := obtenerConcepto(*importeremunerativo.Conceptoid, r)
		cuentaContable = concepto.CuentaContable
		importeUnitario := *importeremunerativo.Importeunitario

		importe := mapCuentasImportes[*cuentaContable]
		mapCuentasImportes[*cuentaContable] = importe + importeUnitario
	}

	for m := 0; m < len(liquidacion.Retenciones); m++ {
		retencion := liquidacion.Retenciones[m]
		concepto := obtenerConcepto(*retencion.Conceptoid, r)
		cuentaContable = concepto.CuentaContable
		importeUnitario := *retencion.Importeunitario

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

func LiquidacionDesContabilizar(w http.ResponseWriter, r *http.Request) {
	var respuestaDescontabilizar = make(map[int]respJson)
	var cantidadLiquidaciones int
	var cantidadAsientoManualTransaccionID int
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
			existeAsientoManualID, liquidaciones := checkAsientoManualTransaccionID(w, asientomanualtransaccionid, respuestaDescontabilizar, db)
			if existeAsientoManualID {
				cantidadLiquidacionesDescontabilizadas, cantidadTransaccionIDDescontabilizada := descontabilizarLiquidaciones(w, r, liquidaciones, asientomanualtransaccionid, tokenAutenticacion, respuestaDescontabilizar, db)
				cantidadLiquidaciones = cantidadLiquidaciones + cantidadLiquidacionesDescontabilizadas
				cantidadAsientoManualTransaccionID = cantidadAsientoManualTransaccionID + cantidadTransaccionIDDescontabilizada
			} else {
				var respuestaJson respJson
				respuestaJson.Codigo = http.StatusConflict
				respuestaJson.Respuesta = "La TransaccionID " + strconv.Itoa(asientomanualtransaccionid) + " no se encuentra contabilizando ninguna liquidación"
				respuestaDescontabilizar[asientomanualtransaccionid] = respuestaJson
			}
		}
		if cantidadLiquidaciones > 0 {
			var respuestaJson respJson
			respuestaJson.Codigo = http.StatusOK
			respuestaJson.Respuesta = "Se descontabilizaron correctamente " + strconv.Itoa(cantidadLiquidaciones) + " Liquidaciones, correspondientes a " + strconv.Itoa(cantidadAsientoManualTransaccionID) + " Asientos Manuales"
			respuestaDescontabilizar[-1] = respuestaJson
		}

		framework.RespondJSON(w, http.StatusOK, respuestaDescontabilizar)
	} else {
		framework.RespondError(w, http.StatusConflict, "El token utilizado es invalido")
	}

}
func checkAsientoManualTransaccionID(w http.ResponseWriter, asientomanualtransaccionid int, respuestaDescontabilizar map[int]respJson, db *gorm.DB) (bool, []structLiquidacion.Liquidacion) {

	liquidaciones := buscarLiquidacionesAsientoManualTransaccion(asientomanualtransaccionid, respuestaDescontabilizar, w, db)
	return len(liquidaciones) > 0, liquidaciones
}

func buscarLiquidacionesAsientoManualTransaccion(asientomanualtransaccionid int, respuestaDescontabilizar map[int]respJson, w http.ResponseWriter, db *gorm.DB) []structLiquidacion.Liquidacion {
	var liquidaciones []structLiquidacion.Liquidacion

	if err := db.Find(&liquidaciones, "asientomanualtransaccionid = ?", asientomanualtransaccionid).Error; gorm.IsRecordNotFoundError(err) {
		framework.RespondError(w, http.StatusNotFound, err.Error())

	}

	return liquidaciones
}

func descontabilizarLiquidaciones(w http.ResponseWriter, r *http.Request, liquidaciones []structLiquidacion.Liquidacion, asientomanualtransaccionid int, tokenAutenticacion *publico.Security, respuestaDescontabilizar map[int]respJson, db *gorm.DB) (int, int) {

	resp := requestMonoliticoContabilizarDescontabilizarLiquidaciones(r, nil, tokenAutenticacion, "", asientomanualtransaccionid, db)
	body, err := ioutil.ReadAll(resp.Body)
	var cantLiquidaciones int
	var cantAsientoManualTransaccionID int
	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer resp.Body.Close()

	var respuestaJson respJson
	if resp.StatusCode == http.StatusOK {
		blanquearAsientoManualTransaccionYNombreEnLiquidaciones(w, liquidaciones, asientomanualtransaccionid, db)
		cantLiquidaciones = cantLiquidaciones + len(liquidaciones)
		cantAsientoManualTransaccionID = cantAsientoManualTransaccionID + 1
	} else {
		str := string(body)
		respuestaJson.Codigo = http.StatusNotFound
		respuestaJson.Respuesta = str
		respuestaDescontabilizar[asientomanualtransaccionid] = respuestaJson

	}

	return cantLiquidaciones, cantAsientoManualTransaccionID

}

func blanquearAsientoManualTransaccionYNombreEnLiquidaciones(w http.ResponseWriter, liquidaciones []structLiquidacion.Liquidacion, asientocontablemanualid int, db *gorm.DB) {

	//Utilice la forma "manual" para updetear porque la otra no me funcionaba! (ver!)
	db.Raw("UPDATE LIQUIDACION SET Asientomanualtransaccionid = 0, Asientomanualnombre = '', Estacontabilizada = false WHERE Asientomanualtransaccionid = " + strconv.Itoa(asientocontablemanualid)).Scan(&liquidaciones)

	//var liquidacion structLiquidacion.Liquidacion
	//db.Model(&liquidacion).Where("Asientomanualtransaccionid = ?", asientocontablemanualid).UpdateColumns(structLiquidacion.Liquidacion{Asientomanualtransaccionid: 0, Asientomanualnombre: "", Estacontabilizada: false})
	//db.Model(&liquidaciones).Updates(structLiquidacion.Liquidacion{Asientomanualtransaccionid: 0, Asientomanualnombre: "", })
}

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

func obtenerConcepto(conceptoid int, r *http.Request) *structConcepto.Concepto {

	var concepto structConcepto.Concepto

	config := configuracion.GetInstance()

	url := configuracion.GetUrlMicroservicio(config.Puertomicroservicioconcepto) + "concepto/conceptos/" + strconv.Itoa(conceptoid)

	//url := "http://localhost:8084/conceptos/" + strconv.Itoa(conceptoid)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	header := r.Header.Get("Authorization")

	req.Header.Add("Authorization", header)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("URL:", url)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	str := string(body)

	json.Unmarshal([]byte(str), &concepto)

	return &concepto

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
