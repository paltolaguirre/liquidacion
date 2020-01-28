package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/xubiosueldos/conexionBD"

	"github.com/xubiosueldos/conexionBD/Concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"

	"git-codecommit.us-east-1.amazonaws.com/v1/repos/sueldos-liquidacion/calculosAutomaticos"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xubiosueldos/autenticacion/apiclientautenticacion"
	"github.com/xubiosueldos/conexionBD/Autenticacion/structAutenticacion"
	"github.com/xubiosueldos/framework"
	"github.com/xubiosueldos/monoliticComunication"
)

type strIdsLiquidacionesAContabilizar struct {
	Idsliquidacionesacontabilizar []int  `json:"idsliquidacionesacontabilizar"`
	Descripcion                   string `json:"descripcion"`
	Fechaasiento                  string `json:"fechaasiento"`
}

type strTransaccionesIdsAsientosContablesManuales struct {
	Transaccionesidsasientoscontablesmanuales []int `json:"transaccionesidsasientoscontablesmanuales"`
}

type strCuentaImporte struct {
	Cuentaid      int     `json:"cuentaid"`
	Importecuenta float64 `json:"importecuenta"`
}

type strCuentaImporteTipoGrilla struct {
	Cuentaid      int     `json:"cuentaid"`
	Importecuenta float64 `json:"importecuenta"`
	Tipogrilla    int     `json:"tipogrilla"`
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

type strCheckLiquidacionesNoContabilizadas struct {
	Cantidadliquidacionesnocontabilizadas int `json:"cantidadliquidacionesnocontabilizadas"`
}

type DuplicarLiquidaciones struct {
	Liquidaciondefaultvalues structLiquidacion.Liquidacion `json:"liquidaciondefaultvalues"`
	Idstoreplicate           []int                         `json:"idstoreplicate"`
}

type ResultProcesamientoMasivo struct {
	Processid string                `json:"processid"`
	Result    []ProcesamientoStatus `json:"result"`
}

type ProcesamientoStatus struct {
	Id      int    `json:"id"`
	Tipo    string `json:"tipo"`
	Codigo  int    `json:"codigo"`
	Mensaje string `json:"mensaje"`
}

type StrDatosAsientoContableManual struct {
	Asientocontablemanualid     int    `json:"asientocontablemanualid"`
	Asientocontablemanualnombre string `json:"asientocontablemanualnombre"`
	Statuscode                  int    `json:"statuscode"`
}

type StrDatosAsientoContableManualBlanquear struct {
	Asientocontablemanualid int    `json:"asientocontablemanualid"`
	Tokensecurityencode     string `json:"tokensecurityencode"`
}

type StrCalculoAutomaticoConceptoId struct {
	Conceptoid      *int     `json:"conceptoid"`
	Importeunitario *float64 `json:"importeunitario" `
}

var nombreMicroservicio string = "liquidacion"

// Sirve para controlar si el server esta OK
func Healthy(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Healthy"))
}

func LiquidacionList(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {
		queries := r.URL.Query()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

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

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		//gorm:auto_preload se usa para que complete todos los struct con su informacion
		if err := db.Set("gorm:auto_preload", true).First(&liquidacion, "id = ?", liquidacion_id).Error; gorm.IsRecordNotFoundError(err) {
			framework.RespondError(w, http.StatusNotFound, err.Error())
			return
		}

		bancoID := liquidacion.Cuentabancoid
		if bancoID != nil {
			cuentaBanco := monoliticComunication.Obtenerbanco(w, r, tokenAutenticacion, strconv.Itoa(*bancoID))
			liquidacion.Cuentabanco = cuentaBanco
		}

		bancoaportejubilatorioID := liquidacion.Bancoaportejubilatorioid
		if bancoaportejubilatorioID != nil {
			bancoAporteJubilatorio := monoliticComunication.Obtenerbanco(w, r, tokenAutenticacion, strconv.Itoa(*bancoaportejubilatorioID))
			liquidacion.Bancoaportejubilatorio = bancoAporteJubilatorio
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

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		if err := monoliticComunication.Checkexistebanco(w, r, tokenAutenticacion, strconv.Itoa(*liquidacion_data.Cuentabancoid)).Error; err != nil {
			framework.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err := monoliticComunication.Checkexistebanco(w, r, tokenAutenticacion, strconv.Itoa(*liquidacion_data.Bancoaportejubilatorioid)).Error; err != nil {
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

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		if !liquidacionContabilizada(p_liquidacionid, db) {
			decoder := json.NewDecoder(r.Body)

			var liquidacion_data structLiquidacion.Liquidacion

			if err := decoder.Decode(&liquidacion_data); err != nil {
				framework.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			defer r.Body.Close()

			liquidacionid := liquidacion_data.ID

			if err := monoliticComunication.Checkexistebanco(w, r, tokenAutenticacion, strconv.Itoa(*liquidacion_data.Cuentabancoid)).Error; err != nil {
				framework.RespondError(w, http.StatusInternalServerError, err.Error())
				return
			}

			if err := monoliticComunication.Checkexistebanco(w, r, tokenAutenticacion, strconv.Itoa(*liquidacion_data.Bancoaportejubilatorioid)).Error; err != nil {
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

				if err := tx.Model(structLiquidacion.Liquidacionitem{}).Unscoped().Where("liquidacionid = ? AND deleted_at is not null", liquidacionid).Delete(structLiquidacion.Liquidacionitem{}).Error; err != nil {
					tx.Rollback()
					framework.RespondError(w, http.StatusInternalServerError, err.Error())
					return
				}
				//despues de modificar, recorro los descuentos asociados a la liquidacion para ver si alguno fue eliminado logicamente y lo elimino de la BD
				/*	if err := tx.Model(structLiquidacion.Descuento{}).Unscoped().Where("liquidacionid = ? AND deleted_at is not null", liquidacionid).Delete(structLiquidacion.Descuento{}).Error; err != nil {
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
					}*/

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

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

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
	var mapCuentasImportes = make(map[int]float64)
	var strCuentaImporteTipoGrillas []strCuentaImporteTipoGrilla
	var strCuentasImportes []strCuentaImporte
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {
		decoder := json.NewDecoder(r.Body)

		var strIdsLiquidaciones strIdsLiquidacionesAContabilizar

		if err := decoder.Decode(&strIdsLiquidaciones); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)
		var liquidaciones_ids string
		descripcion_asiento := strIdsLiquidaciones.Descripcion
		fecha_asiento := strIdsLiquidaciones.Fechaasiento
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
					obtenerCuentasImportesYTipoDeGrillas(liquidaciones[i], &strCuentaImporteTipoGrillas, r)
				}

				obtenerCuentasImportes(strCuentaImporteTipoGrillas, &strCuentasImportes)

				agruparCuentas(strCuentasImportes, mapCuentasImportes)

			} else {
				framework.RespondError(w, http.StatusNotFound, framework.Seleccionaronliquidacionescontabilizadas)
				return
			}
		}

		generarAsientoManualDesdeMonolitico(w, r, liquidaciones, mapCuentasImportes, tokenAutenticacion, descripcion_asiento, fecha_asiento, 0, db)

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

func generarAsientoManualDesdeMonolitico(w http.ResponseWriter, r *http.Request, liquidaciones []structLiquidacion.Liquidacion, mapCuentasImportes map[int]float64, tokenAutenticacion *structAutenticacion.Security, descripcion string, fechaasiento string, asientomanualtransaccionid int, db *gorm.DB) {
	var cuentasImportes []monoliticComunication.StrCuentaImporte
	cuentasImportes = obtenerCuentasImportesLiquidacion(mapCuentasImportes)
	datosAsientoContableManual := monoliticComunication.Generarasientomanual(w, r, cuentasImportes, tokenAutenticacion, descripcion, fechaasiento)

	if err := monoliticComunication.Checkgeneroasientomanual(datosAsientoContableManual).Error; err != nil {
		framework.RespondError(w, http.StatusNotFound, err.Error())
	} else {
		marcarLiquidacionesComoContabilizadas(liquidaciones, datosAsientoContableManual, db)
		var respuestaJson respJson
		respuestaJson.Codigo = http.StatusOK
		respuestaJson.Respuesta = "Se contabilizaron correctamente " + strconv.Itoa(len(liquidaciones)) + " liquidaciones"
		framework.RespondJSON(w, http.StatusOK, respuestaJson)
	}

}

func marcarLiquidacionesComoContabilizadas(liquidaciones []structLiquidacion.Liquidacion, datosAsientoContableManual *monoliticComunication.StrDatosAsientoContableManual, db *gorm.DB) {
	for i := 0; i < len(liquidaciones); i++ {
		db.Model(&liquidaciones[i]).Update("Estacontabilizada", true)
		db.Model(&liquidaciones[i]).Update("Asientomanualtransaccionid", datosAsientoContableManual.Asientocontablemanualid)
		db.Model(&liquidaciones[i]).Update("Asientomanualnombre", datosAsientoContableManual.Asientocontablemanualnombre)
	}
}

func obtenerCuentasImportesYTipoDeGrillas(liquidacion structLiquidacion.Liquidacion, strCuentaImporteTipoGrillas *[]strCuentaImporteTipoGrilla, r *http.Request) {
	fmt.Println("Se obtienen las cuentas de la Liquidacion: " + strconv.Itoa(liquidacion.ID))
	var cuentaContable *int

	/*for i := 0; i < len(liquidacion.Importesremunerativos); i++ {

		importeremunerativo := liquidacion.Importesremunerativos[i]
		concepto := importeremunerativo.Concepto
		cuentaContable = concepto.CuentaContable
		importeUnitario := *importeremunerativo.Importeunitario

		cuentaImporteTipoGrilla := strCuentaImporteTipoGrilla{Cuentaid: *cuentaContable, Importecuenta: importeUnitario, Tipogrilla: 1}
		*strCuentaImporteTipoGrillas = append(*strCuentaImporteTipoGrillas, cuentaImporteTipoGrilla)

	}

	for j := 0; j < len(liquidacion.Importesnoremunerativos); j++ {

		importenoremunerativo := liquidacion.Importesnoremunerativos[j]
		concepto := importenoremunerativo.Concepto
		cuentaContable = concepto.CuentaContable
		importeUnitario := *importenoremunerativo.Importeunitario

		cuentaImporteTipoGrilla := strCuentaImporteTipoGrilla{Cuentaid: *cuentaContable, Importecuenta: importeUnitario, Tipogrilla: 2}
		*strCuentaImporteTipoGrillas = append(*strCuentaImporteTipoGrillas, cuentaImporteTipoGrilla)
	}

	for k := 0; k < len(liquidacion.Descuentos); k++ {

		descuento := liquidacion.Descuentos[k]
		concepto := descuento.Concepto
		cuentaContable = concepto.CuentaContable
		importeUnitario := *descuento.Importeunitario

		cuentaImporteTipoGrilla := strCuentaImporteTipoGrilla{Cuentaid: *cuentaContable, Importecuenta: importeUnitario, Tipogrilla: 3}
		*strCuentaImporteTipoGrillas = append(*strCuentaImporteTipoGrillas, cuentaImporteTipoGrilla)

	}

	for m := 0; m < len(liquidacion.Retenciones); m++ {
		retencion := liquidacion.Retenciones[m]
		concepto := retencion.Concepto
		cuentaContable = concepto.CuentaContable
		importeUnitario := *retencion.Importeunitario

		cuentaImporteTipoGrilla := strCuentaImporteTipoGrilla{Cuentaid: *cuentaContable, Importecuenta: importeUnitario, Tipogrilla: 4}
		*strCuentaImporteTipoGrillas = append(*strCuentaImporteTipoGrillas, cuentaImporteTipoGrilla)
	}

	for n := 0; n < len(liquidacion.Aportespatronales); n++ {

		aportepatronal := liquidacion.Aportespatronales[n]
		concepto := aportepatronal.Concepto
		cuentaContable = concepto.CuentaContable
		importeUnitario := *aportepatronal.Importeunitario

		cuentaImporteTipoGrilla := strCuentaImporteTipoGrilla{Cuentaid: *cuentaContable, Importecuenta: importeUnitario, Tipogrilla: 5}
		*strCuentaImporteTipoGrillas = append(*strCuentaImporteTipoGrillas, cuentaImporteTipoGrilla)
	}
	*/

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		item := liquidacion.Liquidacionitems[i]
		concepto := item.Concepto
		cuentaContable = concepto.CuentaContable
		importeUnitario := *item.Importeunitario

		cuentaImporteTipoGrilla := strCuentaImporteTipoGrilla{Cuentaid: *cuentaContable, Importecuenta: importeUnitario, Tipogrilla: *concepto.Tipoconceptoid}
		*strCuentaImporteTipoGrillas = append(*strCuentaImporteTipoGrillas, cuentaImporteTipoGrilla)
	}

	fmt.Println("Array strCuentaImporteTipoGrillas: ", *strCuentaImporteTipoGrillas)

}

func obtenerCuentasImportes(strCuentaImporteTipoGrillas []strCuentaImporteTipoGrilla, strCuentasImportes *[]strCuentaImporte) {

	sueldosYJornalesAPagar := -49
	cargasSocialesAPagar := -48

	for i := 0; i < len(strCuentaImporteTipoGrillas); i++ {
		cuentaImporteTipoGrilla := strCuentaImporteTipoGrillas[i]
		cuentaID := cuentaImporteTipoGrilla.Cuentaid
		importeUnitario := cuentaImporteTipoGrilla.Importecuenta
		tipoGrilla := cuentaImporteTipoGrilla.Tipogrilla

		switch tipoGrilla {
		case -1:
			cuentaImporteDebe := strCuentaImporte{Cuentaid: cuentaID, Importecuenta: importeUnitario}
			*strCuentasImportes = append(*strCuentasImportes, cuentaImporteDebe)

			cuentaImporteHaber := strCuentaImporte{Cuentaid: sueldosYJornalesAPagar, Importecuenta: importeUnitario * -1}
			*strCuentasImportes = append(*strCuentasImportes, cuentaImporteHaber)

		case -2:
			cuentaImporteDebe := strCuentaImporte{Cuentaid: cuentaID, Importecuenta: importeUnitario}
			*strCuentasImportes = append(*strCuentasImportes, cuentaImporteDebe)

			cuentaImporteHaber := strCuentaImporte{Cuentaid: sueldosYJornalesAPagar, Importecuenta: importeUnitario * -1}
			*strCuentasImportes = append(*strCuentasImportes, cuentaImporteHaber)

		case -3:
			cuentaImporteHaber := strCuentaImporte{Cuentaid: cuentaID, Importecuenta: importeUnitario * -1}
			*strCuentasImportes = append(*strCuentasImportes, cuentaImporteHaber)

			cuentaImporteDebe := strCuentaImporte{Cuentaid: sueldosYJornalesAPagar, Importecuenta: importeUnitario}
			*strCuentasImportes = append(*strCuentasImportes, cuentaImporteDebe)

		case -4:
			cuentaImporteHaber := strCuentaImporte{Cuentaid: cuentaID, Importecuenta: importeUnitario * -1}
			*strCuentasImportes = append(*strCuentasImportes, cuentaImporteHaber)

			cuentaImporteDebe := strCuentaImporte{Cuentaid: sueldosYJornalesAPagar, Importecuenta: importeUnitario}
			*strCuentasImportes = append(*strCuentasImportes, cuentaImporteDebe)

		case -5:
			cuentaImporteDebe := strCuentaImporte{Cuentaid: cuentaID, Importecuenta: importeUnitario}
			*strCuentasImportes = append(*strCuentasImportes, cuentaImporteDebe)

			cuentaImporteHaber := strCuentaImporte{Cuentaid: cargasSocialesAPagar, Importecuenta: importeUnitario * -1}
			*strCuentasImportes = append(*strCuentasImportes, cuentaImporteHaber)
		}
	}
}

func agruparCuentas(strCuentasImportes []strCuentaImporte, mapCuentasImportes map[int]float64) {
	for i := 0; i < len(strCuentasImportes); i++ {
		cuentaContable := strCuentasImportes[i].Cuentaid
		importeUnitario := strCuentasImportes[i].Importecuenta

		importe := mapCuentasImportes[cuentaContable]
		mapCuentasImportes[cuentaContable] = roundTo(importe+importeUnitario, 4)
	}
}

func obtenerCuentasImportesLiquidacion(mapCuentasImportes map[int]float64) []monoliticComunication.StrCuentaImporte {
	var arrayStrCuentaImporte []monoliticComunication.StrCuentaImporte

	for cuenta, importe := range mapCuentasImportes {
		var strcuentaimporte monoliticComunication.StrCuentaImporte
		strcuentaimporte.Cuentaid = cuenta
		strcuentaimporte.Importecuenta = importe
		arrayStrCuentaImporte = append(arrayStrCuentaImporte, strcuentaimporte)
	}

	return arrayStrCuentaImporte
}

/*func LiquidacionDesContabilizar(w http.ResponseWriter, r *http.Request) {
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

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

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

}*/
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

/*func descontabilizarLiquidaciones(w http.ResponseWriter, r *http.Request, liquidaciones []structLiquidacion.Liquidacion, asientomanualtransaccionid int, tokenAutenticacion *structAutenticacion.Security, respuestaDescontabilizar map[int]respJson, db *gorm.DB) (int, int) {

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

}*/

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

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

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

func LiquidacionDuplicarMasivo(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		decoder := json.NewDecoder(r.Body)

		var duplicarLiquidacionesData DuplicarLiquidaciones
		if err := decoder.Decode(&duplicarLiquidacionesData); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		var procesamientoMasivo ResultProcesamientoMasivo
		for index := 0; index < len(duplicarLiquidacionesData.Idstoreplicate); index++ {
			var liquidacionID = duplicarLiquidacionesData.Idstoreplicate[index]
			var liquidacion structLiquidacion.Liquidacion
			var procesamientoStatus ProcesamientoStatus

			//gorm:auto_preload se usa para que complete todos los struct con su informacion
			if err := db.Set("gorm:auto_preload", true).First(&liquidacion, "id = ?", liquidacionID).Error; gorm.IsRecordNotFoundError(err) {
				procesamientoStatus.Id = liquidacionID
				procesamientoStatus.Tipo = "ERROR"
				procesamientoStatus.Codigo = http.StatusNotFound
				procesamientoStatus.Mensaje = err.Error()
				procesamientoMasivo.Result = append(procesamientoMasivo.Result, procesamientoStatus)
			} else {
				/* se modifica liquidacion a duplicar */
				liquidacion.ID = 0
				liquidacion.Tipoid = duplicarLiquidacionesData.Liquidaciondefaultvalues.Tipoid
				liquidacion.Tipo = nil
				liquidacion.Fecha = duplicarLiquidacionesData.Liquidaciondefaultvalues.Fecha
				liquidacion.Fechaultimodepositoaportejubilatorio = duplicarLiquidacionesData.Liquidaciondefaultvalues.Fechaultimodepositoaportejubilatorio
				liquidacion.Fechaperiododepositado = duplicarLiquidacionesData.Liquidaciondefaultvalues.Fechaperiododepositado
				liquidacion.Fechaperiodoliquidacion = duplicarLiquidacionesData.Liquidaciondefaultvalues.Fechaperiodoliquidacion
				liquidacion.Estacontabilizada = false
				liquidacion.Asientomanualtransaccionid = 0
				liquidacion.Asientomanualnombre = ""

				for index := 0; index < len(liquidacion.Liquidacionitems); index++ {
					liquidacion.Liquidacionitems[index].ID = 0
					liquidacion.Liquidacionitems[index].CreatedAt = time.Time{}
					liquidacion.Liquidacionitems[index].UpdatedAt = time.Time{}
					liquidacion.Liquidacionitems[index].Liquidacionid = 0
				}

				/*for index := 0; index < len(liquidacion.Importesremunerativos); index++ {
					liquidacion.Importesremunerativos[index].ID = 0
					liquidacion.Importesremunerativos[index].CreatedAt = time.Time{}
					liquidacion.Importesremunerativos[index].UpdatedAt = time.Time{}
					liquidacion.Importesremunerativos[index].Liquidacionid = 0
				}
				for index := 0; index < len(liquidacion.Importesnoremunerativos); index++ {
					liquidacion.Importesnoremunerativos[index].ID = 0
					liquidacion.Importesnoremunerativos[index].CreatedAt = time.Time{}
					liquidacion.Importesnoremunerativos[index].UpdatedAt = time.Time{}
					liquidacion.Importesnoremunerativos[index].Liquidacionid = 0
				}
				for index := 0; index < len(liquidacion.Descuentos); index++ {
					liquidacion.Descuentos[index].ID = 0
					liquidacion.Descuentos[index].CreatedAt = time.Time{}
					liquidacion.Descuentos[index].UpdatedAt = time.Time{}
					liquidacion.Descuentos[index].Liquidacionid = 0
				}
				for index := 0; index < len(liquidacion.Retenciones); index++ {
					liquidacion.Retenciones[index].ID = 0
					liquidacion.Retenciones[index].CreatedAt = time.Time{}
					liquidacion.Retenciones[index].UpdatedAt = time.Time{}
					liquidacion.Retenciones[index].Liquidacionid = 0
				}
				for index := 0; index < len(liquidacion.Aportespatronales); index++ {
					liquidacion.Aportespatronales[index].ID = 0
					liquidacion.Aportespatronales[index].CreatedAt = time.Time{}
					liquidacion.Aportespatronales[index].UpdatedAt = time.Time{}
					liquidacion.Aportespatronales[index].Liquidacionid = 0
				}*/

				/*liquidacionJSON, _ := json.Marshal(liquidacion)
				fmt.Println(string(liquidacionJSON))*/

				/*decoder2 := json.NewDecoder(strings.NewReader(string(liquidacionJSON)))

				var liquidacion2 structLiquidacion.Liquidacion
				if err := decoder2.Decode(&liquidacion2); err != nil {
					framework.RespondError(w, http.StatusBadRequest, err.Error())
					return
				}*/

				if err := db.Create(&liquidacion).Error; err != nil {
					procesamientoStatus.Id = liquidacionID
					procesamientoStatus.Tipo = "ERROR"
					procesamientoStatus.Codigo = http.StatusInternalServerError
					procesamientoStatus.Mensaje = err.Error()
					procesamientoMasivo.Result = append(procesamientoMasivo.Result, procesamientoStatus)
				} else {
					/* se crea la duplicacion de la liquidacion correctamente */
					procesamientoStatus.Id = liquidacionID
					procesamientoStatus.Tipo = "SUCCESS"
					procesamientoStatus.Codigo = http.StatusOK
					procesamientoStatus.Mensaje = "Duplicado correctamente."
					procesamientoMasivo.Result = append(procesamientoMasivo.Result, procesamientoStatus)
				}
			}
		}

		framework.RespondJSON(w, http.StatusCreated, procesamientoMasivo)
	}
}

func LiquidacionAsientoManualDescontabilizar(w http.ResponseWriter, r *http.Request) {
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error:", err)
		}

		str := string(body)

		var datosAsientoContableManualBlanquear StrDatosAsientoContableManualBlanquear
		json.Unmarshal([]byte(str), &datosAsientoContableManualBlanquear)

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)
		defer conexionBD.CerrarDB(db)

		var tokenSecurityDecode []byte = []byte("blanquearasientomanualidenlasliquidacionesquecontabilizo")
		tokenSecurityEncode := base64.StdEncoding.EncodeToString(tokenSecurityDecode)

		if tokenSecurityEncode == datosAsientoContableManualBlanquear.Tokensecurityencode {
			asientocontablemanualid := datosAsientoContableManualBlanquear.Asientocontablemanualid
			var liquidaciones []structLiquidacion.Liquidacion
			db.Raw("UPDATE LIQUIDACION SET Asientomanualtransaccionid = 0, Asientomanualnombre = '', Estacontabilizada = false WHERE Asientomanualtransaccionid = " + strconv.Itoa(asientocontablemanualid)).Scan(&liquidaciones)

		} else {
			framework.RespondError(w, http.StatusInternalServerError, "Acceso denegado")
			return
		}

	}
	framework.RespondJSON(w, http.StatusCreated, "Liquidaciones descontabilizadas correctamente")

}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func roundTo(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func LiquidacionCalculoAutomatico(w http.ResponseWriter, r *http.Request) {
	var liquidacionCalculoAutomatico structLiquidacion.Liquidacion
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&liquidacionCalculoAutomatico); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		for i := 0; i < len(liquidacionCalculoAutomatico.Liquidacionitems); i++ {
			if liquidacionCalculoAutomatico.Liquidacionitems[i].DeletedAt == nil {
				concepto := *liquidacionCalculoAutomatico.Liquidacionitems[i].Concepto
				if concepto.Codigo == "IMPUESTO_GANANCIAS" || concepto.Codigo == "IMPUESTO_GANANCIAS_DEVOLUCION" {
					importeCalculoImpuestoGanancias := calculosAutomaticos.GetfgRetencionMes(&liquidacionCalculoAutomatico, db)
					*liquidacionCalculoAutomatico.Liquidacionitems[i].Importeunitario = roundTo(math.Abs(importeCalculoImpuestoGanancias), 4)

				} else {
					if concepto.Porcentaje != nil && concepto.Tipodecalculoid != nil {

						calculoAutomatico := calculosAutomaticos.NewCalculoAutomatico(&concepto, &liquidacionCalculoAutomatico)
						calculoAutomatico.Hacercalculoautomatico()
						*liquidacionCalculoAutomatico.Liquidacionitems[i].Importeunitario = roundTo(calculoAutomatico.GetImporteCalculado(), 4)
					}
				}
			}
		}

	}

	framework.RespondJSON(w, http.StatusOK, liquidacionCalculoAutomatico)

}

