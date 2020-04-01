package Ganancias

type CalculoImpuestoFijo struct {
	CalculoGanancias
}

func (cg *CalculoImpuestoFijo) getResultInternal() float64 {
	var importeTotal float64 = 0
	strescalaimpuestoaplicable := *getfgEscalaImpuestoAplicable(cg.Liquidacion, cg.Db)
	totalganancianeta := (&CalculoTotalGananciaNetaImponibleAcumuladaSinHorasExtras{cg.CalculoGanancias}).getResult()

	for i := 0; i < len(strescalaimpuestoaplicable); i++ {
		escalaimpuestoaplicable := strescalaimpuestoaplicable[i]
		if totalganancianeta > escalaimpuestoaplicable.Limiteinferior && totalganancianeta <= escalaimpuestoaplicable.Limitesuperior {
			importeTotal = escalaimpuestoaplicable.Valorfijo
		}
	}
	return importeTotal
}

func (cg *CalculoImpuestoFijo) getResult() float64 {
	return cg.getResultOnDemandTemplate("DETERMINACION_IMPUESTO_FIJO", 50, cg)
}

func (cg *CalculoImpuestoFijo) getTope() *float64 {
	return nil
}

func (cg *CalculoImpuestoFijo) getNombre() string {
	return "Determinacion de impuesto fijo"
}

func (cg *CalculoImpuestoFijo) getEsMostrable() bool {
	return true
}
