package Ganancias

type CalculoAportesObraSocialOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoAportesObraSocialOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig("aporteobrasocial")
}

func (cg *CalculoAportesObraSocialOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("APORTES_OBRA_SOCIAL_OTROS_EMPLEOS", 20, cg)
}

func (cg *CalculoAportesObraSocialOtrosEmpleos) getTope() *float64 {
	return nil
}

func (cg *CalculoAportesObraSocialOtrosEmpleos) getNombre() string {
	return "Aportes obra social – Otros empleos (-)"
}

func (cg *CalculoAportesObraSocialOtrosEmpleos) getEsMostrable() bool {
	return true
}
