package Ganancias

type CalculoAjustesPeriodosAnterioresRemuneracionesExentasNoAlcanzadasOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesExentasNoAlcanzadasOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig("ajusteperiodoanteriorremexentanoalcanzada")
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesExentasNoAlcanzadasOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("AJUSTES_PERIODOS_ANTERIORES_REMUNERACIONES_EXENTAS_NO_ALCANZADAS_OTROS_EMPLEOS", 0, cg)
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesExentasNoAlcanzadasOtrosEmpleos) getTope() *float64 {
	return nil
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesExentasNoAlcanzadasOtrosEmpleos) getNombre() string {
	return "Ajustes Periodos Anteriores Remuneraciones Exentas / No Alcanzadas Otros Empleos"
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesExentasNoAlcanzadasOtrosEmpleos) getEsMostrable() bool {
	return false
}
