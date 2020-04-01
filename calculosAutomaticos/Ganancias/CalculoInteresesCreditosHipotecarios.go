package Ganancias

type CalculoInteresesCreditosHipotecarios struct {
	CalculoGanancias
}

func (cg *CalculoInteresesCreditosHipotecarios) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "INTERESES_PRESTAMO_HIPOTECARIO", "deducciondesgravacionsiradig")
	importeTope := *cg.getTope()
	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoInteresesCreditosHipotecarios) getResult() float64 {
	return cg.getResultOnDemandTemplate("INTERESES_CREDITOS_HIPOTECARIOS", 30, cg)
}

func (cg *CalculoInteresesCreditosHipotecarios) getTope() *float64 {
	importeTope := cg.getfgValorFijoImpuestoGanancia("topemaximodescuento", "topehipotecarios")
	return &importeTope
}

func (cg *CalculoInteresesCreditosHipotecarios) getNombre() string {
	return "Intereses créditos hipotecarios (-)"
}

func (cg *CalculoInteresesCreditosHipotecarios) getEsMostrable() bool {
	return true
}
