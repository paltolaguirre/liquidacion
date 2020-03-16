package Ganancias

type CalculoMovilidadYViaticosRemuneracionExenta struct {
	CalculoGanancias
}

func (cg *CalculoMovilidadYViaticosRemuneracionExenta) getResultInternal() float64 {
	total := cg.GetfgImporteTotalSegunTipoImpuestoGanancias("MOVILIDAD_Y_VIATICOS_REMUNERACION_EXENTA", false)
	return total + cg.obtenerConceptosProrrateoMesesAnteriores()
}

func (cg *CalculoMovilidadYViaticosRemuneracionExenta) getResult() float64 {
	return cg.getResultOnDemandTemplate("MOVILIDAD_Y_VIATICOS_REMUNERACION_EXENTA", 0, cg)
}

func (cg *CalculoMovilidadYViaticosRemuneracionExenta) getTope() *float64 {
	return nil
}

func (cg *CalculoMovilidadYViaticosRemuneracionExenta) getNombre() string {
	return "Movilidad y Viaticos Remuneración Exenta"
}

func (cg *CalculoMovilidadYViaticosRemuneracionExenta) getEsMostrable() bool {
	return false
}
