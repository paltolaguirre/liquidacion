package Ganancias

type CalculoAlquileresInmueblesDestinadosASuCasaHabitacion struct {
	CalculoGanancias
}

func (cg *CalculoAlquileresInmueblesDestinadosASuCasaHabitacion) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "ALQUILER_INMUEBLES_DESTINADOS_A_CASA_HABITACION", "deducciondesgravacionsiradig") * 0.4
	importeTope := *cg.getTope()
	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoAlquileresInmueblesDestinadosASuCasaHabitacion) getResult() float64 {
	return cg.getResultOnDemandTemplate("ALQUILERES_INMUEBLES_DESTINADOS_A_SU_CASA_HABITACION", 32, cg)
}

func (cg *CalculoAlquileresInmueblesDestinadosASuCasaHabitacion) getTope() *float64 {
	importeTope := (&CalculoMinimoNoImponible{cg.CalculoGanancias}).getResult()
	return &importeTope
}

func (cg *CalculoAlquileresInmueblesDestinadosASuCasaHabitacion) getNombre() string {
	return "Alquileres de inmuebles destinados a su casa habitación (-)"
}

func (cg *CalculoAlquileresInmueblesDestinadosASuCasaHabitacion) getEsMostrable() bool {
	return true
}
