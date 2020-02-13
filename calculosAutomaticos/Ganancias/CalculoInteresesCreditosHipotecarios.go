package Ganancias

type CalculoInteresesCreditosHipotecarios struct {
	CalculoGanancias
}

func (cg *CalculoInteresesCreditosHipotecarios) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "INTERESES_PRESTAMO_HIPOTECARIO", "deducciondesgravacionsiradig")
	importeTope := cg.getfgValorFijoImpuestoGanancia( "topemaximodescuento", "topehipotecarios")
	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoInteresesCreditosHipotecarios) getResult() float64 {
	return cg.getResultOnDemandTemplate("Intereses créditos hipotecarios (-)", "INTERESES_CREDITOS_HIPOTECARIOS", 39, cg)
}
