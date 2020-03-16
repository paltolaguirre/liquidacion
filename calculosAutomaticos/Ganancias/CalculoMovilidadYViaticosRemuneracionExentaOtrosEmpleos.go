package Ganancias

type CalculoMovilidadYViaticosRemuneracionExentaOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoMovilidadYViaticosRemuneracionExentaOtrosEmpleos) getResultInternal() float64 {
	return float64(0)
}

func (cg *CalculoMovilidadYViaticosRemuneracionExentaOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("MOVILIDAD_Y_VIATICOS_REMUNERACION_EXENTA_OTROS_EMPLEOS", 0, cg)
}

func (cg *CalculoMovilidadYViaticosRemuneracionExentaOtrosEmpleos) getTope() *float64 {
	return nil
}

func (cg *CalculoMovilidadYViaticosRemuneracionExentaOtrosEmpleos) getNombre() string {
	return "Movilidad y Viaticos Remuneración Exenta Otros Empleos"
}

func (cg *CalculoMovilidadYViaticosRemuneracionExentaOtrosEmpleos) getEsMostrable() bool {
	return false
}
