package Ganancias

type CalculoMaterialDidacticoPersonalDocenteRemuneracionExenta struct {
	CalculoGanancias
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracionExenta) getResultInternal() float64 {
	total := cg.GetfgImporteTotalSegunTipoImpuestoGanancias("MATERIAL_DIDACTICO_PERSONAL_DOCENTE_REMUNERACION_EXENTA", false)
	return total + cg.obtenerConceptosProrrateoMesesAnteriores()
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracionExenta) getResult() float64 {
	return cg.getResultOnDemandTemplate("MATERIAL_DIDACTICO_PERSONAL_DOCENTE_REMUNERACION_EXENTA", 0, cg)
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracionExenta) getTope() *float64 {
	return nil
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracionExenta) getNombre() string {
	return "Material Didactico Personal Docente Remuneración Exenta"
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracionExenta) getEsMostrable() bool {
	return false
}
