package Ganancias

type CalculoSubtotalDeduccionesPersonales struct {
	CalculoGanancias
}

func (cg *CalculoSubtotalDeduccionesPersonales) getResultInternal() float64 {
	var arraySubtotalDeduccionesPersonales []float64
	var subTotalDeduccionesPersonales float64

	arraySubtotalDeduccionesPersonales = append(arraySubtotalDeduccionesPersonales, (&CalculoConyuge{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesPersonales = append(arraySubtotalDeduccionesPersonales, (&CalculoHijos{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesPersonales = append(arraySubtotalDeduccionesPersonales, (&CalculoMinimoNoImponible{cg.CalculoGanancias}).getResult())
	arraySubtotalDeduccionesPersonales = append(arraySubtotalDeduccionesPersonales, (&CalculoDeduccionEspecial{cg.CalculoGanancias}).getResult())

	subTotalDeduccionesPersonales = Sum(arraySubtotalDeduccionesPersonales)
	return subTotalDeduccionesPersonales
}

func (cg *CalculoSubtotalDeduccionesPersonales) getResult() float64 {
	return cg.getResultOnDemandTemplate("SUBTOTAL_DEDUCCIONES_PERSONALES", 43, cg)
}

func (cg *CalculoSubtotalDeduccionesPersonales) getTope() *float64 {
	return nil
}

func (cg *CalculoSubtotalDeduccionesPersonales) getNombre() string {
	return "Subtotal Deducciones Personales"
}

func (cg *CalculoSubtotalDeduccionesPersonales) getEsMostrable() bool {
	return true
}
