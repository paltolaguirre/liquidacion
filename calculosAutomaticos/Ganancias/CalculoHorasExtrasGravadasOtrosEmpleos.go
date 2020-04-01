package Ganancias

type CalculoHorasExtrasGravadasOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoHorasExtrasGravadasOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig("importehorasextrasgravadas")
}

func (cg *CalculoHorasExtrasGravadasOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("HORAS_EXTRAS_GRAVADAS_OTROS_EMPLEOS", 12, cg)
}

func (cg *CalculoHorasExtrasGravadasOtrosEmpleos) getTope() *float64 {
	return nil
}

func (cg *CalculoHorasExtrasGravadasOtrosEmpleos) getNombre() string {
	return "Horas Extras Gravadas Otros Empleos (+)"
}

func (cg *CalculoHorasExtrasGravadasOtrosEmpleos) getEsMostrable() bool {
	return true
}
