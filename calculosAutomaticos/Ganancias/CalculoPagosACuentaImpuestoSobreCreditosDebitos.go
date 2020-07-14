package Ganancias

type CalculoPagosACuentaImpuestoSobreCreditosDebitos struct {
	CalculoGanancias
}

func (cg *CalculoPagosACuentaImpuestoSobreCreditosDebitos) getResultInternal() float64 {
	return cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "IMPUESTOS_SOBRE_CREDITOS_Y_DEBITOS", "retencionpercepcionsiradig")
}

func (cg *CalculoPagosACuentaImpuestoSobreCreditosDebitos) getResult() float64 {
	return cg.getResultOnDemandTemplate("IMPUESTOS_SOBRE_CREDITOS_Y_DEBITOS", 0, cg)
}

func (cg *CalculoPagosACuentaImpuestoSobreCreditosDebitos) getTope() *float64 {
	return nil
}

func (cg *CalculoPagosACuentaImpuestoSobreCreditosDebitos) getNombre() string {
	return "PAGOS A CUENTA - Impuesto sobre créditos y débitos"
}

func (cg *CalculoPagosACuentaImpuestoSobreCreditosDebitos) getEsMostrable() bool {
	return false
}
