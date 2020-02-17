package calculosAutomaticos

import (
	"fmt"
	"strconv"
	s "strings"
	"time"

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

	sql = "SELECT l.fechaperiodoliquidacion FROM liquidacion l INNER JOIN liquidacionitem li ON l.id = li.liquidacionid INNER JOIN  concepto c ON c.id = li.conceptoid INNER JOIN legajo le ON le.id = l.legajoid WHERE c.id = " + strconv.Itoa(concepto.ID) + " AND l.fechaperiodoliquidacion BETWEEN '" + fechadesde + "' AND '" + fechahasta + "' AND le.id = " + strconv.Itoa(*liquidacion.Legajoid) + "AND le.deleted_at IS NULL AND l.deleted_at IS NULL and li.deleted_at IS NULL, c.deleted_at IS NULL ORDER BY fechaperiodoliquidacion ASC LIMIT 1"
	db.Raw(sql).Row().Scan(&fechaliquidacionmasantigua)

	if fechaliquidacionmasantigua != nil {
		mesLiquidacionBD := getfgMes(fechaliquidacionmasantigua)

		if mesLiquidacionBD < mesAProrratear {
			mesAProrratear = mesLiquidacionBD
		}
	}
	fmt.Println("Calculos Automaticos - Mes a Prorratear:", 13-mesAProrratear)
	return 13 - mesAProrratear

}

func getfgImporteTotalSegunTipoImpuestoGanancias(tipoImpuestoALasGanancias string, liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal, importeConcepto float64

	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		concepto := liquidacionitem.Concepto
		tipoimpuesto := obtenerTipoImpuesto(concepto, db)
		var mes float64 = 1

		if tipoimpuesto == tipoImpuestoALasGanancias && concepto.Codigo != "IMPUESTO_GANANCIAS" && concepto.Codigo != "IMPUESTO_GANANCIAS_DEVOLUCION" {
			if concepto.Prorrateo == true {
				mes = float64(getfgMesesAProrratear(concepto, liquidacion, db))
			}
			importeLiquidacionitem := liquidacionitem.Importeunitario
			if importeLiquidacionitem != nil {
				importeConcepto = *importeLiquidacionitem / mes
			}
			importeTotal = importeTotal + importeConcepto

			importeTotal = importeTotal + obtenerConceptosProrrateoMesesAnteriores(liquidacion, db)
		}
	}
	return importeTotal
}

type importeMes struct {
	Importe        float64
	Mesliquidacion string
}

func obtenerConceptosProrrateoMesesAnteriores(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importemes []importeMes
	anioLiquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	legajoID := liquidacion.Legajoid

	sql := "SELECT li.importeunitario, to_char(l.fechaperiodoliquidacion, 'MM') AS mesliquidacion FROM liquidacion l INNER JOIN liquidacionitem li on l.id = li.liquidacionid INNER JOIN legajo le on le.id = l.legajoid INNER JOIN concepto c on c.id = li.conceptoid WHERE li.ID != " + strconv.Itoa(liquidacion.ID) + "AND to_char(l.fechaperiodoliquidacion, 'YYYY') = '" + strconv.Itoa(anioLiquidacion) + "' and le.id = " + strconv.Itoa(*legajoID) + " and c.prorrateo = true ORDER BY to_char(l.fechaperiodoliquidacion, 'MM') ASC"
	db.Raw(sql).Row().Scan(&importemes)
	var mes float64 = 1
	var trece float64 = 13
	var importeTotal float64 = 0
	for i := 0; i < len(importemes); i++ {
		mesLiquidacion, _ := strconv.ParseFloat(importemes[i].Mesliquidacion, 64)
		if mesLiquidacion < mes {
			mes = mesLiquidacion
		}
		importeConcepto := importemes[i].Importe / (trece - mes)

		importeTotal = importeTotal + importeConcepto
	}

	return importeTotal

}

func obtenerTipoImpuesto(concepto *structConcepto.Concepto, db *gorm.DB) string {
	var tipoimpuesto string
	if concepto.Tipoimpuestoganancias != nil {
		tipoimpuesto = concepto.Tipoimpuestoganancias.Codigo
	} else {
		sql := "SELECT codigo FROM tipoimpuestoganancias WHERE id = " + strconv.Itoa(*concepto.Tipoimpuestogananciasid)
		db.Raw(sql).Row().Scan(&tipoimpuesto)

	}

	return tipoimpuesto
}
func getfgRemuneracionBruta(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("REMUNERACION_BRUTA", liquidacion, db)
	fmt.Println("Calculos Automaticos - Remuneracion Bruta:", importeTotal)
	return importeTotal
}

