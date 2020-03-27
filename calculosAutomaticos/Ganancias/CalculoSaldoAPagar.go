package Ganancias

type CalculoSaldoAPagar struct {
	CalculoGanancias
}

func (cg *CalculoSaldoAPagar) getResultInternal() float64 {
	impuestoDeterminado := cg.roundTo((&CalculoImpuestoDeterminado{cg.CalculoGanancias}).getResult(), 4)
	impuestoRetenido := cg.roundTo((&CalculoImpuestoRetenido{cg.CalculoGanancias}).getResult(), 4)
	pagosACuenta := cg.roundTo((&CalculoPagosACuenta{cg.CalculoGanancias}).getResult(), 4)

	importeTotal := impuestoDeterminado - impuestoRetenido - pagosACuenta
	return importeTotal
}

func (cg *CalculoSaldoAPagar) getResult() float64 {
	return cg.getResultOnDemandTemplate("SALDO_A_PAGAR", 0, cg)
}

func (cg *CalculoSaldoAPagar) getTope() *float64 {
	return nil
}

func (cg *CalculoSaldoAPagar) getNombre() string {
	return "Saldo a Pagar"
}

func (cg *CalculoSaldoAPagar) getEsMostrable() bool {
	return false
}
