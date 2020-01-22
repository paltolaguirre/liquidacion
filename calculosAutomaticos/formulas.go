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

func getfgMesesAProrratear(concepto *structConcepto.Concepto, liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) int {
	fechadesde := strconv.Itoa(liquidacion.Fechaperiodoliquidacion.Year()) + "-01-01"
	fechahasta := liquidacion.Fechaperiodoliquidacion.Format("2006-01-02")
	var fechaliquidacionmasantigua *time.Time
	var sql string
	mesAProrratear := getfgMes(&liquidacion.Fechaperiodoliquidacion)

	sql = "SELECT l.fechaperiodoliquidacion FROM liquidacion l INNER JOIN liquidacionitem li ON l.id = li.liquidacionid INNER JOIN  concepto c ON c.id = li.conceptoid INNER JOIN legajo le ON le.id = l.legajoid WHERE c.id = " + strconv.Itoa(concepto.ID) + " AND l.fechaperiodoliquidacion BETWEEN '" + fechadesde + "' AND '" + fechahasta + "' AND le.id = " + strconv.Itoa(*liquidacion.Legajoid) + " ORDER BY fechaperiodoliquidacion ASC LIMIT 1"
	db.Raw(sql).Row().Scan(&fechaliquidacionmasantigua)

	if fechaliquidacionmasantigua != nil {
		mesLiquidacionBD := getfgMes(fechaliquidacionmasantigua)

		if mesLiquidacionBD < mesAProrratear {
			mesAProrratear = mesLiquidacionBD
		}
	}

	return 13 - mesAProrratear

}

func getfgImporteTotalSegunTipoImpuestoGanancias(tipoImpuestoALasGanancias string, liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var mes float64 = 1
	var importeTotal, importeConcepto float64

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		concepto := liquidacionitem.Concepto
		tipoimpuesto := concepto.Tipoimpuestoganancias.Codigo
		if tipoimpuesto == tipoImpuestoALasGanancias {
			if concepto.Prorrateo == true {
				mes = float64(getfgMesesAProrratear(concepto, liquidacion, db))
			}
			importeConcepto = *liquidacionitem.Importeunitario / mes
			importeTotal = importeTotal + importeConcepto
		}
	}
	return importeTotal
}

func getfgRemuneracionBruta(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("REMUNERACION_BRUTA", liquidacion, db)
	return importeTotal
}

func getfgRemuneracionNoHabitual(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("RETRIBUCIONES_NO_HABITUALES", liquidacion, db)
	return importeTotal
}

func getfgSacCuotas(liquidacion *structLiquidacion.Liquidacion, correspondeSemestre bool, db *gorm.DB) float64 {
	var mes float64 = 1
	var importeTotal, importeConcepto float64

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		concepto := liquidacionitem.Concepto
		if correspondeSemestre && concepto.Basesac == true {
			if concepto.Prorrateo == true {
				mes = float64(getfgMesesAProrratear(concepto, liquidacion, db))
			}
			importeConcepto = *liquidacionitem.Importeunitario / mes
			if *concepto.Tipoconceptoid == -4 {
				importeConcepto = importeConcepto * -1
			}
			importeTotal = importeTotal + importeConcepto
		}
	}
	return importeTotal / 12
}

func getfgSacPrimerCuota(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	correspondePrimerSemetre := getfgMes(&liquidacion.Fechaperiodoliquidacion) <= 6
	importeTotal := getfgSacCuotas(liquidacion, correspondePrimerSemetre, db)

	return importeTotal

}

func getfgSacSegundaCuota(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	correspondeSegundoSemetre := getfgMes(&liquidacion.Fechaperiodoliquidacion) > 6
	importeTotal := getfgSacCuotas(liquidacion, correspondeSegundoSemetre, db)

	return importeTotal
}

func getfgHorasExtrasGravadas(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("HORAS_EXTRAS_REMUNERACION_GRAVADA", liquidacion, db)
	return importeTotal
}

func getfgMovilidadYViaticosGravada(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("MOVILIDAD_Y_VIATICOS_REMUNERACION_GRAVADA", liquidacion, db)
	return importeTotal
}

func getfgMaterialDidacticoPersonalDocenteRemuneracion(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("MATERIAL_DIDACTICO_PERSONAL_DOCENTE_REMUNERACION_GRAVADA", liquidacion, db)
	return importeTotal
}

