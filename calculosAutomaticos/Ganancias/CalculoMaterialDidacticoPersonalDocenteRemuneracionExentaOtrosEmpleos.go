package Ganancias

type CalculoMaterialDidacticoPersonalDocenteRemuneracionExentaOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracionExentaOtrosEmpleos) getResultInternal() float64 {
	return float64(0)
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracionExentaOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("MATERIAL_DIDACTICO_PERSONAL_DOCENTE_REMUNERACION_EXENTA", 0, cg)
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracionExentaOtrosEmpleos) getTope() *float64 {
	return nil
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracionExentaOtrosEmpleos) getNombre() string {
	return "Horas Extras Remuneración Exenta Otros Empleos"
}

func (cg *CalculoMaterialDidacticoPersonalDocenteRemuneracionExentaOtrosEmpleos) getEsMostrable() bool {
	return false
}
