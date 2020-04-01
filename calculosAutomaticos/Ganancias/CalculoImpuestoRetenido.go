package Ganancias

type CalculoImpuestoRetenido struct {
	CalculoGanancias
}

func (cg *CalculoImpuestoRetenido) getResultInternal() float64 {
	var arrayImpuestoRetenido []float64

	arrayImpuestoRetenido = append(arrayImpuestoRetenido, (&CalculoRetencionAcumulada{cg.CalculoGanancias}).getResult())
	arrayImpuestoRetenido = append(arrayImpuestoRetenido, (&CalculoRetencionDelMes{cg.CalculoGanancias}).getResult())

	return Sum(arrayImpuestoRetenido)
}

func (cg *CalculoImpuestoRetenido) getResult() float64 {
	return cg.getResultOnDemandTemplate("IMPUESTO_RETENIDO", 0, cg)
}

func (cg *CalculoImpuestoRetenido) getTope() *float64 {
	return nil
}

func (cg *CalculoImpuestoRetenido) getNombre() string {
	return "Impuesto Retenido"
}

func (cg *CalculoImpuestoRetenido) getEsMostrable() bool {
	return false
}