func getfgImporteGananciasOtroEmpleoSiradig(liquidacion *structLiquidacion.Liquidacion, columnaimportegananciasotroempleosiradig string, db *gorm.DB) float64 {
	var importeTotal float64
	fechaliquidacion := liquidacion.Fechaperiodoliquidacion.Format("2006-01-02")

	sql := "SELECT SUM(" + columnaimportegananciasotroempleosiradig + ") FROM importegananciasotroempleosiradig WHERE '" + fechaliquidacion + "' >= mes"
	db.Raw(sql).Row().Scan(&importeTotal)

	return importeTotal
}

func getfgRemuneracionBrutaOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "importegananciasbrutas", db)
	return importeTotal
}

func getfgRemuneracionNoHabitualOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "importeretribucionesnohabituales", db)
	return importeTotal
}

func getfgSacPrimerCuotaOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0
	if getfgMes(&liquidacion.Fechaperiodoliquidacion) <= 6 {
		importeTotal = getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "sac", db)
	}
	return importeTotal
}

func getfgSacSegundaCuotaOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0
	if getfgMes(&liquidacion.Fechaperiodoliquidacion) > 6 {
		importeTotal = getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "sac", db)
	}
	return importeTotal
}

func getfgHorasExtrasGravadasOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "importehorasextrasgravadas", db)
	return importeTotal
}

func getfgMovilidadYViaticosGravadaOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "gastosmovilidad", db)
	return importeTotal
}

func getfgMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "materialdidactico", db)
	return importeTotal
}

