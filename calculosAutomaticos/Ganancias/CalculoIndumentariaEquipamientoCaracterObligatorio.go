package Ganancias

type CalculoIndumentariaEquipamientoCaracterObligatorio struct {
	CalculoGanancias
}

func (cg *CalculoIndumentariaEquipamientoCaracterObligatorio) getResultInternal() float64 {
	return cg.getfgImporteTotalSiradigSegunTipoGrilla( "importe", "GASTOS_ADQUISICION_INDUMENTARIA_Y_EQUIPAMIENTO_PARA_USO_EXCLUSIVO_EN_EL_LUGAR_DE_TRABAJO", "deducciondesgravacionsiradig")
}

func (cg *CalculoIndumentariaEquipamientoCaracterObligatorio) getResult() float64 {
	return cg.getResultOnDemandTemplate("INDUMENTARIA_EQUIPAMIENTO_CARACTER_OBLIGATORIO", 33, cg)
}

func (cg *CalculoIndumentariaEquipamientoCaracterObligatorio) getTope() *float64 {
	return nil
}

func (cg *CalculoIndumentariaEquipamientoCaracterObligatorio) getNombre() string {
	return "Indumentaria/Equipamiento – uso exclusivo – carácter obligatorio (-)"
}