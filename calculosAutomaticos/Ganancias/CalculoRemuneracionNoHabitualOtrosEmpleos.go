package Ganancias

type CalculoRemuneracionNoHabitualOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoRemuneracionNoHabitualOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig("importeretribucionesnohabituales")
}

func (cg *CalculoRemuneracionNoHabitualOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("RETRIBUCIONES_NO_HABITUALES_OTROS_EMPLEOS", 9, cg)
}

func (cg *CalculoRemuneracionNoHabitualOtrosEmpleos) getTope() *float64 {
	return nil
}

func (cg *CalculoRemuneracionNoHabitualOtrosEmpleos) getNombre() string {
	return "Remuneración No Habitual Otros Empleos (+)"
}