func getfgSubtotalIngresos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {

	var arraySubtotalIngresos []float64
	var subtotalIngresos float64

	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgRemuneracionBruta(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgRemuneracionNoHabitual(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgSacPrimerCuota(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgSacSegundaCuota(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgHorasExtrasGravadas(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgMovilidadYViaticosGravada(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgMaterialDidacticoPersonalDocenteRemuneracion(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgRemuneracionBrutaOtrosEmpleos(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgRemuneracionNoHabitualOtrosEmpleos(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgSacPrimerCuotaOtrosEmpleos(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgSacSegundaCuotaOtrosEmpleos(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgHorasExtrasGravadasOtrosEmpleos(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgMovilidadYViaticosGravadaOtrosEmpleos(liquidacion, db))
	arraySubtotalIngresos = append(arraySubtotalIngresos, getfgMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos(liquidacion, db))

	subtotalIngresos = Sum(arraySubtotalIngresos)
	return subtotalIngresos
}

func getfgAportesJubilatoriosRetirosPensionesOSubsidios(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("APORTES_JUBILATORIOS_RETIROS_PENSIONES_O_SUBSIDIOS", liquidacion, db)
	return importeTotal
}

func getfgAportesObraSocial(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("APORTES_OBRA_SOCIAL", liquidacion, db)
	return importeTotal
}

func getfgCuotaSindical(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("CUOTA_SINDICAL", liquidacion, db)
	return importeTotal
}

func getfgDescuentosObligatoriosPorLeyNacionalProvincialOMunicipal(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("DESCUENTOS_OBLIGATORIOS_POR_LEY_NACIONAL_PROVINCIAL_MUNICIPAL", liquidacion, db)
	return importeTotal
}

func getfgGastosMovilidadViaticosAbonadosPorElEmpleador(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arrayGastosMovilidad []float64
	var importeTotal float64

	arrayGastosMovilidad = append(arrayGastosMovilidad, getfgMovilidadYViaticosGravada(liquidacion, db))
	arrayGastosMovilidad = append(arrayGastosMovilidad, getfgMovilidadYViaticosGravadaOtrosEmpleos(liquidacion, db))
	arrayGastosMovilidad = append(arrayGastosMovilidad, getfgMaterialDidacticoPersonalDocenteRemuneracion(liquidacion, db))
	arrayGastosMovilidad = append(arrayGastosMovilidad, getfgMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos(liquidacion, db))

	importeTotal = Sum(arrayGastosMovilidad)
	return importeTotal

}

func getfgAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "aporteseguridadsocial", db)
	return importeTotal
}

func getfgAportesObraSocialOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "aporteobrasocial", db)
	return importeTotal
}

func getfgCuotaSindicalOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "aportesindical", db)
	return importeTotal
}

func getfgPrimasDeSeguroParaCasoDeMuerte(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	return 0
}

func getfgSeguroMuerteMixtosSujetosAlControlSSN(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	return 0
}

func getfgSegurosRetirosPrivadosSujetosAlControlSSN(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	return 0
}

func getfgImporteTotalSiradigSegunTipoGrilla(liquidacion *structLiquidacion.Liquidacion, columnadeducciondesgravacionsiradig string, tipodeducciondesgravacionsiradig string, nombretablasiradig string, db *gorm.DB) float64 {
	var importeTotal float64
	fechaliquidacion := liquidacion.Fechaperiodoliquidacion.Format("2006-01-02")

	sql := "SELECT SUM(" + columnadeducciondesgravacionsiradig + ") FROM " + nombretablasiradig + " ts INNER JOIN siradigtipogrilla stg ON stg.id = ts.siradigtipogrillaid WHERE '" + fechaliquidacion + "' >= mes AND stg.codigo = '" + tipodeducciondesgravacionsiradig + "'"

	db.Raw(sql).Row().Scan(&importeTotal)

	return importeTotal
}

func getfgImporteTotalTope(importeTotal float64, tope float64) float64 {
	if importeTotal > tope {
		return tope
	} else {
		return importeTotal
	}
}

func getfgValorFijoImpuestoGanancia(liquidacion *structLiquidacion.Liquidacion, nombretabla string, nombrecolumna string, db *gorm.DB) float64 {
	var importeTope float64
	anioLiquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	sql := "SELECT " + nombrecolumna + " FROM " + nombretabla + " WHERE anio = " + strconv.Itoa(anioLiquidacion)
	fmt.Println("valor fijo:", sql)
	db.Raw(sql).Row().Scan(&importeTope)

	return importeTope
}

func getfgGastosSepelio(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "GASTOS_DE_SEPELIO", "deducciondesgravacionsiradig", db)
	importeTope := getfgValorFijoImpuestoGanancia(liquidacion, "topemaximodescuento", "topesepelio", db)
	importeTotal = getfgImporteTotalTope(importeTotal, importeTope)
	return importeTotal
}

func getfgGastosAmortizacionEInteresesRodadoCorredoresViajantesComercio(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "GASTOS_DE_REPRESENTACION_E_INTERESES_DE_CORREDORES_Y_VIAJANTES_DE_COMERCIO", "deducciondesgravacionsiradig", db)
	return importeTotal
}

func getfgInteresesCreditosHipotecarios(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "INTERESES_PRESTAMO_HIPOTECARIO", "deducciondesgravacionsiradig", db)
	importeTope := getfgValorFijoImpuestoGanancia(liquidacion, "topemaximodescuento", "topehipotecarios", db)
	importeTotal = getfgImporteTotalTope(importeTotal, importeTope)
	return importeTotal
}

func getfgAportesCapSocFondoRiesgoSociosProtectoresSGR(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "Montoreintegrar + Montoreintegrar3", "REINTEGRO_DE_APORTES_DE_SOCIOS_PROTECTORES_A_SOCIEDADES_DE_GARANTIA_RECIPROCA", "ajustesiradig", db)
	return importeTotal
}

func getfgAlquileresInmueblesDestinadosASuCasaHabitacion(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "ALQUILER_INMUEBLES_DESTINADOS_A_CASA_HABITACION", "deducciondesgravacionsiradig", db)
	importeTope := getfgMinimoNoImponible(liquidacion, db) * 0.4 /*es el 40% de MNI(40)*/
	importeTotal = getfgImporteTotalTope(importeTotal, importeTope)
	return importeTotal
}

func getfgEmpleadosServicioDomestico(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "contribucion + retribucion", "DEDUCCION_DEL_PERSONAL_DOMESTICO", "deducciondesgravacionsiradig", db)
	importeTope := getfgMinimoNoImponible(liquidacion, db) /*es el MNI(40)*/
	importeTotal = getfgImporteTotalTope(importeTotal, importeTope)
	return importeTotal
}

func getfgIndumentariaEquipamientoCaracterObligatorio(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "GASTOS_ADQUISICION_INDUMENTARIA_Y_EQUIPAMIENTO_PARA_USO_EXCLUSIVO_EN_EL_LUGAR_DE_TRABAJO", "deducciondesgravacionsiradig", db)
	return importeTotal
}