func getfgRemuneracionNoHabitual(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("RETRIBUCIONES_NO_HABITUALES", liquidacion, db)
	fmt.Println("Calculos Automaticos - Remuneracion No Habitual:", importeTotal)
	return importeTotal
}

func getfgSacCuotas(liquidacion *structLiquidacion.Liquidacion, correspondeSemestre bool, db *gorm.DB) float64 {
	var importeTotal, importeConcepto float64

	if correspondeSemestre {
		for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
			liquidacionitem := liquidacion.Liquidacionitems[i]
			concepto := liquidacionitem.Concepto
			var mes float64 = 1
			if concepto.Basesac == true {
				if concepto.Prorrateo == true {
					mes = float64(getfgMesesAProrratear(concepto, liquidacion, db))
				}
				importeLiquidacionitem := liquidacionitem.Importeunitario
				if importeLiquidacionitem != nil {
					importeConcepto = *importeLiquidacionitem / mes
				}

				if *concepto.Tipoconceptoid == -4 {
					importeConcepto = importeConcepto * -1
				}
				importeTotal = importeTotal + importeConcepto
			}
		}

		importeTotal = importeTotal + getfgBaseSacOtrosEmpleos(liquidacion, db)

		importeTotal = importeTotal + obtenerConceptosProrrateoMesesAnteriores(liquidacion, db)
	}

	return importeTotal / 12
}

func getfgBaseSacOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {

	return getfgRemuneracionBrutaOtrosEmpleos(liquidacion, db) + getfgRemuneracionNoHabitualOtrosEmpleos(liquidacion, db) + getfgHorasExtrasGravadasOtrosEmpleos(liquidacion, db) + getfgMovilidadYViaticosGravadaOtrosEmpleos(liquidacion, db) + getfgMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos(liquidacion, db) - getfgAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos(liquidacion, db) - getfgAportesObraSocialOtrosEmpleos(liquidacion, db) - getfgCuotaSindicalOtrosEmpleos(liquidacion, db)
}

func GetfgSacPrimerCuota(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	correspondePrimerSemetre := getfgMes(&liquidacion.Fechaperiodoliquidacion) <= 6
	importeTotal := getfgSacCuotas(liquidacion, correspondePrimerSemetre, db)
	fmt.Println("Calculos Automaticos - Sac Primer Cuota:", importeTotal)
	return importeTotal

}

func getfgSacSegundaCuota(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	correspondeSegundoSemetre := getfgMes(&liquidacion.Fechaperiodoliquidacion) > 6
	importeTotal := getfgSacCuotas(liquidacion, correspondeSegundoSemetre, db)
	fmt.Println("Calculos Automaticos - Sac Segunda Cuota:", importeTotal)
	return importeTotal
}

func getfgHorasExtrasGravadas(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("HORAS_EXTRAS_REMUNERACION_GRAVADA", liquidacion, db)
	fmt.Println("Calculos Automaticos - Horas Extras Gravadas:", importeTotal)
	return importeTotal
}

func getfgMovilidadYViaticosGravada(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("MOVILIDAD_Y_VIATICOS_REMUNERACION_GRAVADA", liquidacion, db)
	fmt.Println("Calculos Automaticos - Movilidad y Viaticos Gravada:", importeTotal)
	return importeTotal
}

func getfgMaterialDidacticoPersonalDocenteRemuneracion(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("MATERIAL_DIDACTICO_PERSONAL_DOCENTE_REMUNERACION_GRAVADA", liquidacion, db)
	fmt.Println("Calculos Automaticos - Material Didactico Personal Docente Remuneracion:", importeTotal)
	return importeTotal
}

func getfgImporteGananciasOtroEmpleoSiradig(liquidacion *structLiquidacion.Liquidacion, columnaimportegananciasotroempleosiradig string, db *gorm.DB) float64 {
	var importeTotal float64
	anoLiquidacion := liquidacion.Fechaperiodoliquidacion.Format("2006")
	mesLiquidacion := liquidacion.Fechaperiodoliquidacion.Format("01")
	legajoid := strconv.Itoa(*liquidacion.Legajoid)
	sql := "SELECT SUM(" + columnaimportegananciasotroempleosiradig + ") FROM importegananciasotroempleosiradig WHERE '" + anoLiquidacion + "' = extract(YEAR from mes) and '" + mesLiquidacion + "' = extract(MONTH from mes) " +
		"and siradigid in (SELECT id from siradig where legajoid = " + legajoid + " ) AND importegananciasotroempleosiradig.deleted_at IS NULL"
	db.Raw(sql).Row().Scan(&importeTotal)

	return importeTotal
}

func getfgRemuneracionBrutaOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "importegananciasbrutas", db)
	fmt.Println("Calculos Automaticos - Remuneracion Bruta Otros Empleos:", importeTotal)
	return importeTotal
}

func getfgRemuneracionNoHabitualOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "importeretribucionesnohabituales", db)
	fmt.Println("Calculos Automaticos - Remuneracion No Habitual Otros Empleos:", importeTotal)
	return importeTotal
}

func getfgSacPrimerCuotaOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0
	if getfgMes(&liquidacion.Fechaperiodoliquidacion) <= 6 {
		importeTotal = getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "sac", db)
	}
	fmt.Println("Calculos Automaticos - Sac Primer Cuota Otros Empleos:", importeTotal)
	return importeTotal
}

func getfgSacSegundaCuotaOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0
	if getfgMes(&liquidacion.Fechaperiodoliquidacion) > 6 {
		importeTotal = getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "sac", db)
	}
	fmt.Println("Calculos Automaticos - Sac Segunda Cuota Otros Empleos:", importeTotal)
	return importeTotal
}

func getfgHorasExtrasGravadasOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "importehorasextrasgravadas", db)
	fmt.Println("Calculos Automaticos - Horas Extras Gravadas Otros Empleos:", importeTotal)
	return importeTotal
}

func getfgMovilidadYViaticosGravadaOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "gastosmovilidad", db)
	fmt.Println("Calculos Automaticos - Movilidad y Viaticos Gravada Otros Empleos:", importeTotal)
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
	arraySubtotalIngresos = append(arraySubtotalIngresos, GetfgSacPrimerCuota(liquidacion, db))
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
	fmt.Println("Calculos Automaticos - Subtotal Ingresos:", subtotalIngresos)
	return subtotalIngresos
}

func getfgAportesJubilatoriosRetirosPensionesOSubsidios(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("APORTES_JUBILATORIOS_RETIROS_PENSIONES_O_SUBSIDIOS", liquidacion, db)
	fmt.Println("Calculos Automaticos - Aportes Jubilatorios Retiros, Pensiones o Subsidios:", importeTotal)
	return importeTotal
}

func getfgAportesObraSocial(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("APORTES_OBRA_SOCIAL", liquidacion, db)
	fmt.Println("Calculos Automaticos - Aportes Obra Social:", importeTotal)
	return importeTotal
}

func getfgCuotaSindical(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("CUOTA_SINDICAL", liquidacion, db)
	fmt.Println("Calculos Automaticos - Cuota Sindical:", importeTotal)
	return importeTotal
}

func getfgDescuentosObligatoriosPorLeyNacionalProvincialOMunicipal(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSegunTipoImpuestoGanancias("DESCUENTOS_OBLIGATORIOS_POR_LEY_NACIONAL_PROVINCIAL_MUNICIPAL", liquidacion, db)
	fmt.Println("Calculos Automaticos - Descuentos Obligatorios por ley Nacional, Provincial o Municipal:", importeTotal)
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
	fmt.Println("Calculos Automaticos - Gastos Movilidad Viaticos Abonados por el Empleador:", importeTotal)
	return importeTotal

}

func getfgAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "aporteseguridadsocial", db)
	fmt.Println("Calculos Automaticos - Aportes Jubilatorios Retiros Pensiones o Subsidios Otros Empleos:", importeTotal)
	return importeTotal
}

func getfgAportesObraSocialOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "aporteobrasocial", db)
	fmt.Println("Calculos Automaticos - Aportes Obra Social Otros Empleos:", importeTotal)
	return importeTotal
}

