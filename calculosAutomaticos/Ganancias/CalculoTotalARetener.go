package Ganancias

type CalculoTotalARetener struct {
	CalculoGanancias
}

func (cg *CalculoTotalARetener) getResultInternal() float64 {
	var arrayTotalRetener []float64
	var totalRetener float64

	arrayTotalRetener = append(arrayTotalRetener, (&CalculoImpuestoFijo{cg.CalculoGanancias}).getResult())
	arrayTotalRetener = append(arrayTotalRetener, (&CalculoImpuestoPorEscala{cg.CalculoGanancias}).getResult())

	totalRetener = Sum(arrayTotalRetener)
	return totalRetener
}

func (cg *CalculoTotalARetener) getResult() float64 {
	return cg.getResultOnDemandTemplate("TOTAL_A_RETENER", 51, cg)
}

func (cg *CalculoTotalARetener) getTope() *float64 {
	return nil
}

func (cg *CalculoTotalARetener) getNombre() string {
	return "Total a Retener"
}