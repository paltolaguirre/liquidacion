package Ganancias

type CalculoPrimasDeSeguroParaCasoDeMuerte struct {
	CalculoGanancias
}

func (cg *CalculoPrimasDeSeguroParaCasoDeMuerte) getResultInternal() float64 {
	return 0
}

func (cg *CalculoPrimasDeSeguroParaCasoDeMuerte) getResult() float64 {
	return cg.getResultOnDemandTemplate("Primas de seguro para el caso de muerte (-)", "PRIMAS_DE_SEGURO_PARA_CASO_DE_MUERTE", 1, cg)
}
