package Ganancias

type CalculoMovilidadYViaticosGravadaOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoMovilidadYViaticosGravadaOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig("gastosmovilidad")
}

func (cg *CalculoMovilidadYViaticosGravadaOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("MOVILIDAD_Y_VIATICOS_GRAVADA_OTROS_EMPLEOS", 13, cg)
}

func (cg *CalculoMovilidadYViaticosGravadaOtrosEmpleos) getTope() *float64 {
	return nil
}

func (cg *CalculoMovilidadYViaticosGravadaOtrosEmpleos) getNombre() string {
	return "Movilidad y Viáticos Gravada Otros Empleos (+)"
}