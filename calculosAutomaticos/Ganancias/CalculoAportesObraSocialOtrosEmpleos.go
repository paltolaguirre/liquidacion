package Ganancias

type CalculoAportesObraSocialOtrosEmpleos struct{
	CalculoGanancias
}

func (cg *CalculoAportesObraSocialOtrosEmpleos) getResultInternal() float64{
	return cg.getfgImporteGananciasOtroEmpleoSiradig( "aporteobrasocial")
}

func (cg *CalculoAportesObraSocialOtrosEmpleos) getResult() float64{
	return cg.getResultOnDemandTemplate("Aportes obra social – Otros empleos (-)", "APORTES_OBRA_SOCIAL_OTROS_EMPLEOS", 22, cg)
}

func (cg *CalculoAportesObraSocialOtrosEmpleos) getTope() *float64 {
	return nil
}
