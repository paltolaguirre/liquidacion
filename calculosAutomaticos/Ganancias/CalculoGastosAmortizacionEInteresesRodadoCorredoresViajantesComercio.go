package Ganancias

type CalculoGastosAmortizacionEInteresesRodadoCorredoresViajantesComercio struct {
	CalculoGanancias
}

func (cg *CalculoGastosAmortizacionEInteresesRodadoCorredoresViajantesComercio) getResultInternal() float64 {
	return cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "GASTOS_DE_REPRESENTACION_E_INTERESES_DE_CORREDORES_Y_VIAJANTES_DE_COMERCIO", "deducciondesgravacionsiradig")
}

func (cg *CalculoGastosAmortizacionEInteresesRodadoCorredoresViajantesComercio) getResult() float64 {
	return cg.getResultOnDemandTemplate("Gastos amortización e intereses rodado, corredores y viajantes de comercio (-)", "GASTOS_AMORTIZACION_E_INTERESES_RODADO_CORREDORES_VIAJANTES_COMERCIO", 39, cg)
}
