package Ganancias

type CalculoRemuneracionSujetaAImpuesto struct {
	CalculoGanancias
}

func (cg *CalculoRemuneracionSujetaAImpuesto) getResultInternal() float64 {
	return (&CalculoGananciaNetaAnual{cg.CalculoGanancias}).getResult() - (&CalculoSubtotalDeduccionesPersonalesAnual{cg.CalculoGanancias}).getResult()
}

func (cg *CalculoRemuneracionSujetaAImpuesto) getResult() float64 {
	return cg.getResultOnDemandTemplate("REMUNERACION_SUJETA_A_IMPUESTO", 0, cg)
}

func (cg *CalculoRemuneracionSujetaAImpuesto) getTope() *float64 {
	return nil
}

func (cg *CalculoRemuneracionSujetaAImpuesto) getNombre() string {
	return "Remuneracion sujeta a impuesto"
}

func (cg *CalculoRemuneracionSujetaAImpuesto) getEsMostrable() bool {
	return false
}
