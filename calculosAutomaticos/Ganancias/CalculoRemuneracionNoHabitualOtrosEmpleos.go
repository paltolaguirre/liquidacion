package Ganancias

type CalculoRemuneracionNoHabitualOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoRemuneracionNoHabitualOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig("importeretribucionesnohabituales")
}

func (cg *CalculoRemuneracionNoHabitualOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("Remuneración No Habitual Otros Empleos (+)", "RETRIBUCIONES_NO_HABITUALES_OTROS_EMPLEOS", 8, cg)
}

