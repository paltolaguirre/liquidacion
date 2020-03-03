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
	return cg.getResultOnDemandTemplate("GASTOS_SEPELIO", 27, cg)
}

func (cg *CalculoGastosSepelio) getTope() *float64 {
	importeTope := cg.getfgValorFijoImpuestoGanancia("topemaximodescuento", "topesepelio")
	return &importeTope
}

func (cg *CalculoGastosSepelio) getNombre() string {
	return "Gastos de sepelio (-)"
}

func (cg *CalculoGastosSepelio) getEsMostrable() bool {
	return true
}
