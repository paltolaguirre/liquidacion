package Ganancias

type CalculoDeduccionEspecial struct {
	CalculoGanancias
}

func (cg *CalculoDeduccionEspecial) getResultInternal() float64 {
	valorfijoMNI := cg.getfgValorFijoImpuestoGanancia( "deduccionespersonales", "valorfijoddei")
	mesperiodoliquidacion := getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion)
	return (valorfijoMNI / 12) * float64(mesperiodoliquidacion)
}

func (cg *CalculoDeduccionEspecial) getResult() float64 {
	return cg.getResultOnDemandTemplate("Deducción especial", "DEDUCCION_ESPECIAL", 42, cg)
}

func (cg *CalculoDeduccionEspecial) getTope() *float64 {
	return nil
}