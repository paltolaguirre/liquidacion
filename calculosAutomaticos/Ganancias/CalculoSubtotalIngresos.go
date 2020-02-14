package Ganancias

type CalculoSubtotalIngresos struct {
	CalculoGanancias
}

func (cg *CalculoSubtotalIngresos) getResultInternal() float64 {
	var arraySubtotalIngresos []float64

	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoRemuneracionBruta{cg.CalculoGanancias}).getResult())
	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoRemuneracionNoHabitual{cg.CalculoGanancias}).getResult())
	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoSACPrimerCuota{cg.CalculoGanancias}).getResult())
	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoSACSegundaCuota{cg.CalculoGanancias}).getResult())
	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoHorasExtrasGravadas{cg.CalculoGanancias}).getResult())
	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoMovilidadYViaticosGravada{cg.CalculoGanancias}).getResult())
	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoMaterialDidacticoPersonalDocenteRemuneracion{cg.CalculoGanancias}).getResult())
	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoRemuneracionBrutaOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoRemuneracionNoHabitualOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoSACPrimerCuotaOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoSACSegundaCuotaOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoHorasExtrasGravadasOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoMovilidadYViaticosGravadaOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalIngresos = append(arraySubtotalIngresos, (&CalculoMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos{cg.CalculoGanancias}).getResult())

	return Sum(arraySubtotalIngresos)
}

func (cg *CalculoSubtotalIngresos) getResult() float64 {
	return cg.getResultOnDemandTemplate("SUBTOTAL_INGRESOS", 15, cg)
}

func (cg *CalculoSubtotalIngresos) getTope() *float64 {
	return nil
}

func (cg *CalculoSubtotalIngresos) getNombre() string {
	return "Subtotal Ingresos"
}
