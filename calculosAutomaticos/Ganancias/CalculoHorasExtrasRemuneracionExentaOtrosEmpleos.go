package Ganancias

type CalculoHorasExtrasRemuneracionExentaOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoHorasExtrasRemuneracionExentaOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig("importehorasextrasexentas")
}

func (cg *CalculoHorasExtrasRemuneracionExentaOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("HORAS_EXTRAS_REMUNERACION_EXENTA", 0, cg)
}

func (cg *CalculoHorasExtrasRemuneracionExentaOtrosEmpleos) getTope() *float64 {
	return nil
}

func (cg *CalculoHorasExtrasRemuneracionExentaOtrosEmpleos) getNombre() string {
	return "Horas Extras Remuneración Exenta Otros Empleos"
}

func (cg *CalculoHorasExtrasRemuneracionExentaOtrosEmpleos) getEsMostrable() bool {
	return false
}
