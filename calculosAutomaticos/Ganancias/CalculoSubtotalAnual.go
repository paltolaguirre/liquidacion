package Ganancias

type CalculoSubtotalAnual struct {
	CalculoGanancias
}

func (cg *CalculoSubtotalAnual) getResultInternal() float64 {
	return (&CalculoSubtotalRemuneracionGravada{cg.CalculoGanancias}).getResult() - (&CalculoSubtotalDeduccionesGenerales{cg.CalculoGanancias}).getResult()
}

func (cg *CalculoSubtotalAnual) getResult() float64 {
	return cg.getResultOnDemandTemplate("SUBTOTAL_ANUAL", 35, cg)
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
