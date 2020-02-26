package Ganancias

type CalculoSACSegundaCuotaOtrosEmpleos struct {
	CalculoGanancias
}

func (cg *CalculoSACSegundaCuotaOtrosEmpleos) getResultInternal() float64 {
	var importeTotal float64 = 0
	if getfgMes(&cg.Liquidacion.Fechaperiodoliquidacion) > 6 {
		importeTotal = cg.getfgImporteGananciasOtroEmpleoSiradig("sac")
	}
	return importeTotal
}

func (cg *CalculoSACSegundaCuotaOtrosEmpleos) getResult() float64 {
	return cg.getResultOnDemandTemplate("SAC_SEGUNDA_CUOTA_OTROS_EMPLEOS", 11, cg)
}

func (cg *CalculoSACSegundaCuotaOtrosEmpleos) getTope() *float64 {
	return nil
}

func (cg *CalculoSACSegundaCuotaOtrosEmpleos) getNombre() string {
	return "SAC Segunda Cuota Otros Empleos (+)"
}