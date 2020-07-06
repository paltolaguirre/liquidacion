package Ganancias

type CalculoOtrasDeduccionesActoresRetribucionPagadaALosRepresentantesRG244208 struct {
	CalculoGanancias
}

func (cg *CalculoOtrasDeduccionesActoresRetribucionPagadaALosRepresentantesRG244208) getResultInternal() float64 {
	return cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "OTRAS_DEDUCCIONES_ACTORES_RETRIBUCION_PAGADA_A_LOS_REPRESENTANTES_RG_2442_08", "deducciondesgravacionsiradig")
}

func (cg *CalculoOtrasDeduccionesActoresRetribucionPagadaALosRepresentantesRG244208) getResult() float64 {
	return cg.getResultOnDemandTemplate("OTRAS_DEDUCCIONES_ACTORES_RETRIBUCION_PAGADA_A_LOS_REPRESENTANTES_RG_2442_08", 0, cg)
}

func (cg *CalculoOtrasDeduccionesActoresRetribucionPagadaALosRepresentantesRG244208) getTope() *float64 {
	return nil
}

func (cg *CalculoOtrasDeduccionesActoresRetribucionPagadaALosRepresentantesRG244208) getNombre() string {
	return "Otras Deducciones - Actores - Retribución pagada a los representantes - RG 2442/08"
}

func (cg *CalculoOtrasDeduccionesActoresRetribucionPagadaALosRepresentantesRG244208) getEsMostrable() bool {
	return false
}
