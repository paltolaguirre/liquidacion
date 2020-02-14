package Ganancias

type CalculoImpuestoPorEscala struct {
	CalculoGanancias
}

func (cg *CalculoImpuestoPorEscala) getResultInternal() float64 {
	var importeTotal float64 = 0
	strescalaimpuestoaplicable := *getfgEscalaImpuestoAplicable(cg.Liquidacion, cg.Db)
	totalganancianeta := (&CalculoTotalGananciaNetaImponibleAcumuladaSinHorasExtras{cg.CalculoGanancias}).getResult()
	baseimponible := (&CalculoBaseImponible{cg.CalculoGanancias}).getResult()
	for i := 0; i < len(strescalaimpuestoaplicable); i++ {
		escalaimpuestoaplicable := strescalaimpuestoaplicable[i]
		if totalganancianeta > escalaimpuestoaplicable.Limiteinferior && totalganancianeta <= escalaimpuestoaplicable.Limitesuperior {

			importeTotal = (baseimponible - escalaimpuestoaplicable.Limiteinferior) * escalaimpuestoaplicable.Valorvariable
		}
	}
	return importeTotal
}

func (cg *CalculoImpuestoPorEscala) getResult() float64 {
	return cg.getResultOnDemandTemplate("DETERMINACION_IMPUESTO_POR_ESCALA", 50, cg)
}

func (cg *CalculoImpuestoPorEscala) getTope() *float64 {
	return nil
}

func (cg *CalculoImpuestoPorEscala) getNombre() string {
	return "Determinacion de Impuesto por Escala"
}
