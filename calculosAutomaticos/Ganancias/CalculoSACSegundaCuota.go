package Ganancias

type CalculoSACSegundaCuota struct {
	CalculoGanancias
}

func (cg *CalculoSACSegundaCuota) getResultInternal() float64 {
	correspondeSegundoSemetre := getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion) > 6
	return cg.getfgSacCuotas(correspondeSegundoSemetre)
}

func (cg *CalculoSACSegundaCuota) getResult() float64 {
	return cg.getResultOnDemandTemplate("SAC Segunda Cuota (+)", "SAC_SEGUNDA_CUOTA", 3, cg)
}

