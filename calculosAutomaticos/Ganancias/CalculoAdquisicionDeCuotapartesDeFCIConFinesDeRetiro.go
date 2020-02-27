package Ganancias

type CalculoAdquisicionDeCuotapartesDeFCIConFinesDeRetiro struct {
	CalculoGanancias
}

func (cg *CalculoAdquisicionDeCuotapartesDeFCIConFinesDeRetiro) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "ADQUISICION_DE_CUOTAPARTES_DE_FCI_CON_FINES_DE_RETIRO", "deducciondesgravacionsiradig")
	importeTope := *cg.getTope()
	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoAdquisicionDeCuotapartesDeFCIConFinesDeRetiro) getResult() float64 {
	return cg.getResultOnDemandTemplate("ADQUISICION_DE_CUOTAPARTES_DE_FCI_CON_FINES_DE_RETIRO", 0, cg)
}

func (cg *CalculoAdquisicionDeCuotapartesDeFCIConFinesDeRetiro) getTope() *float64 {
	importeTope := cg.getfgValorFijoImpuestoGanancia("topemaximodescuento", "toperetiroprivado")
	return &importeTope
}

func (cg *CalculoAdquisicionDeCuotapartesDeFCIConFinesDeRetiro) getNombre() string {
	return "Adquisición de Cuotapartes de FCI con fines de retiro"
}

func (cg *CalculoAdquisicionDeCuotapartesDeFCIConFinesDeRetiro) getEsMostrable() bool {
	return false
}
