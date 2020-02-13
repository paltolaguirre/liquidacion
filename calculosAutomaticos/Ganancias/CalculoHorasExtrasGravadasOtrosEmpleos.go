package Ganancias

type CalculoHorasExtrasGravadasOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoHorasExtrasGravadasOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig( "importehorasextrasgravadas")
}

func (cg *CalculoHorasExtrasGravadasOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("Horas Extras Gravadas Otros Empleos (+)", "HORAS_EXTRAS_GRAVADAS_OTROS_EMPLEOS", 11, cg)
}
