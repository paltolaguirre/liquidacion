package Ganancias

type CalculoIndumentariaEquipamientoCaracterObligatorio struct {
	CalculoGanancias
}

func (cg *CalculoIndumentariaEquipamientoCaracterObligatorio) getResultInternal() float64 {
	return cg.getfgImporteTotalSiradigSegunTipoGrilla( "importe", "GASTOS_ADQUISICION_INDUMENTARIA_Y_EQUIPAMIENTO_PARA_USO_EXCLUSIVO_EN_EL_LUGAR_DE_TRABAJO", "deducciondesgravacionsiradig")
}

func (cg *CalculoIndumentariaEquipamientoCaracterObligatorio) getResult() float64 {
	return cg.getResultOnDemandTemplate("Indumentaria/Equipamiento – uso exclusivo – carácter obligatorio (-)", "INDUMENTARIA_EQUIPAMIENTO_CARACTER_OBLIGATORIO", 39, cg)
}