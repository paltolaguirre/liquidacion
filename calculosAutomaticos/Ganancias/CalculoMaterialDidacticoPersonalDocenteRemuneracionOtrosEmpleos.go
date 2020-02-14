package Ganancias


type CalculoMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos) getResultInternal() float64 {
	return cg.getfgImporteGananciasOtroEmpleoSiradig( "materialdidactico")
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("Material didáctico personal docente remuneración Otros Empleos (+)", "MATERIA_DIDACTICO_OTROS_EMPLEOS", 14, cg)
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos) getTope() *float64 {
	return nil
}
