package Ganancias

type CalculoOtrasDeduccionesFondosCompensadoresPrevision struct {
	CalculoGanancias
}

func (cg *CalculoOtrasDeduccionesFondosCompensadoresPrevision) getResultInternal() float64 {
	return cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "OTRAS_DEDUCCIONES_ACTORES_RETRIBUCION_PAGADA_A_LOS_REPRESENTANTES_RG_2442_08", "deducciondesgravacionsiradig")
}

func (cg *CalculoOtrasDeduccionesFondosCompensadoresPrevision) getResult() float64 {
	return cg.getResultOnDemandTemplate("OTRAS_DEDUCCIONES_FONDOS_COMPENSADORES_DE_PREVISION", 0, cg)
}

func (cg *CalculoOtrasDeduccionesFondosCompensadoresPrevision) getTope() *float64 {
	return nil
}

func (cg *CalculoOtrasDeduccionesFondosCompensadoresPrevision) getNombre() string {
	return "Otras Deducciones - Fondos Compensadores de Previsión"
}

func (cg *CalculoOtrasDeduccionesFondosCompensadoresPrevision) getEsMostrable() bool {
	return false
}
