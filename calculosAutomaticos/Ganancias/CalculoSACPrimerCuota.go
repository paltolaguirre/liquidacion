package Ganancias

type CalculoSACPrimerCuota struct {
	CalculoGanancias
}

func (cg *CalculoSACPrimerCuota) getResultInternal() float64 {
	correspondePrimerSemetre := getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion) <= 6
	return cg.getfgSacCuotas(correspondePrimerSemetre)
}

func (cg *CalculoSACPrimerCuota) getResult() float64 {
	return cg.getResultOnDemandTemplate("SAC Primer Cuota (+)", "SAC_PRIMER_CUOTA", 3, cg)
}

func (cg *CalculoSACPrimerCuota) getTope() *float64 {
	return nil
}