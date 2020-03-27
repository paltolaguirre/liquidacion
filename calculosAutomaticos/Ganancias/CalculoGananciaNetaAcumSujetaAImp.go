package Ganancias

type CalculoGananciaNetaAcumSujetaAImp struct {
	CalculoGanancias
}

func (cg *CalculoGananciaNetaAcumSujetaAImp) getResultInternal() float64 {
	return (&CalculoGananciaNeta{cg.CalculoGanancias}).getResult()
}

func (cg *CalculoGananciaNetaAcumSujetaAImp) getResult() float64 {
	return cg.getResultOnDemandTemplate("GANANCIA_NETA_ACUM_SUJETA_A_IMP", 46, cg)
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
