package Ganancias

type CalculoSegurosRetirosPrivadosSujetosAlControlSSN struct {
	CalculoGanancias
}

func (cg *CalculoSegurosRetirosPrivadosSujetosAlControlSSN) getResultInternal() float64 {
	return 0
}

func (cg *CalculoSegurosRetirosPrivadosSujetosAlControlSSN) getResult() float64 {
	return cg.getResultOnDemandTemplate("Seguro de retiro privados – Sujetos al control de la SSN (-)", "SEGUROS_RETIROS_PRIVADOS_SUJETOS_AL_CONTROL_SSN", 1, cg)
}