func getfgOtrasDeducciones(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "OTRAS", "deducciondesgravacionsiradig", db)
	return importeTotal
}

func getfgSubtotal(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arraySubtotal []float64
	var subTotal float64

	arraySubtotal = append(arraySubtotal, getfgAportesJubilatoriosRetirosPensionesOSubsidios(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgAportesObraSocial(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgCuotaSindical(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgDescuentosObligatoriosPorLeyNacionalProvincialOMunicipal(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgGastosMovilidadViaticosAbonadosPorElEmpleador(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgAportesObraSocialOtrosEmpleos(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgCuotaSindicalOtrosEmpleos(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgPrimasDeSeguroParaCasoDeMuerte(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgSeguroMuerteMixtosSujetosAlControlSSN(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgSegurosRetirosPrivadosSujetosAlControlSSN(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgGastosSepelio(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgGastosAmortizacionEInteresesRodadoCorredoresViajantesComercio(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgInteresesCreditosHipotecarios(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgAportesCapSocFondoRiesgoSociosProtectoresSGR(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgAlquileresInmueblesDestinadosASuCasaHabitacion(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgEmpleadosServicioDomestico(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgIndumentariaEquipamientoCaracterObligatorio(liquidacion, db))
	arraySubtotal = append(arraySubtotal, getfgOtrasDeducciones(liquidacion, db))

	subTotal = getfgSubtotalIngresos(liquidacion, db) - Sum(arraySubtotal)

	return subTotal
}

func getfgCuotaMedicoAsistencial(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "CUOTA_MEDICA_ASISTENCIAL", "deducciondesgravacionsiradig", db)
	importeTope := getfgSubtotal(liquidacion, db) * 0.05 //5% de Subtotal
	importeTotal = getfgImporteTotalTope(importeTotal, importeTope)
	return importeTotal
}

func getfgDonacionFiscosNacProvMunArt20(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "DONACIONES", "deducciondesgravacionsiradig", db)
	importeTope := getfgSubtotal(liquidacion, db) * 0.05 //5% de Subtotal
	importeTotal = getfgImporteTotalTope(importeTotal, importeTope)
	return importeTotal
}

func getfgGananciaNeta(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arrayGananciaNeta []float64
	var gananciaNeta float64

	arrayGananciaNeta = append(arrayGananciaNeta, getfgCuotaMedicoAsistencial(liquidacion, db))
	arrayGananciaNeta = append(arrayGananciaNeta, getfgDonacionFiscosNacProvMunArt20(liquidacion, db))

	gananciaNeta = getfgSubtotal(liquidacion, db) - Sum(arrayGananciaNeta)

	return gananciaNeta
}

func getfgDetalleCargoFamiliar(liquidacion *structLiquidacion.Liquidacion, columnaDetalleCargoFamiliar string, valorfijocolumna string, porcentaje float64, db *gorm.DB) float64 {
	var importeTotal float64
	var tienevalorbeneficio bool
	anioperiodoliquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	mesperiodoliquidacion := getfgMes(&liquidacion.Fechaperiodoliquidacion)

	var detallecargofamiliar structSiradig.Detallecargofamiliarsiradig
	sql := "SELECT dcfs.* FROM siradig s INNER JOIN detallecargofamiliarsiradig dcfs ON s.id = dcfs.siradigid where to_char(periodosiradig, 'YYYY') = '" + strconv.Itoa(anioperiodoliquidacion) + "' AND dcfs." + columnaDetalleCargoFamiliar + " NOTNULL AND s.legajoid = " + strconv.Itoa(*liquidacion.Legajoid)
	db.Raw(sql).Scan(&detallecargofamiliar)

	sql = "SELECT valor FROM siradig s INNER JOIN beneficiosiradig bs ON s.id = bs.siradigid WHERE to_number(to_char(bs.mesdesde, 'MM'),'99') <= " + strconv.Itoa(mesperiodoliquidacion) + " AND to_number(to_char(bs.meshasta, 'MM'), '99') > " + strconv.Itoa(mesperiodoliquidacion) + " AND bs.siradigtipogrillaid = -24 AND to_char(s.periodosiradig, 'YYYY') = '" + strconv.Itoa(anioperiodoliquidacion) + "'"
	db.Raw(sql).Row().Scan(&tienevalorbeneficio)

	if detallecargofamiliar.ID != 0 {

		mesdadobaja := getfgMes(detallecargofamiliar.Meshasta)
		mesdadoalta := getfgMes(detallecargofamiliar.Mesdesde)
		valorfijo := getfgValorFijoImpuestoGanancia(liquidacion, "deduccionespersonales", valorfijocolumna, db)

		if tienevalorbeneficio == true {
			valorfijo = valorfijo * 1.22
		}

		if mesdadobaja == 0 {
			importeTotal = (valorfijo / 12) * float64(mesperiodoliquidacion-mesdadoalta) * porcentaje
		} else {
			if mesdadobaja <= mesperiodoliquidacion {
				importeTotal = (valorfijo / 12) * float64(mesdadobaja-mesdadoalta) * porcentaje
			} else {
				if mesdadobaja > mesperiodoliquidacion {
					importeTotal = (valorfijo / 12) * float64(mesperiodoliquidacion-mesdadoalta) * porcentaje
				}
			}
		}
	}

	return importeTotal
}

func getfgConyuge(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgDetalleCargoFamiliar(liquidacion, "conyugeid", "valorfijoconyuge", 1, db)
	return importeTotal
}

func getfgHijos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0
	var detallescargofamiliarsiradig []structSiradig.Detallecargofamiliarsiradig

	valorfijoMNI := getfgValorFijoImpuestoGanancia(liquidacion, "deduccionespersonales", "valorfijomni", db)
	sql := "SELECT * FROM detallecargofamiliar WHERE hijoid NOT NULL AND estaacargo = true AND montoanual < " + strconv.FormatFloat(valorfijoMNI, 'f', 5, 64)
	fmt.Println("hijos", sql)
	db.Raw(sql).Scan(&detallescargofamiliarsiradig)

	for i := 0; i < len(detallescargofamiliarsiradig); i++ {
		porcentaje := detallescargofamiliarsiradig[i].Porcentaje
		importeTotal = importeTotal + getfgDetalleCargoFamiliar(liquidacion, "hijoid", "valorfijohijo", *porcentaje, db)
	}

	return importeTotal
}

func getfgMinimoNoImponible(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	valorfijoMNI := getfgValorFijoImpuestoGanancia(liquidacion, "deduccionespersonales", "valorfijomni", db)
	mesperiodoliquidacion := getfgMes(&liquidacion.Fechaperiodoliquidacion)
	importeTotal := (valorfijoMNI / 12) * float64(mesperiodoliquidacion)
	return importeTotal
}

func getfgDeduccionEspecial(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	valorfijoMNI := getfgValorFijoImpuestoGanancia(liquidacion, "deduccionespersonales", "valorfijoddei", db)
	mesperiodoliquidacion := getfgMes(&liquidacion.Fechaperiodoliquidacion)
	importeTotal := (valorfijoMNI / 12) * float64(mesperiodoliquidacion)
	return importeTotal
}

func getfgSubtotalDeduccionesPersonales(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arraySubtotalDeduccionesPersonales []float64
	var subTotalDeduccionesPersonales float64

	arraySubtotalDeduccionesPersonales = append(arraySubtotalDeduccionesPersonales, getfgConyuge(liquidacion, db))
	arraySubtotalDeduccionesPersonales = append(arraySubtotalDeduccionesPersonales, getfgHijos(liquidacion, db))
	arraySubtotalDeduccionesPersonales = append(arraySubtotalDeduccionesPersonales, getfgMinimoNoImponible(liquidacion, db))
	arraySubtotalDeduccionesPersonales = append(arraySubtotalDeduccionesPersonales, getfgDeduccionEspecial(liquidacion, db))

	subTotalDeduccionesPersonales = Sum(arraySubtotalDeduccionesPersonales)
	return subTotalDeduccionesPersonales
}

func getfgDeduccionesAComputar(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgSubtotalDeduccionesPersonales(liquidacion, db)
	return importeTotal
}

func getfgGananciaNetaAcumSujetaAImp(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var liquidaciones []structLiquidacion.Liquidacion
	var importeTotal float64 = 0
	anioperiodoliquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	sql := "SELECT l.* FROM liquidacion l INNER JOIN legajo le ON le.id = l.legajoid WHERE to_char(l.fechaperiodoliquidacion, 'YYYY') = ' " + strconv.Itoa(anioperiodoliquidacion) + "' AND le.ID = " + strconv.Itoa(*liquidacion.Legajoid)
	db.Raw(sql).Scan(&liquidaciones)

	for i := 0; i < len(liquidaciones); i++ {
		importeTotal = importeTotal + getfgGananciaNeta(&liquidaciones[i], db)
	}

	return importeTotal
}

func GetfgDeduccionesPersonales(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgDeduccionesAComputar(liquidacion, db)
	return importeTotal * -1
}

func getfgBaseImponible(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arrayBaseImponible []float64
	var importeTotal float64

	arrayBaseImponible = append(arrayBaseImponible, getfgGananciaNetaAcumSujetaAImp(liquidacion, db))
	arrayBaseImponible = append(arrayBaseImponible, GetfgDeduccionesPersonales(liquidacion, db))

	importeTotal = Sum(arrayBaseImponible)
	return importeTotal
}

func getfgTotalGananciaNetaImponibleAcumuladaSinHorasExtras(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {

	cuotaSindical := getfgCuotaSindical(liquidacion, db)
	obraSocial := getfgAportesObraSocial(liquidacion, db)
	aportesJubilatorios := getfgAportesJubilatoriosRetirosPensionesOSubsidios(liquidacion, db)
	remunerativosMenosDescuentos := obtenerRemunerativosMenosDescuentos(liquidacion)
	cuotaSindicalOtros := getfgCuotaSindicalOtrosEmpleos(liquidacion, db)
	obraSocialOtros := getfgAportesObraSocialOtrosEmpleos(liquidacion, db)
	aportesJubilatoriosOtros := getfgAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos(liquidacion, db)
	remunerativosOtros := obtenerRemunerativosOtros(liquidacion, db)

	var uno float64 = 1
	porcentaje := uno - (cuotaSindical/remunerativosMenosDescuentos + obraSocial/remunerativosMenosDescuentos + aportesJubilatorios/remunerativosMenosDescuentos)
	porcentajeOtrosEmp := uno - (cuotaSindicalOtros/remunerativosOtros + obraSocialOtros/remunerativosOtros + aportesJubilatoriosOtros/remunerativosOtros)

	importeTotal := getfgBaseImponible(liquidacion, db) - (getfgHorasExtrasGravadas(liquidacion, db)*porcentaje + getfgHorasExtrasGravadasOtrosEmpleos(liquidacion, db)*porcentajeOtrosEmp)
	return importeTotal
}

func obtenerRemunerativosMenosDescuentos(liquidacion *structLiquidacion.Liquidacion) float64 {
	var totalRemunerativos, totalDescuentos float64
	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		tipoconcepto := *liquidacionitem.Concepto.Tipoconceptoid
		importeconcepto := *liquidacionitem.Importeunitario
		if tipoconcepto == -1 {
			totalRemunerativos = totalRemunerativos + importeconcepto
		}
		if tipoconcepto == -3 {
			totalDescuentos = totalDescuentos + importeconcepto
		}

	}
	return totalRemunerativos - totalDescuentos
}

func obtenerRemunerativosOtros(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arrayRemunerativosOtros []float64
	var totalRemunerativosOtros float64

	arrayRemunerativosOtros = append(arrayRemunerativosOtros, getfgRemuneracionBrutaOtrosEmpleos(liquidacion, db))
	arrayRemunerativosOtros = append(arrayRemunerativosOtros, getfgSacPrimerCuotaOtrosEmpleos(liquidacion, db))
	arrayRemunerativosOtros = append(arrayRemunerativosOtros, getfgSacSegundaCuotaOtrosEmpleos(liquidacion, db))
	arrayRemunerativosOtros = append(arrayRemunerativosOtros, getfgHorasExtrasGravadasOtrosEmpleos(liquidacion, db))

	totalRemunerativosOtros = Sum(arrayRemunerativosOtros)
	return totalRemunerativosOtros
}

type strEscalaimpuestoaplicable struct {
	Limiteinferior float64 `json:"limiteinferior"`
	Limitesuperior float64 `json:"limitesuperior"`
	Valorfijo      float64 `json:"valorfijo"`
	Valorvariable  float64 `json:"valorvariable"`
	Mesanio        string  `json:"mesanio"`
}

func getfgEscalaImpuestoAplicable(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) *[]strEscalaimpuestoaplicable {
	var strescalaimpuestoaplicable []strEscalaimpuestoaplicable

	anioLiquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	mesLiquidacion := getfgMes(&liquidacion.Fechaperiodoliquidacion)

	mesAnioLiquidacion := strconv.Itoa(mesLiquidacion) + "/" + strconv.Itoa(anioLiquidacion)

	sql := "SELECT limiteinferior,limitesuperior,valorfijo,valorvariable,mesanio FROM escalaimpuestoaplicable where mesanio = '" + mesAnioLiquidacion + "'"
	db.Raw(sql).Scan(&strescalaimpuestoaplicable)

	return &strescalaimpuestoaplicable

}

func getfgDeterminacionImpuestoFijo(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0
	strescalaimpuestoaplicable := *getfgEscalaImpuestoAplicable(liquidacion, db)
	totalganancianeta := getfgTotalGananciaNetaImponibleAcumuladaSinHorasExtras(liquidacion, db)

	for i := 0; i < len(strescalaimpuestoaplicable); i++ {
		escalaimpuestoaplicable := strescalaimpuestoaplicable[i]
		if totalganancianeta > escalaimpuestoaplicable.Limiteinferior && totalganancianeta <= escalaimpuestoaplicable.Limitesuperior {
			importeTotal = escalaimpuestoaplicable.Valorfijo
		}
	}
	return importeTotal
}

func getfgDeterminacionImpuestoPorEscala(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0
	strescalaimpuestoaplicable := *getfgEscalaImpuestoAplicable(liquidacion, db)
	totalganancianeta := getfgTotalGananciaNetaImponibleAcumuladaSinHorasExtras(liquidacion, db)
	baseimponible := getfgBaseImponible(liquidacion, db)
	for i := 0; i < len(strescalaimpuestoaplicable); i++ {
		escalaimpuestoaplicable := strescalaimpuestoaplicable[i]
		if totalganancianeta > escalaimpuestoaplicable.Limiteinferior && totalganancianeta <= escalaimpuestoaplicable.Limitesuperior {

			importeTotal = (baseimponible - escalaimpuestoaplicable.Limiteinferior) * escalaimpuestoaplicable.Valorvariable
		}
	}
	return importeTotal
}

func getfgTotalRetener(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arrayTotalRetener []float64
	var totalRetener float64

	arrayTotalRetener = append(arrayTotalRetener, getfgDeterminacionImpuestoFijo(liquidacion, db))
	arrayTotalRetener = append(arrayTotalRetener, getfgDeterminacionImpuestoPorEscala(liquidacion, db))

	totalRetener = Sum(arrayTotalRetener)
	return totalRetener
}

func getfgRetencionAcumulada(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	anioperiodoliquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	var totalconceptosimpuestoganancias float64
	sql := "SELECT SUM(li.importeunitario) FROM siradig s INNER JOIN legajo le ON le.id = s.legajoid INNER JOIN liquidacion l ON le.id = l.legajoid INNER JOIN liquidacionitem li ON l.id = li.liquidacionid INNER JOIN  concepto c ON c.id = li.conceptoid WHERE c.codigo = 'IMPUESTO_GANANCIAS' AND to_char(s.periodosiradig, 'YYYY') = ' " + strconv.Itoa(anioperiodoliquidacion)
	db.Raw(sql).Row().Scan(&totalconceptosimpuestoganancias)

	return totalconceptosimpuestoganancias
}

func GetfgRetencionMes(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	totalRetener := getfgTotalRetener(liquidacion, db)
	retencionAcumulada := getfgRetencionAcumulada(liquidacion, db)

	return totalRetener - retencionAcumulada
}

func getfgMes(fecha *time.Time) int {
	var mes int
	if fecha != nil {
		mes, _ = strconv.Atoi(s.Split(fecha.String(), "-")[1])
	}
	return mes
}

func Sum(s []float64) float64 {
	var sum float64
	for _, val := range s {
		sum += val
	}
	return sum
}
