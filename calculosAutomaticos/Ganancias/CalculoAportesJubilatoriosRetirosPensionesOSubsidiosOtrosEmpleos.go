package Ganancias

type CalculoAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos struct{
	CalculoGanancias
}

func (cg *CalculoAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos) getResultInternal() float64{
	return cg.getfgImporteGananciasOtroEmpleoSiradig( "aporteseguridadsocial")
}

func (cg *CalculoAportesJubilatoriosRetirosPensionesOSubsidiosOtrosEmpleos) getResult() float64{
	return cg.getResultOnDemandTemplate("Otros empleos - Aportes jubilatorios, retiros, pensiones o subsidios (-)", "APORTES_JUBILATORIOS_RETIROS_PENSIONES_O_SUBSIDIOS_OTROS_EMPLEOS", 19, cg)
}
