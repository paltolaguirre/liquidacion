package Ganancias

type CalculoCuotaSindicalOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoCuotaSindicalOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig( "aportesindical")
}

func (cg *CalculoCuotaSindicalOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("Cuota sindical – Otros empleos (-)", "CUOTA_SINDICAL_OTROS_EMPLEOS", 21, cg)
}