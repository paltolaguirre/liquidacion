package Ganancias

type CalculoSubtotal struct{
	CalculoGanancias
}

func (cg *CalculoSubtotal) getResultInternal() float64{
	var arraySubtotal []float64

	arraySubtotal = append(arraySubtotal, (&CalculoAportesJubilatoriosRetirosPensionesOSubsidios{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoAportesObraSocial{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoCuotaSindical{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoDescuentosObligatoriosPorLeyNacionalProvincialOMunicipal{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoGastosMovilidadViaticosAbonadosPorElEmpleador{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoAportesObraSocialOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoCuotaSindicalOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoPrimasDeSeguroParaCasoDeMuerte{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoSeguroMuerteMixtosSujetosAlControlSSN{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoSegurosRetirosPrivadosSujetosAlControlSSN{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoGastosSepelio{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoGastosAmortizacionEInteresesRodadoCorredoresViajantesComercio{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoInteresesCreditosHipotecarios{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoAportesCapSocFondoRiesgoSociosProtectoresSGR{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoAlquileresInmueblesDestinadosASuCasaHabitacion{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoEmpleadosServicioDomestico{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoIndumentariaEquipamientoCaracterObligatorio{cg.CalculoGanancias}).getResult())
	arraySubtotal = append(arraySubtotal, (&CalculoOtrasDeducciones{cg.CalculoGanancias}).getResult())

	return (&CalculoSubtotalIngresos{cg.CalculoGanancias}).getResult() - Sum(arraySubtotal)
}

func (cg *CalculoSubtotal) getResult() float64{
	return cg.getResultOnDemandTemplate("SUBTOTAL", 35, cg)
}

func (cg *CalculoSubtotal) getTope() *float64 {
	return nil
}

func (cg *CalculoSubtotal) getNombre() string {
	return "Subtotal"
}