package Ganancias

type CalculoPrimasDeSeguroParaElCasoDeMuerteAnual struct {
	CalculoGanancias
}

func (cg *CalculoPrimasDeSeguroParaElCasoDeMuerteAnual) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "PRIMAS_DE_SEGURO_PARA_EL_CASO_DE_MUERTE", "deducciondesgravacionsiradig")
	importeTope := *cg.getTope()
	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoPrimasDeSeguroParaElCasoDeMuerteAnual) getResult() float64 {
	return cg.getResultOnDemandTemplate("PRIMAS_DE_SEGURO_PARA_EL_CASO_DE_MUERTE", 0, cg)
}

func (cg *CalculoPrimasDeSeguroParaElCasoDeMuerteAnual) getTope() *float64 {
	importeTope := cg.getfgValorFijoImpuestoGanancia("topemaximodescuento", "topecasomuerte")
	return &importeTope
}

func (cg *CalculoPrimasDeSeguroParaElCasoDeMuerteAnual) getNombre() string {
	return "Primas de seguro para el caso de muerte Anual"
}

func (cg *CalculoPrimasDeSeguroParaElCasoDeMuerteAnual) getEsMostrable() bool {
	return false
}
