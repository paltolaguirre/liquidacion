package Ganancias

type CalculoDeduccionEspecial struct {
	CalculoGanancias
}

func (cg *CalculoDeduccionEspecial) getResultInternal() float64 {
	valorfijoMNI := cg.getfgValorFijoImpuestoGanancia("deduccionespersonales", "valorfijoddei")
	mesperiodoliquidacion := getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion)
	return (valorfijoMNI / 12) * float64(mesperiodoliquidacion)
}

func (cg *CalculoDeduccionEspecial) getResult() float64 {
	return cg.getResultOnDemandTemplate("DEDUCCION_ESPECIAL", 43, cg)
}

func (cg *CalculoDeduccionEspecial) getTope() *float64 {
	return nil
}

func (cg *CalculoDeduccionEspecial) getNombre() string {
	return "Deducción especial"
}

func (cg *CalculoDeduccionEspecial) getEsMostrable() bool {
	return true
}