func getfgCuotaSindicalOtrosEmpleos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteGananciasOtroEmpleoSiradig(liquidacion, "aportesindical", db)
	fmt.Println("Calculos Automaticos - Cuota Sindical Otros Empleos:", importeTotal)
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
	mesliquidacion := getfgMes(&liquidacion.Fechaperiodoliquidacion)
	anoliquidacion := liquidacion.Fechaperiodoliquidacion.Year()

	sql := "SELECT SUM(" + columnadeducciondesgravacionsiradig + ") FROM " + nombretablasiradig + " ts INNER JOIN siradigtipogrilla stg ON stg.id = ts.siradigtipogrillaid INNER JOIN siradig sdg on sdg.id = ts.siradigid WHERE to_number(to_char(mes, 'MM'),'99') <= " + strconv.Itoa(mesliquidacion) + " AND stg.codigo = '" + tipodeducciondesgravacionsiradig + "' AND sdg.legajoid = " + strconv.Itoa(*liquidacion.Legajoid) + " AND EXTRACT(year from sdg.periodosiradig) ='" + strconv.Itoa(anoliquidacion) + "' AND + ts.deleted_at IS  NULL AND stg.deleted_at IS NULL AND sdg.deleted_at IS NULL;"
	db.Raw(sql).Row().Scan(&importeTotal)

	return importeTotal
}

func getfgImporteTotalSiradigSegunTipoGrillaSinMes(liquidacion *structLiquidacion.Liquidacion, columnadeducciondesgravacionsiradig string, tipodeducciondesgravacionsiradig string, nombretablasiradig string, db *gorm.DB) float64 {
	var importeTotal float64
	anoliquidacion := liquidacion.Fechaperiodoliquidacion.Year()

	sql := "SELECT SUM(" + columnadeducciondesgravacionsiradig + ") FROM " + nombretablasiradig + " ts INNER JOIN siradigtipogrilla stg ON stg.id = ts.siradigtipogrillaid INNER JOIN siradig sdg on sdg.id = ts.siradigid WHERE stg.codigo = '" + tipodeducciondesgravacionsiradig + "' AND sdg.legajoid = " + strconv.Itoa(*liquidacion.Legajoid) + " AND EXTRACT(year from sdg.periodosiradig) ='" + strconv.Itoa(anoliquidacion) + "' AND stg.deleted_at IS NULL AND sdg.deleted_at IS NULL AND ts.deleted_at;"
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
	db.Raw(sql).Row().Scan(&importeTope)

	return importeTope
}

func getfgGastosSepelio(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "GASTOS_DE_SEPELIO", "deducciondesgravacionsiradig", db)
	importeTope := getfgValorFijoImpuestoGanancia(liquidacion, "topemaximodescuento", "topesepelio", db)
	importeTotal = getfgImporteTotalTope(importeTotal, importeTope)
	fmt.Println("Calculos Automaticos - Gastos Sepelio:", importeTotal)
	return importeTotal
}

func getfgGastosAmortizacionEInteresesRodadoCorredoresViajantesComercio(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "GASTOS_DE_REPRESENTACION_E_INTERESES_DE_CORREDORES_Y_VIAJANTES_DE_COMERCIO", "deducciondesgravacionsiradig", db)
	fmt.Println("Calculos Automaticos - Gastos Amortizacion e Intereses Corredores y Viajantes de Comercio:", importeTotal)
	return importeTotal
}

func getfgInteresesCreditosHipotecarios(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "INTERESES_PRESTAMO_HIPOTECARIO", "deducciondesgravacionsiradig", db)
	importeTope := getfgValorFijoImpuestoGanancia(liquidacion, "topemaximodescuento", "topehipotecarios", db)
	importeTotal = getfgImporteTotalTope(importeTotal, importeTope)
	fmt.Println("Calculos Automaticos - Intereses Creditos Hipotecarios:", importeTotal)
	return importeTotal
}

func getfgAportesCapSocFondoRiesgoSociosProtectoresSGR(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "Montoreintegrar + Montoreintegrar3", "REINTEGRO_DE_APORTES_DE_SOCIOS_PROTECTORES_A_SOCIEDADES_DE_GARANTIA_RECIPROCA", "ajustesiradig", db)
	fmt.Println("Calculos Automaticos - Aportes Cap. Soc. Fondo Riesgo Socios Protectores SGR:", importeTotal)
	return importeTotal
}

func getfgAlquileresInmueblesDestinadosASuCasaHabitacion(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "ALQUILER_INMUEBLES_DESTINADOS_A_CASA_HABITACION", "deducciondesgravacionsiradig", db)
	importeTope := getfgMinimoNoImponible(liquidacion, db) * 0.4 /*es el 40% de MNI(40)*/
	importeTotal = getfgImporteTotalTope(importeTotal, importeTope)
	fmt.Println("Calculos Automaticos - Alquileres Inmuebles Casa Habitacion:", importeTotal)
	return importeTotal
}

