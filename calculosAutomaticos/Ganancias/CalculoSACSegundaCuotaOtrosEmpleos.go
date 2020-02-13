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
	return cg.getResultOnDemandTemplate("SAC Segunda Cuota Otros Empleos (+)", "SAC_SEGUNDA_CUOTA_OTROS_EMPLEOS", 3, cg)
}


