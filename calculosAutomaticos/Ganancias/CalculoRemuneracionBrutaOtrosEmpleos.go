package Ganancias

type CalculoRemuneracionBrutaOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoRemuneracionBrutaOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig("importegananciasbrutas")
}

func (cg *CalculoRemuneracionBrutaOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("Remuneración Bruta Otros Empleos (+)", "REMUNERACION_BRUTA_OTROS_EMPLEOS", 8, cg)
}

func (cg *CalculoRemuneracionBrutaOtrosEmpleos) getTope() *float64 {
	return nil
}