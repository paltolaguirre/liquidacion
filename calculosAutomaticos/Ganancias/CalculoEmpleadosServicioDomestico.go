package Ganancias

type CalculoEmpleadosServicioDomestico struct {
	CalculoGanancias
}

func (cg *CalculoEmpleadosServicioDomestico) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla("contribucion + retribucion", "DEDUCCION_DEL_PERSONAL_DOMESTICO", "deducciondesgravacionsiradig")
	importeTope := *cg.getTope()
	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoEmpleadosServicioDomestico) getResult() float64 {
	return cg.getResultOnDemandTemplate("EMPLEADOS_SERVICIO_DOMESTICO", 32, cg)
}

func (cg *CalculoEmpleadosServicioDomestico) getTope() *float64 {
	importeTope := (&CalculoMinimoNoImponible{cg.CalculoGanancias}).getResult() /*es el MNI(40)*/
	return &importeTope
}

func (cg *CalculoEmpleadosServicioDomestico) getNombre() string {
	return "Empleados del servicio doméstico (-)"
}