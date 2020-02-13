package Ganancias

type CalculoDeduccionesAComputar struct{
	CalculoGanancias
}

func (cg *CalculoDeduccionesAComputar) getResultInternal() float64{
	return (&CalculoSubtotalDeduccionesPersonales{cg.CalculoGanancias}).getResult()
}

func (cg *CalculoDeduccionesAComputar) getResult() float64{
	return cg.getResultOnDemandTemplate("Deducciones a Computar", "DEDUCCIONES_A_COMPUTAR", 45, cg)
}
