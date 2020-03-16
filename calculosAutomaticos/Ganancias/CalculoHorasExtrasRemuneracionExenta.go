package Ganancias

type CalculoHorasExtrasRemuneracionExenta struct {
	CalculoGanancias
}

func (cg *CalculoHorasExtrasRemuneracionExenta) getResultInternal() float64 {
	total := cg.GetfgImporteTotalSegunTipoImpuestoGanancias("HORAS_EXTRAS_REMUNERACION_EXENTA", false)
	return total + cg.obtenerConceptosProrrateoMesesAnteriores()
}

func (cg *CalculoHorasExtrasRemuneracionExenta) getResult() float64 {
	return cg.getResultOnDemandTemplate("HORAS_EXTRAS_REMUNERACION_EXENTA", 0, cg)
}

func (cg *CalculoHorasExtrasRemuneracionExenta) getTope() *float64 {
	return nil
}

func (cg *CalculoHorasExtrasRemuneracionExenta) getNombre() string {
	return "Horas Extras Remuneración Exenta"
}

func (cg *CalculoHorasExtrasRemuneracionExenta) getEsMostrable() bool {
	return false
}
