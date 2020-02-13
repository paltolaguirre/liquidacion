package Ganancias

type CalculoAportesObraSocial struct{
	CalculoGanancias
}

func (cg *CalculoAportesObraSocial) getResultInternal() float64{
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("APORTES_OBRA_SOCIAL")
}

func (cg *CalculoAportesObraSocial) getResult() float64{
	return cg.getResultOnDemandTemplate("Aportes Obra Social", "APORTES_OBRA_SOCIAL", 15, cg)
}
