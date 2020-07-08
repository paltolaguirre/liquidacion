package Ganancias

type CalculoAjustesPeriodosAnterioresRemuneracionesGravadasOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesGravadasOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig("ajusteperiodoanteriorremgravada")
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesGravadasOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("AJUSTES_PERIODOS_ANTERIORES_REMUNERACIONES_GRAVADAS_OTROS_EMPLEOS", 0, cg)
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesGravadasOtrosEmpleos) getTope() *float64 {
	return nil
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesGravadasOtrosEmpleos) getNombre() string {
	return "Ajustes Periodos Anteriores Remuneraciones Gravadas Otros Empleos"
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesGravadasOtrosEmpleos) getEsMostrable() bool {
	return false
}