func LiquidacionCalculoAutomaticoConceptoId(w http.ResponseWriter, r *http.Request) {
	var liquidacionCalculoAutomatico structLiquidacion.Liquidacion
	var importeCalculado StrCalculoAutomaticoConceptoId
	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	if tokenValido {

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&liquidacionCalculoAutomatico); err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		defer r.Body.Close()

		tenant := apiclientautenticacion.ObtenerTenant(tokenAutenticacion)
		db := conexionBD.ObtenerDB(tenant)

		defer conexionBD.CerrarDB(db)

		params := mux.Vars(r)
		param_conceptoid, _ := strconv.ParseInt(params["id"], 10, 64)
		conceptoid := int(param_conceptoid)

		if conceptoid == 0 {
			framework.RespondError(w, http.StatusNotFound, framework.IdParametroVacio)
			return
		}
		var concepto structConcepto.Concepto

		//db.Set("gorm:auto_preload", true).First(&concepto, "id = ?", conceptoid)
		for i := 0; i < len(liquidacionCalculoAutomatico.Liquidacionitems); i++ {

			if liquidacionCalculoAutomatico.Liquidacionitems[i].Concepto.ID == conceptoid {
				concepto = *liquidacionCalculoAutomatico.Liquidacionitems[i].Concepto
				break
			}
		}

		importeCalculado.Conceptoid = &conceptoid
		if concepto.Codigo == "IMPUESTO_GANANCIAS" || concepto.Codigo == "IMPUESTO_GANANCIAS_DEVOLUCION" {
			importeCalculoImpuestoGanancias := roundTo(math.Abs(calculosAutomaticos.GetfgRetencionMes(&liquidacionCalculoAutomatico, db)), 4)
			importeCalculado = StrCalculoAutomaticoConceptoId{&conceptoid, &importeCalculoImpuestoGanancias}

		} else {
			if concepto.Porcentaje != nil && concepto.Tipodecalculoid != nil {
				calculoAutomatico := calculosAutomaticos.NewCalculoAutomatico(&concepto, &liquidacionCalculoAutomatico)
				calculoAutomatico.Hacercalculoautomatico()
				importeCalculadoConceptoID := roundTo(calculoAutomatico.GetImporteCalculado(), 4)
				importeCalculado = StrCalculoAutomaticoConceptoId{&conceptoid, &importeCalculadoConceptoID}
			}
		}

	}

	framework.RespondJSON(w, http.StatusOK, importeCalculado)

}
