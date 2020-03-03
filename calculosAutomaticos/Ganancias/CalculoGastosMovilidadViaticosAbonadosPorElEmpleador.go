package Ganancias

type CalculoGastosMovilidadViaticosAbonadosPorElEmpleador struct {
	CalculoGanancias
}

func (cg *CalculoGastosMovilidadViaticosAbonadosPorElEmpleador) getResultInternal() float64 {
	var arrayGastosMovilidad []float64

	arrayGastosMovilidad = append(arrayGastosMovilidad, (&CalculoMovilidadYViaticosGravada{cg.CalculoGanancias}).getResult())
	arrayGastosMovilidad = append(arrayGastosMovilidad, (&CalculoMovilidadYViaticosGravadaOtrosEmpleos{cg.CalculoGanancias}).getResult())
	arrayGastosMovilidad = append(arrayGastosMovilidad, (&CalculoMaterialDidacticoPersonalDocenteRemuneracion{cg.CalculoGanancias}).getResult())
	arrayGastosMovilidad = append(arrayGastosMovilidad, (&CalculoMaterialDidacticoPersonalDocenteRemuneracionOtrosEmpleos{cg.CalculoGanancias}).getResult())

	return Sum(arrayGastosMovilidad)
}

func (cg *CalculoGastosMovilidadViaticosAbonadosPorElEmpleador) getResult() float64 {
	return cg.getResultOnDemandTemplate("GASTOS_MOVILIDAD_VIATICOS_ABONADOS_POR_EL_EMPLEADOR", 20, cg)
}

func (cg *CalculoGastosMovilidadViaticosAbonadosPorElEmpleador) getTope() *float64 {
	//ESTE TIENE TOPE PERO AUN NO SE TIENE EN CUENTA POR DEFINICION TODO
	return nil
}

func (cg *CalculoGastosMovilidadViaticosAbonadosPorElEmpleador) getNombre() string {
	return "Gastos Movilidad Viaticos Abonados por el Empleador (-)"
}

func (cg *CalculoGastosMovilidadViaticosAbonadosPorElEmpleador) getEsMostrable() bool {
	return true
}
