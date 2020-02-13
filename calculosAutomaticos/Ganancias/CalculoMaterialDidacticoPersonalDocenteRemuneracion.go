package Ganancias

type CalculoMaterialDidacticoPersonalDocenteRemuneracion struct {
	CalculoGanancias
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracion) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("MATERIAL_DIDACTICO_PERSONAL_DOCENTE_REMUNERACION_GRAVADA")
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracion) getResult() float64 {
	return cg.getResultOnDemandTemplate("Material didáctico personal docente remuneración (+)", "MATERIA_DIDACTICO", 13, cg)
}


