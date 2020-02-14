package Ganancias

type CalculoHorasExtrasGravadas struct {
	CalculoGanancias
}

func (cg *CalculoHorasExtrasGravadas) getResultInternal() float64 {
	return cg.GetfgImporteTotalSegunTipoImpuestoGanancias("HORAS_EXTRAS_REMUNERACION_GRAVADA")
}

func (cg *CalculoHorasExtrasGravadas) getResult() float64 {
	return cg.getResultOnDemandTemplate("HORAS_EXTRAS_GRAVADAS", 5, cg)
}

func (cg *CalculoHorasExtrasGravadas) getTope() *float64 {
	return nil
}

func (cg *CalculoHorasExtrasGravadas) getNombre() string {
	return "Horas Extras Gravadas (+)"
}