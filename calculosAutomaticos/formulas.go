package calculosAutomaticos

import (
	"fmt"
	"strconv"
	"time"

	s "strings"

	"github.com/jinzhu/gorm"
	"github.com/xubiosueldos/conexionBD/Concepto/structConcepto"
	"github.com/xubiosueldos/conexionBD/Liquidacion/structLiquidacion"
	"github.com/xubiosueldos/conexionBD/Siradig/structSiradig"
)

func getRemuneracionBruta(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("REMUNERACION_BRUTA", liquidacion, db)
	return importeTotal
}

func getRemuneracionNoHabitual(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("RETRIBUCIONES_NO_HABITUALES", liquidacion, db)
	return importeTotal
}

func getSacPrimerCuota(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	correspondePrimerSemetre := getMes(&liquidacion.Fechaperiodoliquidacion) <= 6
	importeTotal := getSacCuotas(liquidacion, correspondePrimerSemetre, db)

	return importeTotal

}

func getSacSegundaCuota(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	correspondeSegundoSemetre := getMes(&liquidacion.Fechaperiodoliquidacion) > 6
	importeTotal := getSacCuotas(liquidacion, correspondeSegundoSemetre, db)

	return importeTotal
}

func getHorasExtrasGravadas(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("HORAS_EXTRAS_REMUNERACION_GRAVADA", liquidacion, db)
	return importeTotal
}

func getMovilidadYViaticosGravada(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("MOVILIDAD_Y_VIATICOS_REMUNERACION_GRAVADA", liquidacion, db)
	return importeTotal
}

func getMaterialDidacticoPersonalDocenteRemuneracion(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("MATERIAL_DIDACTICO_PERSONAL_DOCENTE_REMUNERACION_GRAVADA", liquidacion, db)
	return importeTotal
}

func getRemuneracionBrutaOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteGananciasOtroEmpleoSiradig(liquidacion, "importegananciasbrutas", db)
	return importeTotal
}

func getRemuneracionNoHabitualOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteGananciasOtroEmpleoSiradig(liquidacion, "importeretribucionesnohabituales", db)
	return importeTotal
}

func getSacPrimerCuotaOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0
	if getMes(&liquidacion.Fechaperiodoliquidacion) <= 6 {
		importeTotal = getImporteGananciasOtroEmpleoSiradig(liquidacion, "sac", db)
	}
	return importeTotal
}

func getSacSegundaCuotaOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0
	if getMes(&liquidacion.Fechaperiodoliquidacion) > 6 {
		importeTotal = getImporteGananciasOtroEmpleoSiradig(liquidacion, "sac", db)
	}
	return importeTotal
}

func getHorasExtrasGravadasOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteGananciasOtroEmpleoSiradig(liquidacion, "importehorasextrasgravadas", db)
	return importeTotal
}

func getMovilidadYViaticosGravadaOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteGananciasOtroEmpleoSiradig(liquidacion, "gastosmovilidad", db)
	return importeTotal
}

func getMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteGananciasOtroEmpleoSiradig(liquidacion, "materialdidactico", db)
	return importeTotal
}

