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

type strIdsLiquidaciones struct {
	Idsliquidaciones []int `json:"idsliquidaciones"`
}

type strCuentaImporte struct {
	Cuentaid      int     `json:"cuentaid"`
	Importecuenta float32 `json:"importecuenta"`
}

type strLiquidacionContabilizar struct {
	Username        string             `json:"username"`
	Tenant          string             `json:"tenant"`
	Token           string             `json:"token"`
	Liquidacionid   int                `json:"liquidacionid"`
	Descripcion     string             `json:"descripcion"`
	Cuentasimportes []strCuentaImporte `json:"cuentasimportes"`
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

		versionMicroservicio := obtenerVersionLiquidacion()
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

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

			versionMicroservicio := obtenerVersionLiquidacion()
			tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

			db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

			//defer db.Close()
			defer apiclientconexionbd.CerrarDB(db)

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
	}

}

func LiquidacionRemove(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		//Para obtener los parametros por la url
		params := mux.Vars(r)
		liquidacion_id := params["id"]

		versionMicroservicio := obtenerVersionLiquidacion()
		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)

		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		//defer db.Close()
		defer apiclientconexionbd.CerrarDB(db)

		//--Borrado Fisico
		if err := db.Unscoped().Where("id = ?", liquidacion_id).Delete(structLiquidacion.Liquidacion{}).Error; err != nil {

			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		framework.RespondJSON(w, http.StatusOK, framework.Liquidacion+liquidacion_id+framework.MicroservicioEliminado)
	}

}

func LiquidacionContabilizar(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		decoder := json.NewDecoder(r.Body)

		var stridsliquidaciones strIdsLiquidaciones

		if err := decoder.Decode(&stridsliquidaciones); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		versionMicroservicio := obtenerVersionLiquidacion()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := apiclientconexionbd.ObtenerDB(tenant, nombreMicroservicio, versionMicroservicio, AutomigrateTablasPrivadas)

		for i := 0; i < len(stridsliquidaciones.Idsliquidaciones); i++ {

			var liquidacion structLiquidacion.Liquidacion
			liquidacion_id := stridsliquidaciones.Idsliquidaciones[i]

			if err := db.Set("gorm:auto_preload", true).First(&liquidacion, "id = ?", liquidacion_id).Error; gorm.IsRecordNotFoundError(err) {
				framework.RespondError(w, http.StatusNotFound, err.Error())
				return
			}
			if !liquidacion.Estacontabilizada {
				generarAsientoManualDesdeMonolitico(w, r, liquidacion, tokenAutenticacion)
			}

		}
	}

}

func generarAsientoManualDesdeMonolitico(w http.ResponseWriter, r *http.Request, liquidacion structLiquidacion.Liquidacion, tokenAutenticacion *publico.Security) {

	var mapCuentasImportes = make(map[int]float32)

	var strLiquidacionContabilizar strLiquidacionContabilizar
	token := *tokenAutenticacion

	strLiquidacionContabilizar.Tenant = token.Tenant
	strLiquidacionContabilizar.Token = token.Token
	strLiquidacionContabilizar.Username = token.Username
	strLiquidacionContabilizar.Liquidacionid = liquidacion.ID
	strLiquidacionContabilizar.Descripcion = "Asiento Generado para el Legajo: " + liquidacion.Legajo.Legajo + " en el mes: " + liquidacion.Fechaperiodoliquidacion.Month().String() + " del año: " + strconv.Itoa(liquidacion.Fechaperiodoliquidacion.Year())

	agruparLasCuentasDeLasGrillasYSusImportes(liquidacion, mapCuentasImportes)

	strLiquidacionContabilizar.Cuentasimportes = obtenerCuentasImportesLiquidacion(mapCuentasImportes)

	pagesJson, err := json.Marshal(strLiquidacionContabilizar)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url := configuracion.GetUrlMonolitico() + "ContabilizarLiquidacionServlet"

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
