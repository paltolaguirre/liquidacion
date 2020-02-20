package Ganancias

type CalculoAportesObraSocial struct{
	CalculoGanancias
}

func (cg *CalculoAportesObraSocial) getResultInternal() float64{
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("APORTES_OBRA_SOCIAL")
}

func (cg *CalculoAportesObraSocial) getResult() float64{
	return cg.getResultOnDemandTemplate("APORTES_OBRA_SOCIAL", 17, cg)
}

func (cg *CalculoAportesObraSocial) getTope() *float64 {
	return nil
}

func (cg *CalculoAportesObraSocial) getNombre() string {
	return "Aportes obra social (-)"
}