func getfgEmpleadosServicioDomestico(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "contribucion + retribucion", "DEDUCCION_DEL_PERSONAL_DOMESTICO", "deducciondesgravacionsiradig", db)
	importeTope := getfgMinimoNoImponible(liquidacion, db) /*es el MNI(40)*/
	importeTotal = getfgImporteTotalTope(importeTotal, importeTope)
	fmt.Println("Calculos Automaticos - Empleados Servicio Domestico:", importeTotal)
	return importeTotal
}

func getfgIndumentariaEquipamientoCaracterObligatorio(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "GASTOS_ADQUISICION_INDUMENTARIA_Y_EQUIPAMIENTO_PARA_USO_EXCLUSIVO_EN_EL_LUGAR_DE_TRABAJO", "deducciondesgravacionsiradig", db)
	fmt.Println("Calculos Automaticos - Indumentaria Equipamiento Caracter Obligatorio:", importeTotal)
	return importeTotal
}

func getfgOtrasDeducciones(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrilla(liquidacion, "importe", "OTRAS", "deducciondesgravacionsiradig", db)
	fmt.Println("Calculos Automaticos - Otras Deducciones:", importeTotal)
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
	fmt.Println("Calculos Automaticos - Subtotal:", subTotal)
	return subTotal
}

/*Estos dos se utilizan sin mes ya que se toma el acumulado anual y luego se le saca el tope, Cualquier cosa consultar con DIEGO*/

func getfgCuotaMedicoAsistencial(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrillaSinMes(liquidacion, "importe", "CUOTA_MEDICA_ASISTENCIAL", "deducciondesgravacionsiradig", db)
	var importeTope float64
	if importeTotal != 0 {
		importeTope = getfgSubtotal(liquidacion, db) * 0.05 //5% de Subtotal
	}

	importeTotal = getfgImporteTotalTope(importeTotal, importeTope)
	fmt.Println("Calculos Automaticos - Cuota Medico Asistencial:", importeTotal)
	return importeTotal
}

func getfgDonacionFiscosNacProvMunArt20(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgImporteTotalSiradigSegunTipoGrillaSinMes(liquidacion, "importe", "DONACIONES", "deducciondesgravacionsiradig", db)
	var importeTope float64
	if importeTotal != 0 {
		importeTope = getfgSubtotal(liquidacion, db) * 0.05 //5% de Subtotal
	}
	importeTotal = getfgImporteTotalTope(importeTotal, importeTope)
	fmt.Println("Calculos Automaticos - Donacion Fisico Nac, Prov, Munic art. 20:", importeTotal)
	return importeTotal
}

/**/

func getfgGananciaNeta(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arrayGananciaNeta []float64
	var gananciaNeta float64

	arrayGananciaNeta = append(arrayGananciaNeta, getfgCuotaMedicoAsistencial(liquidacion, db))
	arrayGananciaNeta = append(arrayGananciaNeta, getfgDonacionFiscosNacProvMunArt20(liquidacion, db))

	gananciaNeta = getfgSubtotal(liquidacion, db) - Sum(arrayGananciaNeta)
	fmt.Println("Calculos Automaticos - Ganancia Neta:", gananciaNeta)
	return gananciaNeta
}

