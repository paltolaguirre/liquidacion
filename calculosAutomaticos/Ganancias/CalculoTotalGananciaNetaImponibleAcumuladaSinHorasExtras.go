package Ganancias

type CalculoTotalGananciaNetaImponibleAcumuladaSinHorasExtras struct {
	CalculoGanancias

}



func (cg *CalculoTotalGananciaNetaImponibleAcumuladaSinHorasExtras) getResultInternal() float64 {

	cuotaSindical := (&CalculoCuotaSindical{cg.CalculoGanancias}).getResult()
	obraSocial := (&CalculoAportesObraSocial{cg.CalculoGanancias}).getResult()
	aportesJubilatorios := (&CalculoAportesJubilatoriosRetirosPensionesOSubsidios{cg.CalculoGanancias}).getResult()
	remunerativosMenosDescuentos := cg.obtenerRemunerativosMenosDescuentos()
	cuotaSindicalOtros := (&CalculoCuotaSindicalOtrosEmpleos{cg.CalculoGanancias}).getResult()
	obraSocialOtros := (&CalculoAportesObraSocialOtrosEmpleos{cg.CalculoGanancias}).getResult()
	aportesJubilatoriosOtros := (&CalculoAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos{cg.CalculoGanancias}).getResult()
	remunerativosOtros := cg.obtenerRemunerativosOtros()

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

	liquidaciones := *cg.obtenerLiquidacionesIgualAnioLegajoMenorMes()

	importeTotalHorasExtrasGravadas = importeTotalHorasExtrasGravadas + (&CalculoHorasExtrasGravadas{cg.CalculoGanancias}).getResult()
	importeTotalHorasExtrasGravadasOtrosEmpleos = importeTotalHorasExtrasGravadasOtrosEmpleos + (&CalculoHorasExtrasGravadasOtrosEmpleos{cg.CalculoGanancias}).getResult()

	for i := 0; i < len(liquidaciones); i++ {
		itemGanancias := obtenerItemGananciaFromLiquidacion(&liquidaciones[i])
		calculoGananciasAnterior := CalculoGanancias{itemGanancias, &liquidaciones[i], cg.Db}
		importeTotalHorasExtrasGravadas = importeTotalHorasExtrasGravadas + (&CalculoHorasExtrasGravadas{calculoGananciasAnterior}).getResult()
		importeTotalHorasExtrasGravadasOtrosEmpleos = importeTotalHorasExtrasGravadasOtrosEmpleos + (&CalculoHorasExtrasGravadasOtrosEmpleos{calculoGananciasAnterior}).getResult()
	}

	importeTotal := (&CalculoBaseImponible{cg.CalculoGanancias}).getResult() - (importeTotalHorasExtrasGravadas*porcentaje + importeTotalHorasExtrasGravadasOtrosEmpleos*porcentajeOtrosEmp)
	return importeTotal
}

func (cg *CalculoTotalGananciaNetaImponibleAcumuladaSinHorasExtras) getResult() float64 {
	//cg.CalculoGanancias.formula = cg
	return cg.getResultOnDemandTemplate("Total Ganancia Neta Imponible Acumulada Sin Horas Extras", "TOTAL_GANANCIA_NETA_IMPONIBLE_ACUMULADA_SIN_HORAS_EXTRAS", 48, cg)
}

func (cg *CalculoTotalGananciaNetaImponibleAcumuladaSinHorasExtras) getTope() *float64 {
	return nil
}