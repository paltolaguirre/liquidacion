package Ganancias

type CalculoSubtotalCargasFamilia struct {
	CalculoGanancias
}

func (cg *CalculoSubtotalCargasFamilia) getResultInternal() float64 {

	var arraySubtotalCargasFamilia []float64

	arraySubtotalCargasFamilia = append(arraySubtotalCargasFamilia, (&CalculoConyugeAnual{cg.CalculoGanancias}).getResult())
	arraySubtotalCargasFamilia = append(arraySubtotalCargasFamilia, (&CalculoHijosAnual{cg.CalculoGanancias}).getResult())

	return Sum(arraySubtotalCargasFamilia)
}

func (cg *CalculoSubtotalCargasFamilia) getResult() float64 {
	return cg.getResultOnDemandTemplate("SUBTOTAL_CARGAS_FAMILIA", 0, cg)
}

func (cg *CalculoSubtotalCargasFamilia) getTope() *float64 {
	return nil
}

func (cg *CalculoSubtotalCargasFamilia) getNombre() string {
	return "Subtotal Cargas de Familia"
}

func (cg *CalculoSubtotalCargasFamilia) getEsMostrable() bool {
	return false
}
