package Ganancias

type CalculoCuotaSindicalOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoCuotaSindicalOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig("aportesindical")
}

func (cg *CalculoCuotaSindicalOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("CUOTA_SINDICAL_OTROS_EMPLEOS", 23, cg)
}

func (cg *CalculoCuotaSindicalOtrosEmpleos) getTope() *float64 {
	return nil
}

func (cg *CalculoCuotaSindicalOtrosEmpleos) getNombre() string {
	return "Cuota sindical – Otros empleos (-)"
}

func (cg *CalculoCuotaSindicalOtrosEmpleos) getEsMostrable() bool {
	return true
}
