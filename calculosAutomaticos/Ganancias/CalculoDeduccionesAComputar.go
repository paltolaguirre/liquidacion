package Ganancias

type CalculoDeduccionesAComputar struct{
	CalculoGanancias
}

func (cg *CalculoDeduccionesAComputar) getResultInternal() float64{
	return (&CalculoSubtotalDeduccionesPersonales{cg.CalculoGanancias}).getResult()
}

func (cg *CalculoDeduccionesAComputar) getResult() float64{
	return cg.getResultOnDemandTemplate("DEDUCCIONES_A_COMPUTAR", 44, cg)
}

func (cg *CalculoDeduccionesAComputar) getTope() *float64 {
	return nil
}

func (cg *CalculoDeduccionesAComputar) getNombre() string {
	return "Deducciones a Computar"
}