func getSubtotalIngresos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {

	var arraySubtotalIngresos []float64
	var subtotalIngresos float64

	arraySubtotalIngresos = append(arraySubtotalIngresos, getRemuneracionBruta(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getRemuneracionNoHabitual(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getSacPrimerCuota(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getSacSegundaCuota(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getHorasExtrasGravadas(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getMovilidadYViaticosGravada(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getMaterialDidacticoPersonalDocenteRemuneracion(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getRemuneracionBrutaOtrosEmpleos(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getRemuneracionNoHabitualOtrosEmpleos(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getSacPrimerCuotaOtrosEmpleos(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getSacSegundaCuotaOtrosEmpleos(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getHorasExtrasGravadasOtrosEmpleos(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getMovilidadYViaticosGravadaOtrosEmpleos(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos(liquidacion, db))

	subtotalIngresos = Sum(arraySubtotalIngresos)
	return subtotalIngresos
}

func getAportesJubilatoriosRetirosPensionesOSubsidios(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("APORTES_JUBILATORIOS_RETIROS_PENSIONES_O_SUBSIDIOS", liquidacion, db)
	return importeTotal
}

func getAportesObraSocial(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("APORTES_OBRA_SOCIAL", liquidacion, db)
	return importeTotal
}

func getCuotaSindical(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("CUOTA_SINDICAL", liquidacion, db)
	return importeTotal
}

func getDescuentosObligatoriosPorLeyNacionalProvincialOMunicipal(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSegunTipoImpuestoGanancias("DESCUENTOS_OBLIGATORIOS_POR_LEY_NACIONAL_PROVINCIAL_MUNICIPAL", liquidacion, db)
	return importeTotal
}

func getGastosMovilidadViaticosAbonadosPorElEmpleador(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arrayGastosMovilidad []float64
	var importeTotal float64

	arrayGastosMovilidad = append(arrayGastosMovilidad, getMovilidadYViaticosGravada(liquidacion, db))
	arrayGastosMovilidad = append(arrayGastosMovilidad, getMovilidadYViaticosGravadaOtrosEmpleos(liquidacion, db))
	arrayGastosMovilidad = append(arrayGastosMovilidad, getMaterialDidacticoPersonalDocenteRemuneracion(liquidacion, db))
	arrayGastosMovilidad = append(arrayGastosMovilidad, getMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos(liquidacion, db))

	importeTotal = Sum(arrayGastosMovilidad)
	return importeTotal

}

func getAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteGananciasOtroEmpleoSiradig(liquidacion, "aporteseguridadsocial", db)
	return importeTotal
}

func getAportesObraSocialOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteGananciasOtroEmpleoSiradig(liquidacion, "aporteobrasocial", db)
	return importeTotal
}

func getCuotaSindicalOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteGananciasOtroEmpleoSiradig(liquidacion, "aportesindical", db)
	return importeTotal
}

func getPrimasDeSeguroParaCasoDeMuerte(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	return 0
}

func getSeguroMuerteMixtosSujetosAlControlSSN(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	return 0
}

func getSegurosRetirosPrivadosSujetosAlControlSSN(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	return 0
}

func getImporteTotalSiradigSegunTipoGrilla(liquidacion *structLiquidacion.Liquidacion, columnadeducciondesgravacionsiradig string, tipodeducciondesgravacionsiradig string, nombretablasiradig string, db *gorm.DB) float64 {
	var importeTotal float64
	fechaliquidacion := liquidacion.Fechaperiodoliquidacion.Format("2006-01-02")

	sql := "SELECT SUM(" + columnadeducciondesgravacionsiradig + ") FROM " + nombretablasiradig + " ts INNER JOIN siradigtipogrilla stg ON stg.id = ts.siradigtipogrillaid WHERE '" + fechaliquidacion + "' >= mes AND stg.codigo = '" + tipodeducciondesgravacionsiradig + "'"

	db.Raw(sql).Row().Scan(&importeTotal)

	return importeTotal
}

func getGastosSepelio(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "GASTOS_DE_SEPELIO", "deducciondesgravacionsiradig", db)
	importeTope := getValorFijoImpuestoGanancia(liquidacion, "topemaximodescuento", "topesepelio", db)
	importeTotal = getImporteTotalTope(importeTotal, importeTope)
	return importeTotal
}

func getGastosAmortizacionEInteresesRodadoCorredoresViajantesComercio(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "GASTOS_DE_REPRESENTACION_E_INTERESES_DE_CORREDORES_Y_VIAJANTES_DE_COMERCIO", "deducciondesgravacionsiradig", db)
	return importeTotal
}

func getInteresesCreditosHipotecarios(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "INTERESES_PRESTAMO_HIPOTECARIO", "deducciondesgravacionsiradig", db)
	importeTope := getValorFijoImpuestoGanancia(liquidacion, "topemaximodescuento", "topehipotecarios", db)
	importeTotal = getImporteTotalTope(importeTotal, importeTope)
	return importeTotal
}

func getAportesCapSocFondoRiesgoSociosProtectoresSGR(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSiradigSegunTipoGrilla(liquidacion, "Montoreintegrar + Montoreintegrar3", "REINTEGRO_DE_APORTES_DE_SOCIOS_PROTECTORES_A_SOCIEDADES_DE_GARANTIA_RECIPROCA", "ajustesiradig", db)
	return importeTotal
}

func getAlquileresInmueblesDestinadosASuCasaHabitacion(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "ALQUILER_INMUEBLES_DESTINADOS_A_CASA_HABITACION", "deducciondesgravacionsiradig", db)
	importeTope := 11 / 0.4 /*es el 40% de MNI(40)*/
	importeTotal = getImporteTotalTope(importeTotal, importeTope)
	return importeTotal
}

func getEmpleadosServicioDomestico(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSiradigSegunTipoGrilla(liquidacion, "contribucion + retribucion", "DEDUCCION_DEL_PERSONAL_DOMESTICO", "deducciondesgravacionsiradig", db)
	importeTope := float64(11) /*es el MNI(40)*/
	importeTotal = getImporteTotalTope(importeTotal, importeTope)
	return importeTotal
}

func getIndumentariaEquipamientoCaracterObligatorio(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "GASTOS_ADQUISICION_INDUMENTARIA_Y_EQUIPAMIENTO_PARA_USO_EXCLUSIVO_EN_EL_LUGAR_DE_TRABAJO", "deducciondesgravacionsiradig", db)
	return importeTotal
}

func getOtrasDeducciones(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "OTRAS", "deducciondesgravacionsiradig", db)
	return importeTotal
}

func getSubtotal(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arraySubtotal []float64
	var subTotal float64

	arraySubtotal = append(arraySubtotal, getAportesJubilatoriosRetirosPensionesOSubsidios(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getAportesObraSocial(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getCuotaSindical(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getDescuentosObligatoriosPorLeyNacionalProvincialOMunicipal(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getGastosMovilidadViaticosAbonadosPorElEmpleador(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getAportesObraSocialOtrosEmpleos(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getCuotaSindicalOtrosEmpleos(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getPrimasDeSeguroParaCasoDeMuerte(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getSeguroMuerteMixtosSujetosAlControlSSN(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getSegurosRetirosPrivadosSujetosAlControlSSN(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getGastosSepelio(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getGastosAmortizacionEInteresesRodadoCorredoresViajantesComercio(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getInteresesCreditosHipotecarios(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getAportesCapSocFondoRiesgoSociosProtectoresSGR(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getAlquileresInmueblesDestinadosASuCasaHabitacion(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getEmpleadosServicioDomestico(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getIndumentariaEquipamientoCaracterObligatorio(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getOtrasDeducciones(liquidacion, db))

	subTotal = getSubtotalIngresos(liquidacion, db) - Sum(arraySubtotal)

	return subTotal
}

func getCuotaMedicoAsistencial(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "CUOTA_MEDICA_ASISTENCIAL", "deducciondesgravacionsiradig", db)
	importeTope := getSubtotal(liquidacion, db) * 0.05 //5% de Subtotal
	importeTotal = getImporteTotalTope(importeTotal, importeTope)
	return importeTotal
}

func getDonacionFiscosNacProvMunArt20(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "DONACIONES", "deducciondesgravacionsiradig", db)
	importeTope := getSubtotal(liquidacion, db) * 0.05 //5% de Subtotal
	importeTotal = getImporteTotalTope(importeTotal, importeTope)
	return importeTotal
}

func getGananciaNeta(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arrayGananciaNeta []float64
	var gananciaNeta float64

	arrayGananciaNeta = append(arrayGananciaNeta, getCuotaMedicoAsistencial(liquidacion, db))
	arrayGananciaNeta = append(arrayGananciaNeta, getDonacionFiscosNacProvMunArt20(liquidacion, db))

	gananciaNeta = getSubtotal(liquidacion, db) - Sum(arrayGananciaNeta)

	return gananciaNeta
}

func getDetalleCargoFamiliar(liquidacion *structLiquidacion.Liquidacion, columnaDetalleCargoFamiliar string, db *gorm.DB) float64 {
	var importeTotal, valorbeneficio float64
	anioperiodoliquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	mesperiodoliquidacion := getMes(&liquidacion.Fechaperiodoliquidacion)

	var detallecargofamiliar *structSiradig.Detallecargofamiliarsiradig
	sql := "SELECT dcfs.* FROM siradig s INNER JOIN detallecargofamiliarsiradig dcfs ON s.id = dcfs.siradigid where to_char(periodosiradig, 'YYYY') = ' " + strconv.Itoa(anioperiodoliquidacion) + "' AND dcfs." + columnaDetalleCargoFamiliar + " NOTNULL AND s.legajoid = " + strconv.Itoa(*liquidacion.Legajoid)
	db.Raw(sql).Scan(&detallecargofamiliar)

	/*sql = "SELECT valor FROM siradig s INNER JOIN beneficiosiradig bs ON s.id = bs.siradigid WHERE " + mesperiodoliquidacion + " BETWEEN bs.mesdesde AND bs.meshasta AND to_char(periodosiradig, 'YYYY') = ' " + strconv.Itoa(anioperiodoliquidacion) + "'"
	db.Raw(sql).Row().Scan(&valorbeneficio)
	*/
	if detallecargofamiliar.ID != 0 {

		mesdadobaja := getMes(detallecargofamiliar.Meshasta)
		mesdadoalta := getMes(detallecargofamiliar.Mesdesde)
		valorfijoconyuge := getValorFijoImpuestoGanancia(liquidacion, "deduccionespersonales", "valorfijoconyuge", db)

		/*if valorbeneficio != false {
			valorfijoconyuge = valorfijoconyuge * 1.22
		}*/

		if mesdadobaja == 0 {
			importeTotal = (valorfijoconyuge / 12) * float64(mesperiodoliquidacion-mesdadoalta)
		}

		if mesdadobaja <= mesperiodoliquidacion {
			importeTotal = (valorfijoconyuge / 12) * float64(mesdadobaja-mesdadoalta)
		}

		if mesdadobaja > mesperiodoliquidacion {
			importeTotal = (valorfijoconyuge / 12) * float64(mesperiodoliquidacion-mesdadoalta)
		}
	}

	return importeTotal
}

func getImporteTotalTope(importeTotal float64, tope float64) float64 {
	if importeTotal > tope {
		return tope
	} else {
		return importeTotal
	}
}

func getValorFijoImpuestoGanancia(liquidacion *structLiquidacion.Liquidacion, nombretabla string, nombrecolumna string, db *gorm.DB) float64 {
	var importeTope float64
	anioLiquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	sql := "SELECT" + nombrecolumna + "FROM " + nombretabla + " WHERE anio = " + strconv.Itoa(anioLiquidacion)
	db.Raw(sql).Row().Scan(&importeTope)

	return importeTope
}

func getMesesAProrratear(concepto *structConcepto.Concepto, liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) int {
	fechadesde := strconv.Itoa(liquidacion.Fechaperiodoliquidacion.Year()) + "-01-01"
	fechahasta := liquidacion.Fechaperiodoliquidacion.Format("2006-01-02")
	var fechaliquidacionmasantigua *time.Time
	var sql string
	mesAProrratear := getMes(&liquidacion.Fechaperiodoliquidacion)

	sql = "SELECT l.fechaperiodoliquidacion FROM liquidacion l INNER JOIN liquidacionitem li ON l.id = li.liquidacionid INNER JOIN  concepto c ON c.id = li.conceptoid INNER JOIN legajo le ON le.id = l.legajoid WHERE c.id = " + strconv.Itoa(concepto.ID) + " AND l.fechaperiodoliquidacion BETWEEN '" + fechadesde + "' AND '" + fechahasta + "' AND le.id = " + strconv.Itoa(*liquidacion.Legajoid) + " ORDER BY fechaperiodoliquidacion ASC LIMIT 1"
	fmt.Println("sql:", sql)
	db.Raw(sql).Row().Scan(&fechaliquidacionmasantigua)

	if fechaliquidacionmasantigua != nil {
		mesLiquidacionBD := getMes(fechaliquidacionmasantigua)

		if mesLiquidacionBD < mesAProrratear {
			mesAProrratear = mesLiquidacionBD
		}
	}

	return 13 - mesAProrratear

}

func getMes(fecha *time.Time) int {
	mes, _ := strconv.Atoi(s.Split(fecha.String(), "-")[1])
	return mes
}

func getImporteTotalSegunTipoImpuestoGanancias(tipoImpuestoALasGanancias string, liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var mes float64 = 1
	var importeTotal, importeConcepto float64

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		concepto := liquidacionitem.Concepto
		if concepto.Codigo == tipoImpuestoALasGanancias {
			if concepto.Prorrateo == true {
				mes = float64(getMesesAProrratear(concepto, liquidacion, db))
			}
			importeConcepto = *liquidacionitem.Importeunitario / mes
			importeTotal = importeTotal + importeConcepto
		}
	}
	return importeTotal
}

func getSacCuotas(liquidacion *structLiquidacion.Liquidacion, correspondeSemestre bool, db *gorm.DB) float64 {
	var mes float64 = 1
	var importeTotal, importeConcepto float64

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		concepto := liquidacionitem.Concepto
		if correspondeSemestre && concepto.Basesac == true {
			if concepto.Prorrateo == true {
				mes = float64(getMesesAProrratear(concepto, liquidacion, db))
			}
			importeConcepto = *liquidacionitem.Importeunitario / mes
			importeTotal = importeTotal + importeConcepto
		}
	}
	return importeTotal / 12
}

func getImporteGananciasOtroEmpleoSiradig(liquidacion *structLiquidacion.Liquidacion, columnaimportegananciasotroempleosiradig string, db *gorm.DB) float64 {
	var importeTotal float64
	fechaliquidacion := liquidacion.Fechaperiodoliquidacion.Format("2006-01-02")

	sql := "SELECT SUM(" + columnaimportegananciasotroempleosiradig + ") FROM importegananciasotroempleosiradig WHERE '" + fechaliquidacion + "' >= mes"
	fmt.Println("remuneracion bruta:", sql)
	db.Raw(sql).Row().Scan(&importeTotal)

	return importeTotal
}

func Sum(s []float64) float64 {
	var sum float64
	for _, val := range s {
		sum += val
	}
	return sum
}
