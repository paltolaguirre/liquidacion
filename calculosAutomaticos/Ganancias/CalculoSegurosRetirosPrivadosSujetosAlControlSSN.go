package Ganancias

type CalculoSegurosRetirosPrivadosSujetosAlControlSSN struct {
	CalculoGanancias
}

func (cg *CalculoSegurosRetirosPrivadosSujetosAlControlSSN) getResultInternal() float64 {
	return 0
}

func (cg *CalculoSegurosRetirosPrivadosSujetosAlControlSSN) getResult() float64 {
	return cg.getResultOnDemandTemplate("SEGUROS_RETIROS_PRIVADOS_SUJETOS_AL_CONTROL_SSN", 26, cg)
}

func (cg *CalculoSegurosRetirosPrivadosSujetosAlControlSSN) getTope() *float64 {
	//ESTO TIENE TOPE PERO AUN NO SE IMPLEMENTO TODO
	return nil
}

func (cg *CalculoSegurosRetirosPrivadosSujetosAlControlSSN) getNombre() string {
	return "Seguro de retiro privados – Sujetos al control de la SSN (-)"
}

func (cg *CalculoSegurosRetirosPrivadosSujetosAlControlSSN) getEsMostrable() bool {
	return true
}
