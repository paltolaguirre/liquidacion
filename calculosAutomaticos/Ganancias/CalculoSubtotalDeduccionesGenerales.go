package Ganancias

type CalculoSubtotalDeduccionesGenerales struct {
	CalculoGanancias
}

func (cg *CalculoSubtotalDeduccionesGenerales) getResultInternal() float64 {
	var arraySubtotalDeduccionesGenerales []float64
	var importeTotal float64

	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoAportesJubilatoriosRetirosPensionesOSubsidios{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoAportesObraSocial{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoAportesObraSocialOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoCuotaSindical{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoCuotaSindicalOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoPrimasDeSeguroParaElCasoDeMuerteAnual{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoSeguroMuerteMixtosSujetosAlControlSSNAnual{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoSegurosRetirosPrivadosSujetosAlControlSSNAnual{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoAdquisicionDeCuotapartesDeFCIConFinesDeRetiro{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoGastosSepelio{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoGastosAmortizacionEInteresesRodadoCorredoresViajantesComercio{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoDescuentosObligatoriosPorLeyNacionalProvincialOMunicipal{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoHonorariosServAsistenciaSanitariaMedicaYParamedica{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoInteresesCreditosHipotecarios{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoAportesCapSocFondoRiesgoSociosProtectoresSGR{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoAportesCajasComplementariasFondosCompensadoresDePrevSimilares{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoAlquileresInmueblesDestinadosASuCasaHabitacion{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoEmpleadosServicioDomestico{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoGastosMovilidadViaticosAbonadosPorElEmpleador{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoIndumentariaEquipamientoCaracterObligatorio{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoOtrasDeducciones{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoOtrasDeduccionesActoresRetribucionPagadaALosRepresentantesRG244208{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesGenerales = append(arraySubtotalDeduccionesGenerales, (&CalculoOtrasDeduccionesFondosCompensadoresPrevision{cg.CalculoGanancias}).getResult())

	importeAcumuladorMesAnterior := cg.obtenerAcumuladorLiquidacionItemMesAnteriorSegunCodigo("SUBTOTAL_DEDUCCIONES_GENERALES")
	importeTotal = Sum(arraySubtotalDeduccionesGenerales) + importeAcumuladorMesAnterior
	return importeTotal
}

func (cg *CalculoSubtotalDeduccionesGenerales) getResult() float64 {
	return cg.getResultOnDemandTemplate("SUBTOTAL_DEDUCCIONES_GENERALES", 0, cg)
}

func (cg *CalculoSubtotalDeduccionesGenerales) getTope() *float64 {
	return nil
}

func (cg *CalculoSubtotalDeduccionesGenerales) getNombre() string {
	return "Subtotal Deducciones Generales"
}

func (cg *CalculoSubtotalDeduccionesGenerales) getEsMostrable() bool {
	return false
}
