package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"git-codecommit.us-east-1.amazonaws.com/v1/repos/sueldos-liquidacion/apiClientFormula"
	"git-codecommit.us-east-1.amazonaws.com/v1/repos/sueldos-liquidacion/calculosAutomaticos/Ganancias"
	"github.com/xubiosueldos/conexionBD/Novedad/structNovedad"
	"github.com/xubiosueldos/conexionBD/structGormModel"
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
	Conceptoid      *int                           `json:"conceptoid"`
	Importeunitario *float64                       `json:"importeunitario" `
	Acumuladores    []structLiquidacion.Acumulador `json:"acumuladores"`
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

		var where string

		if queries["fechadesde"] != nil || queries["fechahasta"] != nil {
			var p_fechadesde string = r.URL.Query()["fechadesde"][0] + " 00:00:00-03"
			var p_fechahasta string = r.URL.Query()["fechahasta"][0] + " 00:00:00-03"
			db.Set("gorm:auto_preload", true).Where("fecha BETWEEN ? AND ?", p_fechadesde, p_fechahasta).Find(&liquidaciones)
		} else {

			if queries["periododesde"] != nil {
				if where != "" {
					where += " AND "
				}
				where += fmt.Sprintf("to_char(fechaperiodoliquidacion, 'YYYY-MM') >= '%s'", queries["periododesde"][0])
			}

			if queries["periodohasta"] != nil {
				if where != "" {
					where += " AND "
				}
				where += fmt.Sprintf("to_char(fechaperiodoliquidacion, 'YYYY-MM') <= '%s'", queries["periodohasta"][0])
			}

			if queries["liquidaciontipoid"] != nil && queries["liquidaciontipoid"][0] != "0" {
				if where != "" {
					where += " AND "
				}
				where += fmt.Sprintf("tipoid = %s", queries["liquidaciontipoid"][0])
			}

			if where == "" {
				db.Set("gorm:auto_preload", true).Order("fechaperiodoliquidacion desc").Find(&liquidaciones)
			} else {
				db.Set("gorm:auto_preload", true).Order("fechaperiodoliquidacion desc").Where(where).Find(&liquidaciones)
			}
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

	autenticacion := r.Header.Get("Authorization")
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

		err := canInsertUpdate(liquidacion_data, db)
		if err != nil {
			framework.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		for i, liquidacionItem := range liquidacion_data.Liquidacionitems {
			if !liquidacionItem.Concepto.Eseditable && liquidacionItem.DeletedAt == nil {
				recalcularLiquidacionItem(&liquidacionItem, liquidacion_data, db, autenticacion)
				if roundTo(*liquidacion_data.Liquidacionitems[i].Importeunitario, 2) != roundTo(*liquidacionItem.Importeunitario, 2) {
					//framework.RespondError(w, http.StatusBadRequest, "El concepto " + *liquidacion_data.Liquidacionitems[i].Concepto.Nombre + " es no editable y su calculo automatico (" + fmt.Sprintf("%f" , roundTo(*liquidacionItem.Importeunitario, 2)) + ") no coincide con el valor actual " + fmt.Sprintf("%f", roundTo(*liquidacion_data.Liquidacionitems[i].Importeunitario,2)) + ". Intente recalcular.")
					framework.RespondError(w, http.StatusBadRequest, "Alguno de los importes de los conceptos no editables no coincide con el importe calculado automaticamente. Presione el botón Recalcular Conceptos Automaticos.")
					return
				}
			}
		}

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

func estaCargandoSacComoCorresponde(liquidacion structLiquidacion.Liquidacion, db *gorm.DB) error {
	var liquidacionNoPermitida structLiquidacion.Liquidacion
	if liquidacion.Fechaperiodoliquidacion.Month() == time.June && *liquidacion.Tipoid == liquidacionTipoSacID {

		db.Set("gorm:auto_preload", true).First(&liquidacionNoPermitida, "id != ? AND (tipoid = -1 OR tipoid = -3) AND legajoid = ? AND deleted_at is null AND to_char(fechaperiodoliquidacion, 'YYYY') = ? AND to_char(fechaperiodoliquidacion, 'MM') = '06'", liquidacion.ID, liquidacion.Legajoid, liquidacion.Fechaperiodoliquidacion.Format("2006"))

		if liquidacionNoPermitida.ID == 0 {
			return errors.New("Para cargar una liquidacion de tipo SAC en junio, primero debe cargar la liquidacion mensual/segunda quincena del legajo para ese mes")
		}
		return nil
	} else if liquidacion.Fechaperiodoliquidacion.Month() == time.December && *liquidacion.Tipoid == liquidacionTipoSacID {
		db.Set("gorm:auto_preload", true).First(&liquidacionNoPermitida, "id != ? AND (tipoid = -1 OR tipoid = -3) AND legajoid = ? AND deleted_at is null  AND to_char(fechaperiodoliquidacion, 'YYYY') = ? AND to_char(fechaperiodoliquidacion, 'MM') = '12'", liquidacion.ID, liquidacion.Legajoid, liquidacion.Fechaperiodoliquidacion.Format("2006"))

		if liquidacionNoPermitida.ID == 0 {
			return nil
		}

		return errors.New("Para cargar una liquidacion de tipo SAC en diciembre, no pueden existir liquidaciones de tipo mensual/segunda quincena en ese mismo mes")
	}

	return nil
}

const (
	liquidacionTipoMensualID         = -1
	liquidacionTipoPrimeraQuincenaID = -2
	liquidacionTipoSegundaQuincenaID = -3
	liquidacionTipoSacID             = -5
)

func LiquidacionUpdate(w http.ResponseWriter, r *http.Request) {

	tokenValido, tokenAutenticacion := apiclientautenticacion.CheckTokenValido(w, r)
	autenticacion := r.Header.Get("Authorization")
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

		db2 := conexionBD.ObtenerDB(tenant)
		defer conexionBD.CerrarDB(db2)

		if !esUltimaLiquidacionDelAño(p_liquidacionid, db) {
			framework.RespondError(w, http.StatusBadRequest, "Para modificar esta liquidacion primero debe eliminar la liquidacion siguiente")
			return
		}

		if !liquidacionContabilizada(p_liquidacionid, db) {
			decoder := json.NewDecoder(r.Body)

			var liquidacion_data structLiquidacion.Liquidacion

			if err := decoder.Decode(&liquidacion_data); err != nil {
				framework.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}
			defer r.Body.Close()

			err := canInsertUpdate(liquidacion_data, db)

			if err != nil {
				framework.RespondError(w, http.StatusBadRequest, err.Error())
				return
			}

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
				defer tx.Rollback()

				//Actualizo los Calculos necesarios y refresco los acumuladores de los mismos
				for i, liquidacionItem := range liquidacion_data.Liquidacionitems {

					db.Set("gorm:auto_preload", true).First(liquidacionItem.Concepto, "id = ?", liquidacionItem.Conceptoid)
					if !liquidacionItem.Concepto.Eseditable && liquidacionItem.DeletedAt == nil {
						recalcularLiquidacionItem(&liquidacionItem, liquidacion_data, db2, autenticacion)
						if roundTo(*liquidacion_data.Liquidacionitems[i].Importeunitario, 2) != roundTo(*liquidacionItem.Importeunitario, 2) {
							//framework.RespondError(w, http.StatusBadRequest, "El concepto " + *liquidacion_data.Liquidacionitems[i].Concepto.Nombre + " es no editable y su calculo automatico (" + fmt.Sprintf("%f" ,roundTo(*liquidacionItem.Importeunitario,2)) + ") no coincide con el valor actual " + fmt.Sprintf("%f", roundTo(*liquidacion_data.Liquidacionitems[i].Importeunitario,2)) + ". Intente recalcular.")
							framework.RespondError(w, http.StatusBadRequest, "Alguno de los importes de los conceptos no editables no coincide con el importe calculado automaticamente. Presione el botón Recalcular Conceptos Automaticos.")
							return
						}
					}

					if liquidacionItem.Concepto.Esganancias {
						for _, acumulador := range liquidacionItem.Acumuladores {
							acumulador.ID = 0
						}
						if err := tx.Model(structLiquidacion.Acumulador{}).Unscoped().Where("liquidacionitemid = ?", liquidacionItem.ID).Delete(structLiquidacion.Acumulador{}).Error; err != nil {
							framework.RespondError(w, http.StatusInternalServerError, err.Error())
							return
						}
					}
				}

				//modifico el legajo de acuerdo a lo enviado en el json
				if err := tx.Save(&liquidacion_data).Error; err != nil {
					framework.RespondError(w, http.StatusInternalServerError, err.Error())
					return
				}

				if err := tx.Model(structLiquidacion.Liquidacionitem{}).Unscoped().Where("liquidacionid = ? AND deleted_at is not null", liquidacionid).Delete(structLiquidacion.Liquidacionitem{}).Error; err != nil {
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

func esUltimaLiquidacionDelAño(liquidacionid int, db *gorm.DB) bool {
	var liquidacionActual structLiquidacion.Liquidacion
	var liquidacionMasReciente structLiquidacion.Liquidacion
	db.First(&liquidacionActual, "id = "+strconv.Itoa(liquidacionid))
	db.Order("to_number(to_char(fechaperiodoliquidacion, 'MM'),'99') desc, fecha desc, created_at desc").Set("gorm:auto_preload", true).First(&liquidacionMasReciente, "to_char(fechaperiodoliquidacion, 'YYYY') = ? AND legajoid = ?", strconv.Itoa(liquidacionActual.Fechaperiodoliquidacion.Year()), *liquidacionActual.Legajoid)
	return liquidacionActual.ID == liquidacionMasReciente.ID
}

func recalcularLiquidacionItem(liquidacionItem *structLiquidacion.Liquidacionitem, liquidacion structLiquidacion.Liquidacion, db *gorm.DB, autenticacion string) {
	solucionCalculo := calcularConcepto(liquidacionItem.Concepto.ID, &liquidacion, liquidacionItem, db, autenticacion)
	liquidacionItem.Importeunitario = solucionCalculo.Importeunitario
	liquidacionItem.Acumuladores = solucionCalculo.Acumuladores
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
					obtenerCuentasImportes(liquidaciones[i], &strCuentasImportes, r)
				}
				agruparCuentas(strCuentasImportes, mapCuentasImportes)
				eliminarCuentasImporteCero(mapCuentasImportes)

			} else {
				framework.RespondError(w, http.StatusNotFound, framework.Seleccionaronliquidacionescontabilizadas)
				return
			}
		}

		generarAsientoManualDesdeMonolitico(w, r, liquidaciones, mapCuentasImportes, tokenAutenticacion, descripcion_asiento, fecha_asiento, 0, db)

	}

}

func eliminarCuentasImporteCero(mapCuentasImportes map[int]float64) {

	for cuentaContableid, importeUnitario := range mapCuentasImportes {
		if importeUnitario == 0 {
			delete(mapCuentasImportes, cuentaContableid)
		}
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

func obtenerCuentasImportes(liquidacion structLiquidacion.Liquidacion, strCuentasImportes *[]strCuentaImporte, r *http.Request) {
	fmt.Println("Se obtienen las cuentas de la Liquidacion: " + strconv.Itoa(liquidacion.ID))

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		item := liquidacion.Liquidacionitems[i]
		concepto := item.Concepto
		cuentaContableDebe := concepto.CuentaContable
		cuentaContableHaber := concepto.Cuentacontablepasivoid
		importeUnitario := *item.Importeunitario

		cuentaImporte := strCuentaImporte{Cuentaid: *cuentaContableDebe, Importecuenta: importeUnitario}
		cuentaImportePasivo := strCuentaImporte{Cuentaid: *cuentaContableHaber, Importecuenta: importeUnitario * -1}
		*strCuentasImportes = append(*strCuentasImportes, cuentaImporte)
		*strCuentasImportes = append(*strCuentasImportes, cuentaImportePasivo)
	}

	fmt.Println("Array strCuentaImporte: ", *strCuentasImportes)

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
			var liquidacionTipo structLiquidacion.Liquidaciontipo
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
				if err := db.Set("gorm:auto_preload", true).First(&liquidacionTipo, "id = ?", *liquidacion.Tipoid).Error; gorm.IsRecordNotFoundError(err) {
					procesamientoStatus.Id = *liquidacion.Tipoid
					procesamientoStatus.Tipo = "ERROR"
					procesamientoStatus.Codigo = http.StatusNotFound
					procesamientoStatus.Mensaje = err.Error()
					procesamientoMasivo.Result = append(procesamientoMasivo.Result, procesamientoStatus)
				}
				liquidacion.Tipo = &liquidacionTipo
				liquidacion.Fecha = duplicarLiquidacionesData.Liquidaciondefaultvalues.Fecha
				liquidacion.Fechaultimodepositoaportejubilatorio = duplicarLiquidacionesData.Liquidaciondefaultvalues.Fechaultimodepositoaportejubilatorio
				liquidacion.Fechaperiododepositado = duplicarLiquidacionesData.Liquidaciondefaultvalues.Fechaperiododepositado
				liquidacion.Fechaperiodoliquidacion = duplicarLiquidacionesData.Liquidaciondefaultvalues.Fechaperiodoliquidacion
				liquidacion.Estacontabilizada = false
				liquidacion.Asientomanualtransaccionid = 0
				liquidacion.Asientomanualnombre = ""
				liquidacion.Cantidaddiastrabajados = duplicarLiquidacionesData.Liquidaciondefaultvalues.Cantidaddiastrabajados
				liquidacion.Situacionrevistauno = duplicarLiquidacionesData.Liquidaciondefaultvalues.Situacionrevistauno
				liquidacion.Fechasituacionrevistauno = &duplicarLiquidacionesData.Liquidaciondefaultvalues.Fechaperiodoliquidacion
				if duplicarLiquidacionesData.Liquidaciondefaultvalues.Situacionrevistados != nil {
					liquidacion.Situacionrevistados = duplicarLiquidacionesData.Liquidaciondefaultvalues.Situacionrevistados
					liquidacion.Fechasituacionrevistados = &duplicarLiquidacionesData.Liquidaciondefaultvalues.Fechaperiodoliquidacion
				}
				if duplicarLiquidacionesData.Liquidaciondefaultvalues.Situacionrevistatres != nil {
					liquidacion.Situacionrevistatres = duplicarLiquidacionesData.Liquidaciondefaultvalues.Situacionrevistatres
					liquidacion.Fechasituacionrevistatres = &duplicarLiquidacionesData.Liquidaciondefaultvalues.Fechaperiodoliquidacion
				}

				var liquidacionItems []structLiquidacion.Liquidacionitem
				for index := 0; index < len(liquidacion.Liquidacionitems); index++ {
					var concepto structConcepto.Concepto
					if err := db.Set("gorm:auto_preload", true).First(&concepto, "id = ?", *liquidacion.Liquidacionitems[index].Conceptoid).Error; gorm.IsRecordNotFoundError(err) {
						procesamientoStatus.Id = *liquidacion.Liquidacionitems[index].Conceptoid
						procesamientoStatus.Tipo = "ERROR"
						procesamientoStatus.Codigo = http.StatusNotFound
						procesamientoStatus.Mensaje = err.Error()
						procesamientoMasivo.Result = append(procesamientoMasivo.Result, procesamientoStatus)
					}
					if !concepto.Esnovedad {
						liquidacion.Liquidacionitems[index].ID = 0
						liquidacion.Liquidacionitems[index].CreatedAt = time.Time{}
						liquidacion.Liquidacionitems[index].UpdatedAt = time.Time{}
						liquidacion.Liquidacionitems[index].Liquidacionid = 0
						liquidacion.Liquidacionitems[index].Acumuladores = nil
						liquidacionItems = append(liquidacionItems, liquidacion.Liquidacionitems[index])
					}
				}
				agregarNovedades(&liquidacionItems, liquidacion.Fechaperiodoliquidacion, liquidacionTipo, db, *liquidacion.Legajoid)
				liquidacion.Liquidacionitems = liquidacionItems

				for index := 0; index < len(liquidacion.Liquidacionitems); index++ {
					if !liquidacion.Liquidacionitems[index].Concepto.Eseditable {
						recalcularLiquidacionItem(&liquidacion.Liquidacionitems[index], liquidacion, db, r.Header.Get("Authorization"))
					}
				}

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

func agregarNovedades(liquidacionItems *[]structLiquidacion.Liquidacionitem, fechaperiodoliquidacion time.Time, tipo structLiquidacion.Liquidaciontipo, db *gorm.DB, legajoid int) {
	var novedades []structNovedad.Novedad
	switch tipo.ID {
	case liquidacionTipoMensualID:
		db.Set("gorm:auto_preload", true).Find(&novedades, "to_char(fecha, 'YYYY') = ? AND to_char(fecha, 'MM') = ? AND legajoid = ?", fechaperiodoliquidacion.Year(), fechaperiodoliquidacion.Format("01"), legajoid)
		break
	case liquidacionTipoPrimeraQuincenaID:
		db.Set("gorm:auto_preload", true).Find(&novedades, "to_char(fecha, 'YYYY') = ? AND to_char(fecha, 'MM') = ? AND to_char(fecha, 'DD') <= '15' AND legajoid = ?", fechaperiodoliquidacion.Year(), fechaperiodoliquidacion.Format("01"), legajoid)
		break
	case liquidacionTipoSegundaQuincenaID:
		db.Set("gorm:auto_preload", true).Find(&novedades, "to_char(fecha, 'YYYY') = ? AND to_char(fecha, 'MM') = ? AND to_char(fecha, 'DD') > '15' AND legajoid = ?", fechaperiodoliquidacion.Year(), fechaperiodoliquidacion.Format("01"), legajoid)
		break
	}

	for _, novedad := range novedades {
		var novedadImporte *float64
		if *novedad.Concepto.Tipocalculoautomaticoid == -3 {

		} else {
			novedadImporteFloat := float64(*novedad.Importe)
			novedadImporte = &novedadImporteFloat
		}

		liquidacionitem := structLiquidacion.Liquidacionitem{
			GormModel: structGormModel.GormModel{
				ID:        0,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: nil,
			},
			Concepto:        novedad.Concepto,
			Conceptoid:      novedad.Conceptoid,
			Importeunitario: novedadImporte,
			Liquidacionid:   0,
			Cantidad:        novedad.Cantidad,
			Acumuladores:    nil,
		}

		*liquidacionItems = append(*liquidacionItems, liquidacionitem)

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
	autenticacion := r.Header.Get("Authorization")
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

		defer func() {
			if r := recover(); r != nil {
				err := r.(error)
				framework.RespondError(w, http.StatusBadRequest, err.Error())
			}
		}()

		for i, liquidacionitem := range liquidacionCalculoAutomatico.Liquidacionitems {
			if liquidacionitem.DeletedAt == nil {
				concepto := *liquidacionitem.Concepto

				liquidacionCalculoAutomaticoCopia := liquidacionCalculoAutomatico
				resultado := calcularConcepto(concepto.ID, &liquidacionCalculoAutomaticoCopia, &liquidacionitem, db, autenticacion)

				if resultado != nil {
					liquidacionCalculoAutomatico.Liquidacionitems[i].Importeunitario = resultado.Importeunitario
					liquidacionCalculoAutomatico.Liquidacionitems[i].Acumuladores = resultado.Acumuladores
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

	autenticacion := r.Header.Get("Authorization")
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

		defer func() {
			if r := recover(); r != nil {
				err := r.(error)
				framework.RespondError(w, http.StatusBadRequest, err.Error())
			}
		}()

		var liquidacionitem *structLiquidacion.Liquidacionitem

		for i := 0; i < len(liquidacionCalculoAutomatico.Liquidacionitems); i++ {

			if liquidacionCalculoAutomatico.Liquidacionitems[i].Concepto.ID == conceptoid {
				liquidacionitem = &liquidacionCalculoAutomatico.Liquidacionitems[i]
				break
			}
		}

		calculo := calcularConcepto(conceptoid, &liquidacionCalculoAutomatico, liquidacionitem, db, autenticacion)

		if calculo != nil {
			importeCalculado = *calculo
		}

	}

	framework.RespondJSON(w, http.StatusOK, importeCalculado)

}

func calcularConcepto(conceptoid int, liquidacionCalculoAutomatico *structLiquidacion.Liquidacion, liquidacionitem *structLiquidacion.Liquidacionitem, db *gorm.DB, autenticacion string) *StrCalculoAutomaticoConceptoId {

	importeCalculado := StrCalculoAutomaticoConceptoId{}
	//db.Set("gorm:auto_preload", true).First(&concepto, "id = ?", conceptoid)
	importeCalculado.Conceptoid = &conceptoid

	var concepto structConcepto.Concepto
	var Tipocalculoautomatico structConcepto.Tipocalculoautomatico

	db.First(&concepto, "id = ?", conceptoid)

	if concepto.ID == 0 || liquidacionitem == nil {
		panic(errors.New("Error al obtener el concepto de la liquidacion"))
	}

	if concepto.Tipocalculoautomaticoid != nil && concepto.Tipocalculoautomatico == nil {
		ID := *concepto.Tipocalculoautomaticoid
		db.Set("gorm:auto_preload", true).First(&Tipocalculoautomatico, ID)
		concepto.Tipocalculoautomatico = &Tipocalculoautomatico
	}

	if concepto.Tipocalculoautomatico == nil {
		concepto.Tipocalculoautomatico = &structConcepto.Tipocalculoautomatico{}
		concepto.Tipocalculoautomatico.Codigo = "NO_APLICA"
	}

	if concepto.Tipocalculoautomatico.Codigo == "NO_APLICA" {
		return nil
	}

	if concepto.Tipocalculoautomatico.Codigo == "FORMULA" {

		if *concepto.Formulanombre == "ImpuestoALasGanancias" {
			importeCalculado = ImpuestoALasGanancias(concepto, liquidacionCalculoAutomatico, liquidacionitem, db)
		} else if *concepto.Formulanombre == "ImpuestoALasGananciasDevolucion" {
			importeCalculado = ImpuestoALasGananciasDevolucion(concepto, liquidacionCalculoAutomatico, liquidacionitem, db)
		} else {
			//CODIGO PARA EJECUTAR LAS FORMULAS
			resultadoFormula, err := apiClientFormula.ExecuteFormulaLiquidacion(autenticacion, liquidacionCalculoAutomatico, *concepto.Formulanombre, &concepto, liquidacionitem)
			if err != nil {
				panic(err)
			}

			importeCalculado.Importeunitario = &resultadoFormula
		}

	} else if concepto.Tipocalculoautomatico.Codigo == "PORCENTAJE" {
		if concepto.Porcentaje != nil && concepto.Tipodecalculoid != nil {
			calculoAutomatico := calculosAutomaticos.NewCalculoAutomatico(&concepto, liquidacionCalculoAutomatico)
			calculoAutomatico.Hacercalculoautomatico()
			importeCalculadoConceptoID := roundTo(calculoAutomatico.GetImporteCalculado(), 4)
			importeCalculado.Importeunitario = &importeCalculadoConceptoID
		}
	}

	return &importeCalculado
}

func ImpuestoALasGanancias(concepto structConcepto.Concepto, liquidacionCalculoAutomatico *structLiquidacion.Liquidacion, liquidacionitem *structLiquidacion.Liquidacionitem, db *gorm.DB) StrCalculoAutomaticoConceptoId {
	importeCalculado := StrCalculoAutomaticoConceptoId{}

	liquidacionitem.Acumuladores = nil

	if liquidacionCalculoAutomatico.Tipo.Codigo != "PRIMER_QUINCENA" && liquidacionCalculoAutomatico.Tipo.Codigo != "VACACIONES" {
		importeCalculoImpuestoGanancias := roundTo((&Ganancias.CalculoGanancias{liquidacionitem, liquidacionCalculoAutomatico, db, true}).Calculate(), 2)
		importeCalculado.Importeunitario = &importeCalculoImpuestoGanancias
	} else {
		panic(errors.New("La Liquidación de tipo Primer Quincena o Vacaciones no permite los conceptos de Impuesto a las Ganancias"))
	}

	importeCalculado.Acumuladores = liquidacionitem.Acumuladores
	importeCalculado.Conceptoid = &concepto.ID

	return importeCalculado
}

func ImpuestoALasGananciasDevolucion(concepto structConcepto.Concepto, liquidacionCalculoAutomatico *structLiquidacion.Liquidacion, liquidacionitem *structLiquidacion.Liquidacionitem, db *gorm.DB) StrCalculoAutomaticoConceptoId {
	importeCalculado := ImpuestoALasGanancias(concepto, liquidacionCalculoAutomatico, liquidacionitem, db)
	importeFinal := (*importeCalculado.Importeunitario) * -1
	importeCalculado.Importeunitario = &importeFinal
	return importeCalculado
}

func canInsertUpdate(liquidacion structLiquidacion.Liquidacion, db *gorm.DB) error {

	fechaperiodoliquidacionigualfechasituacionrevistados := true
	fechaperiodoliquidacionigualfechasituacionrevistatres := true

	mesLiquidacion := liquidacion.Fechaperiodoliquidacion.Format("01")
	anioLiquidacion := liquidacion.Fechaperiodoliquidacion.Year()

	mesSituacionRevistaUno := liquidacion.Fechasituacionrevistauno.Format("01")
	anioSituacionRevistaUno := liquidacion.Fechasituacionrevistauno.Year()

	fechaperiodoliquidacionigualfechasituacionrevistauno := (mesLiquidacion == mesSituacionRevistaUno && anioLiquidacion == anioSituacionRevistaUno)

	if liquidacion.Fechasituacionrevistados != nil {
		mesSituacionRevistaDos := liquidacion.Fechasituacionrevistados.Format("01")
		anioSituacionRevistaDos := liquidacion.Fechasituacionrevistados.Year()

		fechaperiodoliquidacionigualfechasituacionrevistados = (mesLiquidacion == mesSituacionRevistaDos && anioLiquidacion == anioSituacionRevistaDos)

	}

	if liquidacion.Fechasituacionrevistatres != nil {
		mesSituacionRevistaTres := liquidacion.Fechasituacionrevistatres.Format("01")
		anioSituacionRevistaTres := liquidacion.Fechasituacionrevistatres.Year()

		fechaperiodoliquidacionigualfechasituacionrevistatres = (mesLiquidacion == mesSituacionRevistaTres && anioLiquidacion == anioSituacionRevistaTres)

	}

	if !(fechaperiodoliquidacionigualfechasituacionrevistauno && fechaperiodoliquidacionigualfechasituacionrevistados && fechaperiodoliquidacionigualfechasituacionrevistatres) {
		return errors.New("La Fecha Desde de Situación Revista debe pertenecer al Periodo Liquidación")
	}

	var cantidadItemsConGanancias = 0
	for _, liquidacionitem := range liquidacion.Liquidacionitems {
		db.Set("gorm:auto_preload", true).First(&liquidacionitem.Concepto, "id = ?", liquidacionitem.Conceptoid)
		if liquidacionitem.DeletedAt == nil && liquidacionitem.Concepto.Esganancias {
			if *liquidacionitem.Importeunitario < 0 {
				return errors.New("El concepto de impuesto a las ganancias no puede tener importe negativo.")
			}
			cantidadItemsConGanancias++
		}
	}

	if cantidadItemsConGanancias > 1 {
		return errors.New("Solo se puede utilizar un único concepto de Impuesto a las Ganancias por liquidación")
	}

	if (liquidacion.Tipo.Codigo == "PRIMER_QUINCENA" || liquidacion.Tipo.Codigo == "VACACIONES") && cantidadItemsConGanancias == 1 {
		return errors.New("La Liquidación de tipo Primer Quincena o Vacaciones no permite los conceptos de Impuesto a las Ganancias")
	}

	err := estaCargandoSacComoCorresponde(liquidacion, db)

	if err != nil {
		return err
	}

	return nil
}
