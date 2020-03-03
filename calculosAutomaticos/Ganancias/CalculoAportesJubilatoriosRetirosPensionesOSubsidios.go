package Ganancias

type CalculoAportesJubilatoriosRetirosPensionesOSubsidios struct {
	CalculoGanancias
}

func (cg *CalculoAportesJubilatoriosRetirosPensionesOSubsidios) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("APORTES_JUBILATORIOS_RETIROS_PENSIONES_O_SUBSIDIOS")
}

func (cg *CalculoAportesJubilatoriosRetirosPensionesOSubsidios) getResult() float64 {
	return cg.getResultOnDemandTemplate("APORTES_JUBILATORIOS_RETIROS_PENSIONES_O_SUBSIDIOS", 16, cg)
}

func (cg *CalculoAportesJubilatoriosRetirosPensionesOSubsidios) getTope() *float64 {
	return nil
}

func (cg *CalculoAportesJubilatoriosRetirosPensionesOSubsidios) getNombre() string {
	return "Aportes jubilatorios, retiros, pensiones o subsidios (-)"
}

func (cg *CalculoAportesJubilatoriosRetirosPensionesOSubsidios) getEsMostrable() bool {
	return true
}
