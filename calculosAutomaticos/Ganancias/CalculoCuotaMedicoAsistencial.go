package Ganancias

type CalculoCuotaMedicoAsistencial struct{
	CalculoGanancias
}

func (cg *CalculoCuotaMedicoAsistencial) getResultInternal() float64{
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrillaSinMes( "importe", "CUOTA_MEDICA_ASISTENCIAL", "deducciondesgravacionsiradig")
	var importeTope float64
	if importeTotal != 0 {
		importeTope = (&CalculoSubtotal{cg.CalculoGanancias}).getResult() * 0.05 //5% de Subtotal
	}

	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoCuotaMedicoAsistencial) getResult() float64{
	return cg.getResultOnDemandTemplate("Cuota médico asistencial (-)", "CUOTA_MEDICO_ASISTENCIAL", 34, cg)
}

