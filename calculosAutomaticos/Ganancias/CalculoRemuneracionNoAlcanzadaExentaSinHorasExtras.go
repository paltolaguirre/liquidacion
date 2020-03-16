package Ganancias

type CalculoRemuneracionNoAlcanzadaExentaSinHorasExtras struct {
	CalculoGanancias
}

func (cg *CalculoRemuneracionNoAlcanzadaExentaSinHorasExtras) getResultInternal() float64 {
	total := cg.GetfgImporteTotalSegunTipoImpuestoGanancias("REMUNERACION_NO_ALCANZADA_O_EXENTA")
	return total + cg.obtenerConceptosProrrateoMesesAnteriores()
}

func (cg *CalculoRemuneracionNoAlcanzadaExentaSinHorasExtras) getResult() float64 {
	return cg.getResultOnDemandTemplate("REMUNERACION_NO_ALCANZADA_O_EXENTA", 0, cg)
}

func (cg *CalculoRemuneracionNoAlcanzadaExentaSinHorasExtras) getTope() *float64 {
	return nil
}

func (cg *CalculoRemuneracionNoAlcanzadaExentaSinHorasExtras) getNombre() string {
	return "Remuneración No Alcanzada o Exenta"
}

func (cg *CalculoRemuneracionNoAlcanzadaExentaSinHorasExtras) getEsMostrable() bool {
	return false
}
