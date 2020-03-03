package Ganancias

type CalculoSeguroMuerteMixtosSujetosAlControlSSN struct {
	CalculoGanancias
}

func (cg *CalculoSeguroMuerteMixtosSujetosAlControlSSN) getResultInternal() float64 {
	return 0
}

func (cg *CalculoSeguroMuerteMixtosSujetosAlControlSSN) getResult() float64 {
	return cg.getResultOnDemandTemplate("SEGURO_MUERTE_MIXTOS_SUJETOS_AL_CONTROL_SSN", 25, cg)
}

func (cg *CalculoSeguroMuerteMixtosSujetosAlControlSSN) getTope() *float64 {
	//ESTO TIENE TOPE PERO AUN NO SE IMPLEMENTO TODO
	return nil
}

func (cg *CalculoSeguroMuerteMixtosSujetosAlControlSSN) getNombre() string {
	return "Seguro muerte/mixtos sujetos al control de la SSN (-)"
}

func (cg *CalculoSeguroMuerteMixtosSujetosAlControlSSN) getEsMostrable() bool {
	return true
}
