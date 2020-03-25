package Ganancias

type CalculoSubtotalAnual struct {
	CalculoGanancias
}

func (cg *CalculoSubtotalAnual) getResultInternal() float64 {
	var importeTotal float64

	importeSubtotalAnual := (&CalculoSubtotalRemuneracionGravada{cg.CalculoGanancias}).getResult() - (&CalculoSubtotalDeduccionesGenerales{cg.CalculoGanancias}).getResult()
	importeAcumuladorMesAnterior := cg.obtenerAcumuladorLiquidacionItemMesAnteriorSegunCodigo("SUBTOTAL_ANUAL")
	importeTotal = importeSubtotalAnual + importeAcumuladorMesAnterior
	return importeTotal
}

func (cg *CalculoSubtotalAnual) getResult() float64 {
	return cg.getResultOnDemandTemplate("SUBTOTAL_ANUAL", 0, cg)
}

func (cg *CalculoSubtotalAnual) getTope() *float64 {
	return nil
}

func (cg *CalculoSubtotalAnual) getNombre() string {
	return "Subtotal Anual"
}

func (cg *CalculoSubtotalAnual) getEsMostrable() bool {
	return false
}
