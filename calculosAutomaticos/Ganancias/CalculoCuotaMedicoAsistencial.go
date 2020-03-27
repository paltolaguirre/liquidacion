package Ganancias

type CalculoCuotaMedicoAsistencial struct {
	CalculoGanancias
}

func (cg *CalculoCuotaMedicoAsistencial) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrillaMesDesdeHasta("importe", "CUOTA_MEDICA_ASISTENCIAL", "deducciondesgravacionsiradig")
	var importeTope float64
	if importeTotal != 0 {
		importeTope = *cg.getTope()
	}

	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoCuotaMedicoAsistencial) getResult() float64 {
	return cg.getResultOnDemandTemplate("CUOTA_MEDICO_ASISTENCIAL", 36, cg)
}

func (cg *CalculoCuotaMedicoAsistencial) getTope() *float64 {
	importeTope := (&CalculoSubtotal{cg.CalculoGanancias}).getResult() * 0.05 //5% de Subtotal
	return &importeTope
}

func (cg *CalculoCuotaMedicoAsistencial) getNombre() string {
	return "Cuota médico asistencial (-)"
}

func (cg *CalculoCuotaMedicoAsistencial) getEsMostrable() bool {
	return true
}
