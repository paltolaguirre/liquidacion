package Ganancias

type CalculoSubtotalRemuneracionGravada struct {
	CalculoGanancias
}

func (cg *CalculoSubtotalRemuneracionGravada) getResultInternal() float64 {
	var arraySubtotalRemuneracionGravada []float64
	var importeTotal float64
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoRemuneracionBruta{cg.CalculoGanancias}).getResult())
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoRemuneracionNoHabitual{cg.CalculoGanancias}).getResult())
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoSACPrimerCuota{cg.CalculoGanancias}).getResult())
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoSACSegundaCuota{cg.CalculoGanancias}).getResult())
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoHorasExtrasGravadas{cg.CalculoGanancias}).getResult())
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoMovilidadYViaticosGravada{cg.CalculoGanancias}).getResult())
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoMaterialDidacticoPersonalDocenteRemuneracion{cg.CalculoGanancias}).getResult())
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoRemuneracionBrutaOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoRemuneracionNoHabitualOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoSACPrimerCuotaOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoSACSegundaCuotaOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoHorasExtrasGravadasOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoMovilidadYViaticosGravadaOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arraySubtotalRemuneracionGravada = append(arraySubtotalRemuneracionGravada, (&CalculoMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos{cg.CalculoGanancias}).getResult())

	importeAcumuladorMesAnterior := cg.obtenerAcumuladorLiquidacionItemMesAnteriorSegunCodigo("SUBTOTAL_REMUNERACION_GRAVADA")
	importeTotal = Sum(arraySubtotalRemuneracionGravada) + importeAcumuladorMesAnterior

	return importeTotal
}

func (cg *CalculoSubtotalRemuneracionGravada) getResult() float64 {
	return cg.getResultOnDemandTemplate("SUBTOTAL_REMUNERACION_GRAVADA", 0, cg)
}

func (cg *CalculoSubtotalRemuneracionGravada) getTope() *float64 {
	return nil
}

func (cg *CalculoSubtotalRemuneracionGravada) getNombre() string {
	return "Subtotal Remuneración Gravada"
}

func (cg *CalculoSubtotalRemuneracionGravada) getEsMostrable() bool {
	return false
}
