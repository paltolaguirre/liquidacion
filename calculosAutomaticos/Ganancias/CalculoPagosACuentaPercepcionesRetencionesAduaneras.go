package Ganancias

type CalculoPagosACuentaPercepcionesRetencionesAduaneras struct {
	CalculoGanancias
}

func (cg *CalculoPagosACuentaPercepcionesRetencionesAduaneras) getResultInternal() float64 {
	return cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "PERCEPCIONES_RETENCIONES_ADUANERAS", "retencionpercepcionsiradig")
}

func (cg *CalculoPagosACuentaPercepcionesRetencionesAduaneras) getResult() float64 {
	return cg.getResultOnDemandTemplate("PERCEPCIONES_RETENCIONES_ADUANERAS", 0, cg)
}

func (cg *CalculoPagosACuentaPercepcionesRetencionesAduaneras) getTope() *float64 {
	return nil
}

func (cg *CalculoPagosACuentaPercepcionesRetencionesAduaneras) getNombre() string {
	return "PAGOS A CUENTA - Percepciones / Retenciones aduaneras"
}

func (cg *CalculoPagosACuentaPercepcionesRetencionesAduaneras) getEsMostrable() bool {
	return false
}
