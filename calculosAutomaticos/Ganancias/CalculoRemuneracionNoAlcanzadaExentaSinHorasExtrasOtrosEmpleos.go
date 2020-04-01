package Ganancias

type CalculoRemuneracionNoAlcanzadaExentaSinHorasExtrasOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoRemuneracionNoAlcanzadaExentaSinHorasExtrasOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig("importeconceptosexentos")
}

func (cg *CalculoRemuneracionNoAlcanzadaExentaSinHorasExtrasOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("REMUNERACION_NO_ALCANZADA_O_EXENTA_OTROS_EMPLEOS", 0, cg)
}

func (cg *CalculoRemuneracionNoAlcanzadaExentaSinHorasExtrasOtrosEmpleos) getTope() *float64 {
	return nil
}

func (cg *CalculoRemuneracionNoAlcanzadaExentaSinHorasExtrasOtrosEmpleos) getNombre() string {
	return "Remuneración No Alcanzada/Exenta Otros Empleos"
}

func (cg *CalculoRemuneracionNoAlcanzadaExentaSinHorasExtrasOtrosEmpleos) getEsMostrable() bool {
	return false
}
