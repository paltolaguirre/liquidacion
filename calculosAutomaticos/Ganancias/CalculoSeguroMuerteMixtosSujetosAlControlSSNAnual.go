package Ganancias

type CalculoSeguroMuerteMixtosSujetosAlControlSSNAnual struct {
	CalculoGanancias
}

func (cg *CalculoSeguroMuerteMixtosSujetosAlControlSSNAnual) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "PRIMAS_DE_SEGURO_DE_AHORRO_O_MIXTO", "deducciondesgravacionsiradig")
	importeTope := *cg.getTope()
	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoSeguroMuerteMixtosSujetosAlControlSSNAnual) getResult() float64 {
	return cg.getResultOnDemandTemplate("PRIMAS_DE_SEGURO_DE_AHORRO_O_MIXTO", 0, cg)
}

func (cg *CalculoSeguroMuerteMixtosSujetosAlControlSSNAnual) getTope() *float64 {
	importeTope := cg.getfgValorFijoImpuestoGanancia("topemaximodescuento", "topeseguroahorro")
	return &importeTope
}

func (cg *CalculoSeguroMuerteMixtosSujetosAlControlSSNAnual) getNombre() string {
	return "Seguro muerte/mixtos sujetos al control de la SSN Anual"
}

func (cg *CalculoSeguroMuerteMixtosSujetosAlControlSSNAnual) getEsMostrable() bool {
	return false
}
