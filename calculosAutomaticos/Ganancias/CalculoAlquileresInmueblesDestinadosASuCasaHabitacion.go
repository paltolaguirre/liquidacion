package Ganancias

type CalculoAlquileresInmueblesDestinadosASuCasaHabitacion struct {
	CalculoGanancias
}

func (cg *CalculoAlquileresInmueblesDestinadosASuCasaHabitacion) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla( "importe", "ALQUILER_INMUEBLES_DESTINADOS_A_CASA_HABITACION", "deducciondesgravacionsiradig")
	importeTope := (&CalculoMinimoNoImponible{cg.CalculoGanancias}).getResult() * 0.4 /*es el 40% de MNI(40)*/
	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoAlquileresInmueblesDestinadosASuCasaHabitacion) getResult() float64 {
	return cg.getResultOnDemandTemplate("Alquileres de inmuebles destinados a su casa habitación (-)", "ALQUILERES_INMUEBLES_DESTINADOS_A_SU_CASA_HABITACION", 39, cg)
}
