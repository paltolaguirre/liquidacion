package Ganancias

type CalculoMaterialDidacticoPersonalDocenteRemuneracion struct {
	CalculoGanancias
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracion) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("MATERIAL_DIDACTICO_PERSONAL_DOCENTE_REMUNERACION_GRAVADA", false)
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracion) getResult() float64 {
	return cg.getResultOnDemandTemplate("MATERIA_DIDACTICO", 7, cg)
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracion) getTope() *float64 {
	return nil
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracion) getNombre() string {
	return "Material didáctico personal docente remuneración (+)"
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracion) getEsMostrable() bool {
	return true
}
