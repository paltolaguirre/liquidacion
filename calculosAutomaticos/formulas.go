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

func getConyuge(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getDetalleCargoFamiliar(liquidacion, "conyugeid", "valorfijoconyuge", 1, db)
	return importeTotal
}

func getHijos(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0
	var detallescargofamiliarsiradig []structSiradig.Detallecargofamiliarsiradig

	valorfijoMNI := getValorFijoImpuestoGanancia(liquidacion, "deduccionespersonales", "valorfijomni", db)
	sql := "SELECT * FROM detallecargofamiliar WHERE hijoid NOT NULL AND estaacargo = true AND montoanual < " + strconv.FormatFloat(valorfijoMNI, 'f', 5, 64)
	db.Raw(sql).Scan(&detallescargofamiliarsiradig)

	for i := 0; i < len(detallescargofamiliarsiradig); i++ {
		porcentaje := detallescargofamiliarsiradig[i].Porcentaje
		importeTotal = importeTotal + getDetalleCargoFamiliar(liquidacion, "hijoid", "valorfijohijo", *porcentaje, db)
	}

	return importeTotal
}

func getMinimoNoImponible(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	valorfijoMNI := getValorFijoImpuestoGanancia(liquidacion, "deduccionespersonales", "valorfijomni", db)
	mesperiodoliquidacion := getMes(&liquidacion.Fechaperiodoliquidacion)
	importeTotal := (valorfijoMNI / 12) * float64(mesperiodoliquidacion)
	return importeTotal
}

func getDeduccionEspecial(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	valorfijoMNI := getValorFijoImpuestoGanancia(liquidacion, "deduccionespersonales", "valorfijoddei", db)
	mesperiodoliquidacion := getMes(&liquidacion.Fechaperiodoliquidacion)
	importeTotal := (valorfijoMNI / 12) * float64(mesperiodoliquidacion)
	return importeTotal
}

func getSubtotalDeduccionesPersonales(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arraySubtotalDeduccionesPersonales []float64
	var subTotalDeduccionesPersonales float64

	arraySubtotalDeduccionesPersonales = append(arraySubtotalDeduccionesPersonales, getConyuge(liquidacion, db))
	arraySubtotalDeduccionesPersonales = append(arraySubtotalDeduccionesPersonales, getHijos(liquidacion, db))
	arraySubtotalDeduccionesPersonales = append(arraySubtotalDeduccionesPersonales, getMinimoNoImponible(liquidacion, db))
	arraySubtotalDeduccionesPersonales = append(arraySubtotalDeduccionesPersonales, getDeduccionEspecial(liquidacion, db))

	subTotalDeduccionesPersonales = Sum(arraySubtotalDeduccionesPersonales)
	return subTotalDeduccionesPersonales
}

func getDeduccionesAComputar(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getSubtotalDeduccionesPersonales(liquidacion, db)
	return importeTotal
}

func getGananciaNetaAcumSujetaAImp(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	//Ver con Rodri como manejar este caso
	return 0.6
}

func getDeduccionesPersonales(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	importeTotal := getDeduccionesAComputar(liquidacion, db)
	return importeTotal * -1
}

func getBaseImponible(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arrayBaseImponible []float64
	var importeTotal float64

	arrayBaseImponible = append(arrayBaseImponible, getGananciaNetaAcumSujetaAImp(liquidacion, db))
	arrayBaseImponible = append(arrayBaseImponible, getDeduccionesPersonales(liquidacion, db))

	importeTotal = Sum(arrayBaseImponible)
	return importeTotal
}

