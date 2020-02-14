package Ganancias


type CalculoMovilidadYViaticosGravada struct {
	CalculoGanancias
}

func (cg *CalculoMovilidadYViaticosGravada) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("MOVILIDAD_Y_VIATICOS_REMUNERACION_GRAVADA")
}

func (cg *CalculoMovilidadYViaticosGravada) getResult() float64 {
	return cg.getResultOnDemandTemplate("Movilidad y Viáticos Gravada (+)", "MOVILIDAD_Y_VIATICOS_GRAVADA", 6, cg)
}

func (cg *CalculoMovilidadYViaticosGravada) getTope() *float64 {
	return nil
}
