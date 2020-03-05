package Ganancias

type CalculoSegurosRetirosPrivadosSujetosAlControlSSNAnual struct {
	CalculoGanancias
}

func (cg *CalculoSegurosRetirosPrivadosSujetosAlControlSSNAnual) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "APORTES_A_PLANES_DE_SEGURO_DE_RETIRO_PRIVADO", "deducciondesgravacionsiradig")
	importeTope := *cg.getTope()
	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoSegurosRetirosPrivadosSujetosAlControlSSNAnual) getResult() float64 {
	return cg.getResultOnDemandTemplate("APORTES_A_PLANES_DE_SEGURO_DE_RETIRO_PRIVADO", 0, cg)
}

func (cg *CalculoSegurosRetirosPrivadosSujetosAlControlSSNAnual) getTope() *float64 {
	importeTope := cg.getfgValorFijoImpuestoGanancia("topemaximodescuento", "toperetiroprivado")
	return &importeTope
}

func (cg *CalculoSegurosRetirosPrivadosSujetosAlControlSSNAnual) getNombre() string {
	return "Seguro de retiro privados – Sujetos al control de la SSN Anual"
}

func (cg *CalculoSegurosRetirosPrivadosSujetosAlControlSSNAnual) getEsMostrable() bool {
	return false
}
