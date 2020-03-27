package Ganancias

type CalculoAlicuotaAplicableSinIncluirHorasExtras struct {
	CalculoGanancias
}

func (cg *CalculoAlicuotaAplicableSinIncluirHorasExtras) getResultInternal() float64 {
	var importeTotal float64 = 0
	strescalaimpuestoaplicable := *getfgEscalaImpuestoAplicable(cg.Liquidacion, cg.Db)
	importeRemuneracionSujetaImpuestoSinIncluirHorasExtras := (&CalculoRemuneracionSujetaAImpuestoSinIncluirHorasExtras{cg.CalculoGanancias}).getResult()

	for i := 0; i < len(strescalaimpuestoaplicable); i++ {
		escalaimpuestoaplicable := strescalaimpuestoaplicable[i]
		if importeRemuneracionSujetaImpuestoSinIncluirHorasExtras > escalaimpuestoaplicable.Limiteinferior && importeRemuneracionSujetaImpuestoSinIncluirHorasExtras <= escalaimpuestoaplicable.Limitesuperior {
			importeTotal = escalaimpuestoaplicable.Valorvariable * 100
		}
	}
	return importeTotal
}

func (cg *CalculoAlicuotaAplicableSinIncluirHorasExtras) getResult() float64 {
	return cg.getResultOnDemandTemplate("ALICUOTA_APLICABLE_SIN_INCLUIR_HORAS_EXTRAS", 0, cg)
}

func (cg *CalculoAlicuotaAplicableSinIncluirHorasExtras) getTope() *float64 {
	return nil
}

func (cg *CalculoAlicuotaAplicableSinIncluirHorasExtras) getNombre() string {
	return "Alicuota aplicable (sin incluir horas extras)"
}

func (cg *CalculoAlicuotaAplicableSinIncluirHorasExtras) getEsMostrable() bool {
	return false
}
