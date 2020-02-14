package Ganancias

type CalculoDeduccionesPersonales struct {
	CalculoGanancias
}

func (cg *CalculoDeduccionesPersonales) getResultInternal() float64 {
	return (&CalculoDeduccionesAComputar{cg.CalculoGanancias}).getResult() * -1
}

func (cg *CalculoDeduccionesPersonales) getResult() float64 {
	return cg.getResultOnDemandTemplate("Deducciones Personales", "DEDUCCIONES_PERSONALES", 46, cg)
}

func (cg *CalculoDeduccionesPersonales) getTope() *float64 {
	return nil
}