func getfgDetalleCargoFamiliar(liquidacion *structLiquidacion.Liquidacion, columnaDetalleCargoFamiliar string, valorfijocolumna string, porcentaje float64, db *gorm.DB) float64 {
	var importeTotal float64
	var tienevalorbeneficio bool
	anioperiodoliquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	mesperiodoliquidacion := getfgMes(&liquidacion.Fechaperiodoliquidacion)

	var detallecargofamiliar structSiradig.Detallecargofamiliarsiradig
	sql := "SELECT dcfs.* FROM siradig s INNER JOIN detallecargofamiliarsiradig dcfs ON s.id = dcfs.siradigid where to_char(periodosiradig, 'YYYY') = '" + strconv.Itoa(anioperiodoliquidacion) + "' AND dcfs." + columnaDetalleCargoFamiliar + " NOTNULL AND s.legajoid = " + strconv.Itoa(*liquidacion.Legajoid) + " AND s.deleted_at IS NULL AND dcfs.deleted_at IS NULL"
	db.Raw(sql).Scan(&detallecargofamiliar)
	sql = "SELECT valor FROM siradig s INNER JOIN beneficiosiradig bs ON s.id = bs.siradigid WHERE to_number(to_char(bs.mesdesde, 'MM'),'99') <= " + strconv.Itoa(mesperiodoliquidacion) + " AND to_number(to_char(bs.meshasta, 'MM'), '99') > " + strconv.Itoa(mesperiodoliquidacion) + " AND bs.siradigtipogrillaid = -24 AND to_char(s.periodosiradig, 'YYYY') = '" + strconv.Itoa(anioperiodoliquidacion) + "' AND s.deleted_at IS NULL AND bs.deleted_at IS NULL"
	db.Raw(sql).Row().Scan(&tienevalorbeneficio)

	if detallecargofamiliar.ID != 0 {

		mesdadobaja := getfgMes(detallecargofamiliar.Meshasta)
		mesdadoalta := getfgMes(detallecargofamiliar.Mesdesde)
		valorfijo := getfgValorFijoImpuestoGanancia(liquidacion, "deduccionespersonales", valorfijocolumna, db)

		if tienevalorbeneficio == true {
			valorfijo = valorfijo * 1.22
		}

		if mesdadobaja == 0 {
			importeTotal = (valorfijo / 12) * float64(mesperiodoliquidacion-(mesdadoalta-1)) * (porcentaje / 100)
		} else {
			if mesdadobaja <= mesperiodoliquidacion {
				importeTotal = (valorfijo / 12) * float64(mesdadobaja-mesdadoalta) * (porcentaje / 100)
			} else {
				if mesdadobaja > mesperiodoliquidacion {
					importeTotal = (valorfijo / 12) * float64(mesperiodoliquidacion-(mesdadoalta-1)) * (porcentaje / 100)
				}
			}
		}
	}

	return importeTotal
}

func getfgConyuge(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgDetalleCargoFamiliar(liquidacion, "conyugeid", "valorfijoconyuge", 100, db)
	fmt.Println("Calculos Automaticos - Conyuge:", importeTotal)
	return importeTotal
}

func getfgHijos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0
	var detallescargofamiliarsiradig []structSiradig.Detallecargofamiliarsiradig

	valorfijoMNI := getfgValorFijoImpuestoGanancia(liquidacion, "deduccionespersonales", "valorfijomni", db)
	sql := "SELECT * FROM detallecargofamiliarsiradig WHERE hijoid NOTNULL AND estaacargo = true AND montoanual < " + strconv.FormatFloat(valorfijoMNI, 'f', 5, 64) + "AND detallecargofamiliarsiradig.deleted_at IS NULL"
	db.Raw(sql).Scan(&detallescargofamiliarsiradig)

	for i := 0; i < len(detallescargofamiliarsiradig); i++ {
		porcentaje := detallescargofamiliarsiradig[i].Porcentaje
		importeTotal = importeTotal + getfgDetalleCargoFamiliar(liquidacion, "hijoid", "valorfijohijo", *porcentaje, db)
	}
	fmt.Println("Calculos Automaticos - Hijos:", importeTotal)
	return importeTotal
}

func getfgMinimoNoImponible(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	valorfijoMNI := getfgValorFijoImpuestoGanancia(liquidacion, "deduccionespersonales", "valorfijomni", db)
	mesperiodoliquidacion := getfgMes(&liquidacion.Fechaperiodoliquidacion)
	importeTotal := (valorfijoMNI / 12) * float64(mesperiodoliquidacion)
	fmt.Println("Calculos Automaticos - Minimo No Imponible:", importeTotal)
	return importeTotal
}

func getfgDeduccionEspecial(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	valorfijoMNI := getfgValorFijoImpuestoGanancia(liquidacion, "deduccionespersonales", "valorfijoddei", db)
	mesperiodoliquidacion := getfgMes(&liquidacion.Fechaperiodoliquidacion)
	importeTotal := (valorfijoMNI / 12) * float64(mesperiodoliquidacion)
	fmt.Println("Calculos Automaticos - Deduccion Especial:", importeTotal)
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
	fmt.Println("Calculos Automaticos - Subtotal Deducciones Personales:", subTotalDeduccionesPersonales)
	return subTotalDeduccionesPersonales
}

func getfgDeduccionesAComputar(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgSubtotalDeduccionesPersonales(liquidacion, db)
	fmt.Println("Calculos Automaticos - Deducciones a Computar:", importeTotal)
	return importeTotal
}

