package Ganancias

type CalculoAlicuotaArt90LeyGanancias struct {
	CalculoGanancias
}

func (cg *CalculoAlicuotaArt90LeyGanancias) getResultInternal() float64 {
	var importeTotal float64 = 0
	strescalaimpuestoaplicable := *getfgEscalaImpuestoAplicable(cg.Liquidacion, cg.Db)
	importeRemuneracionSujetaImpuesto := (&CalculoRemuneracionSujetaAImpuesto{cg.CalculoGanancias}).getResult()

	for i := 0; i < len(strescalaimpuestoaplicable); i++ {
		escalaimpuestoaplicable := strescalaimpuestoaplicable[i]
		if importeRemuneracionSujetaImpuesto > escalaimpuestoaplicable.Limiteinferior && importeRemuneracionSujetaImpuesto <= escalaimpuestoaplicable.Limitesuperior {
			importeTotal = escalaimpuestoaplicable.Valorfijo
		}
	}
	return importeTotal
}

func (cg *CalculoAlicuotaArt90LeyGanancias) getResult() float64 {
	return cg.getResultOnDemandTemplate("ALICUOTA_ART_90_LEY_GANANCIAS", 0, cg)
}

func (cg *CalculoAlicuotaArt90LeyGanancias) getTope() *float64 {
	return nil
}

func (cg *CalculoAlicuotaArt90LeyGanancias) getNombre() string {
	return "Alicuota art. 90 Ley Ganancias"
}

func (cg *CalculoAlicuotaArt90LeyGanancias) getEsMostrable() bool {
	return false
}
