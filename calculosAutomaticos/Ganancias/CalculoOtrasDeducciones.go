package Ganancias

type CalculoOtrasDeducciones struct {
	CalculoGanancias
}

func (cg *CalculoOtrasDeducciones) getResultInternal() float64 {
	return cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "OTRAS", "deducciondesgravacionsiradig")
}

func (cg *CalculoOtrasDeducciones) getResult() float64 {
	return cg.getResultOnDemandTemplate("Otras Deducciones", "OTRAS_DEDUCCIONES", 34, cg)
}

func (cg *CalculoOtrasDeducciones) getTope() *float64 {
	return nil
}