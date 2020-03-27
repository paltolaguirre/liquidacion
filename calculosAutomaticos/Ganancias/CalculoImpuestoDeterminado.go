package Ganancias

type CalculoImpuestoDeterminado struct {
	CalculoGanancias
}

func (cg *CalculoImpuestoDeterminado) getResultInternal() float64 {
	var importeTotal float64 = 0
	strescalaimpuestoaplicable := *getfgEscalaImpuestoAplicable(cg.Liquidacion, cg.Db)
	importeRemuneracionSujetaImpuestoSinIncluirHorasExtras := (&CalculoRemuneracionSujetaAImpuestoSinIncluirHorasExtras{cg.CalculoGanancias}).getResult()

	for i := 0; i < len(strescalaimpuestoaplicable); i++ {
		escalaimpuestoaplicable := strescalaimpuestoaplicable[i]
		if importeRemuneracionSujetaImpuestoSinIncluirHorasExtras > escalaimpuestoaplicable.Limiteinferior && importeRemuneracionSujetaImpuestoSinIncluirHorasExtras <= escalaimpuestoaplicable.Limitesuperior {
			importeRemuneracionSujetaImpuesto := (&CalculoRemuneracionSujetaAImpuesto{cg.CalculoGanancias}).getResult()
			importeTotal = escalaimpuestoaplicable.Valorfijo + ((importeRemuneracionSujetaImpuesto - escalaimpuestoaplicable.Limiteinferior) * escalaimpuestoaplicable.Valorvariable)
		}
	}
	return importeTotal
}

func (cg *CalculoImpuestoDeterminado) getResult() float64 {
	return cg.getResultOnDemandTemplate("IMPUESTO_DETERMINADO", 0, cg)
}

func (cg *CalculoImpuestoDeterminado) getTope() *float64 {
	return nil
}

func (cg *CalculoImpuestoDeterminado) getNombre() string {
	return "Impuesto determinado"
}

func (cg *CalculoImpuestoDeterminado) getEsMostrable() bool {
	return false
}
