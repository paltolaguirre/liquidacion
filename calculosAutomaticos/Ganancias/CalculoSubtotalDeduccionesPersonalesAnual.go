package Ganancias

type CalculoSubtotalDeduccionesPersonalesAnual struct {
	CalculoGanancias
}

func (cg *CalculoSubtotalDeduccionesPersonalesAnual) getResultInternal() float64 {
	var arraySubtotalDeduccionesPersonalesAnual []float64
	var subTotalDeduccionesPersonalesAnual float64

	arraySubtotalDeduccionesPersonalesAnual = append(arraySubtotalDeduccionesPersonalesAnual, (&CalculoMinimoNoImponible{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesPersonalesAnual = append(arraySubtotalDeduccionesPersonalesAnual, (&CalculoDeduccionEspecial{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesPersonalesAnual = append(arraySubtotalDeduccionesPersonalesAnual, (&CalculoConyugeAnual{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesPersonalesAnual = append(arraySubtotalDeduccionesPersonalesAnual, (&CalculoHijosAnual{cg.CalculoGanancias}).getResult())

	subTotalDeduccionesPersonalesAnual = Sum(arraySubtotalDeduccionesPersonalesAnual)
	return subTotalDeduccionesPersonalesAnual
}

func (cg *CalculoSubtotalDeduccionesPersonalesAnual) getResult() float64 {
	return cg.getResultOnDemandTemplate("SUBTOTAL_DEDUCCIONES_PERSONALES_ANUAL", 0, cg)
}

func (cg *CalculoSubtotalDeduccionesPersonalesAnual) getTope() *float64 {
	return nil
}

func (cg *CalculoSubtotalDeduccionesPersonalesAnual) getNombre() string {
	return "Subtotal Deducciones Personales Anual"
}

func (cg *CalculoSubtotalDeduccionesPersonalesAnual) getEsMostrable() bool {
	return false
}