func obtenerLiquidacionesIgualAnioLegajoMenorMes(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) *[]structLiquidacion.Liquidacion {
	var liquidaciones []structLiquidacion.Liquidacion
	anioperiodoliquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	mesliquidacion := getfgMes(&liquidacion.Fechaperiodoliquidacion)
	db.Set("gorm:auto_preload", true).Find(&liquidaciones, "to_number(to_char(fechaperiodoliquidacion, 'MM'),'99') < ? AND to_char(fechaperiodoliquidacion, 'YYYY') = ? AND legajoid = ?", mesliquidacion, strconv.Itoa(anioperiodoliquidacion), *liquidacion.Legajoid)

	return &liquidaciones
}

func getfgGananciaNetaAcumSujetaAImp(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0

	liquidaciones := *obtenerLiquidacionesIgualAnioLegajoMenorMes(liquidacion, db)
	importeTotal = importeTotal + getfgGananciaNeta(liquidacion, db)

	for i := 0; i < len(liquidaciones); i++ {
		importeTotal = importeTotal + getfgGananciaNeta(&liquidaciones[i], db)
	}
	fmt.Println("Calculos Automaticos - Ganancia Neta Acum. Sujeta a Impuestos:", importeTotal)
	return importeTotal
}

func getfgDeduccionesPersonales(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getfgDeduccionesAComputar(liquidacion, db)
	fmt.Println("Calculos Automaticos - Deducciones Personales:", importeTotal)
	return importeTotal * -1
}

func getfgBaseImponible(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arrayBaseImponible []float64
	var importeTotal float64

	arrayBaseImponible = append(arrayBaseImponible, getfgGananciaNetaAcumSujetaAImp(liquidacion, db))
	arrayBaseImponible = append(arrayBaseImponible, getfgDeduccionesPersonales(liquidacion, db))

	importeTotal = Sum(arrayBaseImponible)
	fmt.Println("Calculos Automaticos - Base Imponible:", importeTotal)
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
	var porcentaje float64 = 1
	var porcentajeOtrosEmp float64 = 1

	if remunerativosMenosDescuentos != 0 {
		porcentaje = uno - (cuotaSindical/remunerativosMenosDescuentos + obraSocial/remunerativosMenosDescuentos + aportesJubilatorios/remunerativosMenosDescuentos)
	}
	if remunerativosOtros != 0 {
		porcentajeOtrosEmp = uno - (cuotaSindicalOtros/remunerativosOtros + obraSocialOtros/remunerativosOtros + aportesJubilatoriosOtros/remunerativosOtros)
	}

	var importeTotalHorasExtrasGravadas float64 = 0
	var importeTotalHorasExtrasGravadasOtrosEmpleos float64 = 0

	liquidaciones := *obtenerLiquidacionesIgualAnioLegajoMenorMes(liquidacion, db)

	importeTotalHorasExtrasGravadas = importeTotalHorasExtrasGravadas + getfgHorasExtrasGravadas(liquidacion, db)
	importeTotalHorasExtrasGravadasOtrosEmpleos = importeTotalHorasExtrasGravadasOtrosEmpleos + getfgHorasExtrasGravadasOtrosEmpleos(liquidacion, db)

	for i := 0; i < len(liquidaciones); i++ {
		importeTotalHorasExtrasGravadas = importeTotalHorasExtrasGravadas + getfgHorasExtrasGravadas(&liquidaciones[i], db)
		importeTotalHorasExtrasGravadasOtrosEmpleos = importeTotalHorasExtrasGravadasOtrosEmpleos + getfgHorasExtrasGravadasOtrosEmpleos(&liquidaciones[i], db)
	}

	importeTotal := getfgBaseImponible(liquidacion, db) - (importeTotalHorasExtrasGravadas*porcentaje + importeTotalHorasExtrasGravadasOtrosEmpleos*porcentajeOtrosEmp)
	fmt.Println("Calculos Automaticos - Total Ganancia Neta Imponible Acumulada sin Horas Extras:", importeTotal)
	fmt.Println("Calculos Automaticos - Importe Horas Extras Gravadas en Total Ganancia Neta Acum:", importeTotalHorasExtrasGravadas)
	fmt.Println("Calculos Automaticos - Importe Horas Extras Gravadas Otros Empleos en Total Ganancia Neta Acum:", importeTotalHorasExtrasGravadasOtrosEmpleos)
	return importeTotal
}

