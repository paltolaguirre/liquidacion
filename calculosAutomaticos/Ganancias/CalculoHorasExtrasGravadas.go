package Ganancias

type CalculoHorasExtrasGravadas struct {
	CalculoGanancias
}

func (cg *CalculoHorasExtrasGravadas) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("MOVILIDAD_Y_VIATICOS_REMUNERACION_GRAVADA")
}

func (cg *CalculoHorasExtrasGravadas) getResult() float64 {
	return cg.getResultOnDemandTemplate("Horas Extras Gravadas (+)", "HORAS_EXTRAS_GRAVADAS", 5, cg)
}