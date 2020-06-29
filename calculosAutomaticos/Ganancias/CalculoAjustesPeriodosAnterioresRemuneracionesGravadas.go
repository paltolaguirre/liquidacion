package Ganancias

type CalculoAjustesPeriodosAnterioresRemuneracionesGravadas struct {
	CalculoGanancias
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesGravadas) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("AJUSTES_PERÍODOS_ANTERIORES_REMUNERACIONES_GRAVADAS")
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesGravadas) getResult() float64 {
	return cg.getResultOnDemandTemplate("AJUSTES_PERÍODOS_ANTERIORES_REMUNERACIONES_GRAVADAS", 0, cg)
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesGravadas) getTope() *float64 {
	return nil
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesGravadas) getNombre() string {
	return "Ajustes Períodos Anteriores - Remuneraciones Gravadas"
}

func (cg *CalculoAjustesPeriodosAnterioresRemuneracionesGravadas) getEsMostrable() bool {
	return false
}