func obtenerRemunerativosMenosDescuentos(liquidacion *structLiquidacion.Liquidacion) float64 {
	var totalRemunerativos, totalDescuentos float64
	for i := 0; i < len(liquidacion.Liquidacionitems); i++ {
		liquidacionitem := liquidacion.Liquidacionitems[i]
		tipoconcepto := *liquidacionitem.Concepto.Tipoconceptoid
		importeconcepto := liquidacionitem.Importeunitario
		if importeconcepto != nil {

			if tipoconcepto == -1 {
				totalRemunerativos = totalRemunerativos + *importeconcepto
			}
			if tipoconcepto == -3 {
				totalDescuentos = totalDescuentos + *importeconcepto
			}
		}
	}
	fmt.Println("Calculos Automaticos - RemunerativosMenosDescuentos:", totalRemunerativos-totalDescuentos)
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
	fmt.Println("Calculos Automaticos - Remunerativos Otros:", totalRemunerativosOtros)
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
	mesLiquidacion := s.Split(liquidacion.Fechaperiodoliquidacion.String(), "-")[1]

	mesAnioLiquidacion := mesLiquidacion + "/" + strconv.Itoa(anioLiquidacion)

	sql := "SELECT limiteinferior,limitesuperior,valorfijo,valorvariable,mesanio FROM escalaimpuestoaplicable where mesanio = '" + mesAnioLiquidacion + "' and escalaimpuestoaplicable.deleted_at IS NULL"
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
	fmt.Println("Calculos Automaticos - Determinacion Impuesto Fijo:", importeTotal)
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
	fmt.Println("Calculos Automaticos - Determinacion Impuesto por Escala:", importeTotal)
	return importeTotal
}

func getfgTotalRetener(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arrayTotalRetener []float64
	var totalRetener float64

	arrayTotalRetener = append(arrayTotalRetener, getfgDeterminacionImpuestoFijo(liquidacion, db))
	arrayTotalRetener = append(arrayTotalRetener, getfgDeterminacionImpuestoPorEscala(liquidacion, db))

	totalRetener = Sum(arrayTotalRetener)
	fmt.Println("Calculos Automaticos - Total a Retener:", totalRetener)
	return totalRetener
}

func getfgRetencionAcumulada(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	anioperiodoliquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	mesliquidacion := getfgMes(&liquidacion.Fechaperiodoliquidacion)
	var totalconceptosimpuestoganancias float64
	var totalconceptosimpuestogananciasdevolucion float64

	sql := "SELECT SUM(li.importeunitario) FROM liquidacion l INNER JOIN liquidacionitem li ON l.id = li.liquidacionid INNER JOIN legajo le ON le.id = l.legajoid INNER JOIN concepto c ON c.id = li.conceptoid WHERE to_number(to_char(l.fechaperiodoliquidacion, 'MM'),'99') < " + strconv.Itoa(mesliquidacion) + " AND to_char(l.fechaperiodoliquidacion, 'YYYY') = '" + strconv.Itoa(anioperiodoliquidacion) + "' AND le.id = " + strconv.Itoa(*liquidacion.Legajoid) + " AND c.codigo = 'IMPUESTO_GANANCIAS' AND l.deleted_at IS NULL AND le.deleted_at IS NULL AND li.deleted_at IS NULL AND c.deleted_at IS NULL"
	db.Raw(sql).Row().Scan(&totalconceptosimpuestoganancias)

	sql = "SELECT SUM(li.importeunitario) FROM liquidacion l INNER JOIN liquidacionitem li ON l.id = li.liquidacionid INNER JOIN legajo le ON le.id = l.legajoid INNER JOIN concepto c ON c.id = li.conceptoid WHERE to_number(to_char(l.fechaperiodoliquidacion, 'MM'),'99') < " + strconv.Itoa(mesliquidacion) + " AND to_char(l.fechaperiodoliquidacion, 'YYYY') = '" + strconv.Itoa(anioperiodoliquidacion) + "' AND le.id = " + strconv.Itoa(*liquidacion.Legajoid) + " AND c.codigo = 'IMPUESTO_GANANCIAS_DEVOLUCION' AND l.deleted_at IS NULL AND le.deleted_at IS NULL AND li.deleted_at IS NULL AND c.deleted_at IS NULL"
	db.Raw(sql).Row().Scan(&totalconceptosimpuestogananciasdevolucion)
	fmt.Println("Calculos Automaticos - Retencion acumulada:", totalconceptosimpuestoganancias-totalconceptosimpuestogananciasdevolucion)
	return totalconceptosimpuestoganancias - totalconceptosimpuestogananciasdevolucion
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
