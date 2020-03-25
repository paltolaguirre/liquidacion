package Ganancias

type CalculoTotalRemuneraciones struct {
	CalculoGanancias
}

func (cg *CalculoTotalRemuneraciones) getResultInternal() float64 {
	var arrayTotalRemuneraciones []float64
	var importeTotal float64

	arrayTotalRemuneraciones = append(arrayTotalRemuneraciones, (&CalculoSubtotalRemuneracionGravada{cg.CalculoGanancias}).getResult())
	arrayTotalRemuneraciones = append(arrayTotalRemuneraciones, (&CalculoSubtotalRemuneracionNoGravadaNoAlcanzadaExenta{cg.CalculoGanancias}).getResult())

	importeAcumuladorMesAnterior := cg.obtenerAcumuladorLiquidacionItemMesAnteriorSegunCodigo("TOTAL_REMUNERACIONES")
	importeTotal = Sum(arrayTotalRemuneraciones) + importeAcumuladorMesAnterior
	return importeTotal
}

func (cg *CalculoTotalRemuneraciones) getResult() float64 {
	return cg.getResultOnDemandTemplate("TOTAL_REMUNERACIONES", 0, cg)
}

func (cg *CalculoTotalRemuneraciones) getTope() *float64 {
	return nil
}

func (cg *CalculoTotalRemuneraciones) getNombre() string {
	return "Total Remuneraciones"
}

func (cg *CalculoTotalRemuneraciones) getEsMostrable() bool {
	return true
}
