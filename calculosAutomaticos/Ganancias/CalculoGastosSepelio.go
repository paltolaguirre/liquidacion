package Ganancias

type CalculoGastosSepelio struct {
	CalculoGanancias
}

func (cg *CalculoGastosSepelio) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "GASTOS_DE_SEPELIO", "deducciondesgravacionsiradig")
	importeTope := *cg.getTope()
	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoGastosSepelio) getResult() float64 {
	return cg.getResultOnDemandTemplate("Gastos de sepelio", "GASTOS_SEPELIO", 27, cg)
}

func (cg *CalculoGastosSepelio) getTope() *float64 {
	importeTope := cg.getfgValorFijoImpuestoGanancia("topemaximodescuento", "topesepelio")
	return &importeTope
}