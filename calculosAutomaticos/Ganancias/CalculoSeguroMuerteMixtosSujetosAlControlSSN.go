package Ganancias

type CalculoSeguroMuerteMixtosSujetosAlControlSSN struct {
	CalculoGanancias
}

func (cg *CalculoSeguroMuerteMixtosSujetosAlControlSSN) getResultInternal() float64 {
	return 0
}

func (cg *CalculoSeguroMuerteMixtosSujetosAlControlSSN) getResult() float64 {
	return cg.getResultOnDemandTemplate("Seguro muerte/mixtos sujetos al control de la SSN (-)", "SEGURO_MUERTE_MIXTOS_SUJETOS_AL_CONTROL_SSN", 1, cg)
}
