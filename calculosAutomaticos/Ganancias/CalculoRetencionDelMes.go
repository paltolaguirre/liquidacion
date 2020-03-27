package Ganancias

type CalculoRetencionDelMes struct {
	CalculoGanancias
}

func (cg *CalculoRetencionDelMes) getResultInternal() float64 {

	return (&CalculoTotalARetener{cg.CalculoGanancias}).getResult() - (&CalculoRetencionAcumulada{cg.CalculoGanancias}).getResult()
}

func (cg *CalculoRetencionDelMes) getResult() float64 {
	return cg.getResultOnDemandTemplate("RETENCION_DEL_MES", 54, cg)
}

func (cg *CalculoRetencionDelMes) getTope() *float64 {
	return nil
}

func (cg *CalculoRetencionDelMes) getNombre() string {
	return "Retencion del mes"
}

func (cg *CalculoRetencionDelMes) getEsMostrable() bool {
	return true
}