func getTotalGananciaNetaImponibleAcumuladaSinHorasExtras(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {

	cuotaSindical := getCuotaSindical(liquidacion, db)
	obraSocial := getAportesObraSocial(liquidacion, db)
	aportesJubilatorios := getAportesJubilatoriosRetirosPensionesOSubsidios(liquidacion, db)
	remunerativosMenosDescuentos := obtenerRemunerativosMenosDescuentos(liquidacion)
	cuotaSindicalOtros := getCuotaSindicalOtrosEmpleos(liquidacion, db)
	obraSocialOtros := getAportesObraSocialOtrosEmpleos(liquidacion, db)
	aportesJubilatoriosOtros := getAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos(liquidacion, db)
	remunerativosOtros := obtenerRemunerativosOtros(liquidacion, db)

	var uno float64 = 1
	porcentaje := uno - (cuotaSindical/remunerativosMenosDescuentos + obraSocial/remunerativosMenosDescuentos + aportesJubilatorios/remunerativosMenosDescuentos)
	porcentajeOtrosEmp := uno - (cuotaSindicalOtros/remunerativosOtros + obraSocialOtros/remunerativosOtros + aportesJubilatoriosOtros/remunerativosOtros)

	importeTotal := getBaseImponible(liquidacion, db) - (getHorasExtrasGravadas(liquidacion, db)*porcentaje + getHorasExtrasGravadasOtrosEmpleos(liquidacion, db)*porcentajeOtrosEmp)
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

	arrayRemunerativosOtros = append(arrayRemunerativosOtros, getRemuneracionBrutaOtrosEmpleos(liquidacion, db))
	arrayRemunerativosOtros = append(arrayRemunerativosOtros, getSacPrimerCuotaOtrosEmpleos(liquidacion, db))
	arrayRemunerativosOtros = append(arrayRemunerativosOtros, getSacSegundaCuotaOtrosEmpleos(liquidacion, db))
	arrayRemunerativosOtros = append(arrayRemunerativosOtros, getHorasExtrasGravadasOtrosEmpleos(liquidacion, db))

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

func getEscalaImpuestoAplicable(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) *[]strEscalaimpuestoaplicable {
	var strescalaimpuestoaplicable []strEscalaimpuestoaplicable

	anioLiquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	mesLiquidacion := getMes(&liquidacion.Fechaperiodoliquidacion)

	mesAnioLiquidacion := strconv.Itoa(mesLiquidacion) + "/" + strconv.Itoa(anioLiquidacion)

	sql := "SELECT limiteinferior,limitesuperior,valorfijo,valorvariable,mesanio FROM escalaimpuestoaplicable where mesanio = '" + mesAnioLiquidacion + "'"
	db.Raw(sql).Scan(&strescalaimpuestoaplicable)

	return &strescalaimpuestoaplicable

}

func getDeterminacionImpuestoFijo(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0
	strescalaimpuestoaplicable := *getEscalaImpuestoAplicable(liquidacion, db)
	totalganancianeta := getTotalGananciaNetaImponibleAcumuladaSinHorasExtras(liquidacion, db)

	for i := 0; i < len(strescalaimpuestoaplicable); i++ {
		escalaimpuestoaplicable := strescalaimpuestoaplicable[i]
		if totalganancianeta > escalaimpuestoaplicable.Limiteinferior && totalganancianeta <= escalaimpuestoaplicable.Limitesuperior {
			importeTotal = escalaimpuestoaplicable.Valorfijo
		}
	}
	return importeTotal
}

func getDeterminacionImpuestoPorEscala(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var importeTotal float64 = 0
	strescalaimpuestoaplicable := *getEscalaImpuestoAplicable(liquidacion, db)
	totalganancianeta := getTotalGananciaNetaImponibleAcumuladaSinHorasExtras(liquidacion, db)
	baseimponible := getBaseImponible(liquidacion, db)
	for i := 0; i < len(strescalaimpuestoaplicable); i++ {
		escalaimpuestoaplicable := strescalaimpuestoaplicable[i]
		if totalganancianeta > escalaimpuestoaplicable.Limiteinferior && totalganancianeta <= escalaimpuestoaplicable.Limitesuperior {

			importeTotal = (baseimponible - escalaimpuestoaplicable.Limiteinferior) * escalaimpuestoaplicable.Valorvariable
		}
	}
	return importeTotal
}

func getTotalRetener(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	var arrayTotalRetener []float64
	var totalRetener float64

	arrayTotalRetener = append(arrayTotalRetener, getDeterminacionImpuestoFijo(liquidacion, db))
	arrayTotalRetener = append(arrayTotalRetener, getDeterminacionImpuestoPorEscala(liquidacion, db))

	totalRetener = Sum(arrayTotalRetener)
	return totalRetener
}

func getRetencionAcumulada(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	anioperiodoliquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	var totalconceptosimpuestoganancias float64
	sql := "SELECT SUM(li.importeunitario) FROM siradig s INNER JOIN legajo le ON le.id = s.legajoid INNER JOIN liquidacion l ON le.id = l.legajoid INNER JOIN liquidacionitem li ON l.id = li.liquidacionid INNER JOIN  concepto c ON c.id = li.conceptoid WHERE c.codigo = 'IMPUESTO_GANANCIAS' AND to_char(s.periodosiradig, 'YYYY') = ' " + strconv.Itoa(anioperiodoliquidacion)
	db.Raw(sql).Row().Scan(&totalconceptosimpuestoganancias)

	return totalconceptosimpuestoganancias
}

func getRetencionMes(liquidacion *structLiquidacion.Liquidacion, db *gorm.DB) float64 {
	totalRetener := getTotalRetener(liquidacion, db)
	retencionAcumulada := getRetencionAcumulada(liquidacion, db)

	return totalRetener - retencionAcumulada
}

func getDetalleCargoFamiliar(liquidacion *structLiquidacion.Liquidacion, columnaDetalleCargoFamiliar string, valorfijocolumna string, porcentaje float64, db *gorm.DB) float64 {
	var importeTotal float64
	var tienevalorbeneficio bool
	anioperiodoliquidacion := liquidacion.Fechaperiodoliquidacion.Year()
	mesperiodoliquidacion := getMes(&liquidacion.Fechaperiodoliquidacion)

	var detallecargofamiliar *structSiradig.Detallecargofamiliarsiradig
	sql := "SELECT dcfs.* FROM siradig s INNER JOIN detallecargofamiliarsiradig dcfs ON s.id = dcfs.siradigid where to_char(periodosiradig, 'YYYY') = ' " + strconv.Itoa(anioperiodoliquidacion) + "' AND dcfs." + columnaDetalleCargoFamiliar + " NOTNULL AND s.legajoid = " + strconv.Itoa(*liquidacion.Legajoid)
	db.Raw(sql).Scan(&detallecargofamiliar)

	sql = "SELECT valor FROM siradig s INNER JOIN beneficiosiradig bs ON s.id = bs.siradigid WHERE '" + strconv.Itoa(mesperiodoliquidacion) + "' > to_char(bs.mesdesde, 'MM') AND '" + strconv.Itoa(mesperiodoliquidacion) + "' < to_char(bs.meshasta, 'MM') AND to_char(s.periodosiradig, 'YYYY') = ' " + strconv.Itoa(anioperiodoliquidacion) + "'"
	db.Raw(sql).Row().Scan(&tienevalorbeneficio)

	if detallecargofamiliar.ID != 0 {

		mesdadobaja := getMes(detallecargofamiliar.Meshasta)
		mesdadoalta := getMes(detallecargofamiliar.Mesdesde)
		valorfijo := getValorFijoImpuestoGanancia(liquidacion, "deduccionespersonales", valorfijocolumna, db)

		if tienevalorbeneficio == false {
			valorfijo = valorfijo * 1.22
		}

		if mesdadobaja == 0 {
			importeTotal = (valorfijo / 12) * float64(mesperiodoliquidacion-mesdadoalta) * porcentaje
		}

		if mesdadobaja <= mesperiodoliquidacion {
			importeTotal = (valorfijo / 12) * float64(mesdadobaja-mesdadoalta) * porcentaje
		}

		if mesdadobaja > mesperiodoliquidacion {
			importeTotal = (valorfijo / 12) * float64(mesperiodoliquidacion-mesdadoalta) * porcentaje
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
