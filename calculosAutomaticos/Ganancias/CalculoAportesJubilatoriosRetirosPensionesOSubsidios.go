package Ganancias


type CalculoAportesJubilatoriosRetirosPensionesOSubsidios struct{
	CalculoGanancias
}

func (cg *CalculoAportesJubilatoriosRetirosPensionesOSubsidios) getResultInternal() float64{
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("APORTES_JUBILATORIOS_RETIROS_PENSIONES_O_SUBSIDIOS")
}

func (cg *CalculoAportesJubilatoriosRetirosPensionesOSubsidios) getResult() float64{
	return cg.getResultOnDemandTemplate("Aportes Jubilatorios Retiros, Pensiones o Subsidios", "APORTES_JUBILATORIOS_RETIROS_PENSIONES_O_SUBSIDIOS", 14, cg)
}