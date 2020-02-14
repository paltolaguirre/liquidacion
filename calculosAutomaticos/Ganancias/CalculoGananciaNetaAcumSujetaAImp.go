package Ganancias

type CalculoGananciaNetaAcumSujetaAImp struct {
	CalculoGanancias
}

func (cg *CalculoGananciaNetaAcumSujetaAImp) getResultInternal() float64 {
	var importeTotal float64 = 0

	liquidaciones := *cg.obtenerLiquidacionesIgualAnioLegajoMenorMes()
	importeTotal = importeTotal + (&CalculoGananciaNeta{cg.CalculoGanancias}).getResult()

	for i := 0; i < len(liquidaciones); i++ {
		itemGanancia := obtenerItemGananciaFromLiquidacion(&liquidaciones[i]);
		importeTotal = importeTotal + (&CalculoGananciaNetaAcumSujetaAImp{CalculoGanancias{itemGanancia, &liquidaciones[i], cg.Db}}).getResult()
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