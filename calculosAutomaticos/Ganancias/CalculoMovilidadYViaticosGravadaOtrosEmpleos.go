package Ganancias

type CalculoMovilidadYViaticosGravadaOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoMovilidadYViaticosGravadaOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig("gastosmovilidad")
}

func (cg *CalculoMovilidadYViaticosGravadaOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("Movilidad y Viáticos Gravada Otros Empleos (+)", "MOVILIDAD_Y_VIATICOS_GRAVADA_OTROS_EMPLEOS", 12, cg)
}

