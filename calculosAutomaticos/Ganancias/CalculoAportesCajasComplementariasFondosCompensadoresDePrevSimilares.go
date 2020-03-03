package Ganancias

type CalculoAportesCajasComplementariasFondosCompensadoresDePrevSimilares struct {
	CalculoGanancias
}

func (cg *CalculoAportesCajasComplementariasFondosCompensadoresDePrevSimilares) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "APORTES_CAJAS_COMPLEMENTARIAS_FONDOS_COMPENSADORES_DE_PREV_SIMILARES", "deducciondesgravacionsiradig")
	return importeTotal
}

func (cg *CalculoAportesCajasComplementariasFondosCompensadoresDePrevSimilares) getResult() float64 {
	return cg.getResultOnDemandTemplate("APORTES_CAJAS_COMPLEMENTARIAS_FONDOS_COMPENSADORES_DE_PREV_SIMILARES", 0, cg)
}

func (cg *CalculoAportesCajasComplementariasFondosCompensadoresDePrevSimilares) getTope() *float64 {
	return nil
}

func (cg *CalculoAportesCajasComplementariasFondosCompensadoresDePrevSimilares) getNombre() string {
	return "Aportes Cajas complementarias / Fondos compensadores de prev. / Similares"
}

func (cg *CalculoAportesCajasComplementariasFondosCompensadoresDePrevSimilares) getEsMostrable() bool {
	return false
}
