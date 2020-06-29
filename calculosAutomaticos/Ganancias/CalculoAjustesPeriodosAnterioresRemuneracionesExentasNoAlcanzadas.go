package Ganancias

type CalculoAjustesPeriodosAnterioresRemuneracionesExentasNoAlcanzadas struct {
	CalculoGanancias
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesExentasNoAlcanzadas) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("AJUSTES_PERÍODOS_ANTERIORES_REMUNERACIONES_EXENTAS_NO_ALCANZADAS")
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesExentasNoAlcanzadas) getResult() float64 {
	return cg.getResultOnDemandTemplate("AJUSTES_PERÍODOS_ANTERIORES_REMUNERACIONES_EXENTAS_NO_ALCANZADAS", 0, cg)
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesExentasNoAlcanzadas) getTope() *float64 {
	return nil
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesExentasNoAlcanzadas) getNombre() string {
	return "Ajustes Períodos Anteriores - Remuneraciones Exentas/No Alcanzadas"
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesExentasNoAlcanzadas) getEsMostrable() bool {
	return false
}
