package Ganancias

type CalculoSACSegundaCuota struct {
	CalculoGanancias
}

func (cg *CalculoSACSegundaCuota) getResultInternal() float64 {
	correspondeSegundoSemetre := getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion) > 6
	return cg.getSac(correspondeSegundoSemetre)
}

func (cg *CalculoSACSegundaCuota) getResult() float64 {
	return cg.getResultOnDemandTemplate("SAC_SEGUNDA_CUOTA", 4, cg)
}

func (cg *CalculoSACSegundaCuota) getTope() *float64 {
	return nil
}

func (cg *CalculoSACSegundaCuota) getNombre() string {
	return "SAC Segunda Cuota (+)"
}

func (cg *CalculoSACSegundaCuota) getEsMostrable() bool {
	return true
}
