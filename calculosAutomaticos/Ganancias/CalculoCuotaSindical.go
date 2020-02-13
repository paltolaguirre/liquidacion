package Ganancias

type CalculoCuotaSindical struct{
	CalculoGanancias
}

func (cg *CalculoCuotaSindical) getResultInternal() float64{
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("CUOTA_SINDICAL")
}

func (cg *CalculoCuotaSindical) getResult() float64{
	return cg.getResultOnDemandTemplate("Cuota Sindical", "CUOTA_SINDICAL", 16, cg)
}