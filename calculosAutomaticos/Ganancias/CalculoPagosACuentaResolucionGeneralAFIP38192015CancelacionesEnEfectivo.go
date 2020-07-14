package Ganancias

type CalculoPagosACuentaResolucionGeneralAFIP38192015CancelacionesEnEfectivo struct {
	CalculoGanancias
}

func (cg *CalculoPagosACuentaResolucionGeneralAFIP38192015CancelacionesEnEfectivo) getResultInternal() float64 {
	return cg.getfgImporteTotalSiradigSegunTipoGrilla("importe", "PAGO_A_CUENTA_RESOLUCION_GENERAL_AFIP_3819_2015_CANCELACIONES_EN_EFECTIVO", "retencionpercepcionsiradig")
}

func (cg *CalculoPagosACuentaResolucionGeneralAFIP38192015CancelacionesEnEfectivo) getResult() float64 {
	return cg.getResultOnDemandTemplate("PAGO_A_CUENTA_RESOLUCION_GENERAL_AFIP_3819_2015_CANCELACIONES_EN_EFECTIVO", 0, cg)
}

func (cg *CalculoPagosACuentaResolucionGeneralAFIP38192015CancelacionesEnEfectivo) getTope() *float64 {
	return nil
}

func (cg *CalculoPagosACuentaResolucionGeneralAFIP38192015CancelacionesEnEfectivo) getNombre() string {
	return "PAGOS A CUENTA - Resolución General (AFIP) 3819/2015 - Cancelaciones en Efectivo"
}

func (cg *CalculoPagosACuentaResolucionGeneralAFIP38192015CancelacionesEnEfectivo) getEsMostrable() bool {
	return false
}
