package Ganancias

type CalculoGananciaNetaAcumSujetaAImp struct {
	CalculoGanancias
}

func (cg *CalculoGananciaNetaAcumSujetaAImp) getResultInternal() float64 {
	var importeTotal float64 = 0

	liquidacion := *cg.obtenerLiquidacionIgualAnioLegajoMesAnterior()
	importeTotal = importeTotal + (&CalculoGananciaNeta{cg.CalculoGanancias}).getResult()

	itemGanancia := obtenerItemGananciaFromLiquidacion(&liquidacion)
	if itemGanancia.ID != 0 {
		importeTotal = importeTotal + (&CalculoGananciaNetaAcumSujetaAImp{CalculoGanancias{itemGanancia, &liquidacion, cg.Db, false}}).getResult()
	}

	return importeTotal
}

func (cg *CalculoGananciaNetaAcumSujetaAImp) getResult() float64 {
	return cg.getResultOnDemandTemplate("GANANCIA_NETA_ACUM_SUJETA_A_IMP", 45, cg)
}

func (cg *CalculoGananciaNetaAcumSujetaAImp) getTope() *float64 {
	return nil
}

func (cg *CalculoGananciaNetaAcumSujetaAImp) getNombre() string {
	return "Ganancia neta acum. sujeta a Imp."
}

func (cg *CalculoGananciaNetaAcumSujetaAImp) getEsMostrable() bool {
	return true
}
