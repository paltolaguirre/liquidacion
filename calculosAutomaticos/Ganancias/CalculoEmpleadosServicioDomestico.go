package Ganancias

type CalculoEmpleadosServicioDomestico struct {
	CalculoGanancias
}

func (cg *CalculoEmpleadosServicioDomestico) getResultInternal() float64 {
	importeTotal := cg.getfgImporteTotalSiradigSegunTipoGrilla("contribucion + retribucion", "DEDUCCION_DEL_PERSONAL_DOMESTICO", "deducciondesgravacionsiradig")
	importeTope := (&CalculoMinimoNoImponible{cg.CalculoGanancias}).getResult() /*es el MNI(40)*/
	return getfgImporteTotalTope(importeTotal, importeTope)
}

func (cg *CalculoEmpleadosServicioDomestico) getResult() float64 {
	return cg.getResultOnDemandTemplate("Empleados del servicio doméstico (-)", "EMPLEADOS_SERVICIO_DOMESTICO", 39, cg)
}
